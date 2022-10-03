package feishu

import (
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

var (
	CustomUserAgent    = "tt-agent"
	defaultRetryCount  = 2
	retryMinTime       = 1
	retryMaxTime       = 5
	defaultHttpTimeout = 10
	client             *req.Client
	once               sync.Once
)

const (
	UserOpenIdApiAddr  = "https://open.feishu.cn/open-apis/contact/v3/users/batch_get_id"
	AccessTokenApiAddr = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
)

type Config struct {
	Host      string // feishu api server addr
	AppId     string // app id for fetch user info
	AppSecret string //
}

type Api struct {
	client      *req.Client
	accessToken string
	config      Config
}

func (a *Api) FetchAccessToken() error {
	r := AccessTokenReq{
		AppId:     a.config.AppId,
		AppSecret: a.config.AppSecret,
	}
	tokenRsp := AccessTokenRsp{}
	rsp, err := req.C().DevMode().
		SetTimeout(time.Duration(defaultHttpTimeout) * time.Second).
		R().
		SetContentType("application/json; charset=utf-8").
		SetBody(&r).
		SetResult(&tokenRsp).
		Post(AccessTokenApiAddr)

	if err != nil {
		log.Err(err).Msg("failed to fetch access token")
		return err
	}
	if !rsp.IsSuccess() {
		log.Error().Int("statusCode", rsp.GetStatusCode()).Msg("fetch access token rsp code err")
		return fmt.Errorf("fetch access token rsp code:%d", rsp.GetStatusCode())
	}
	if tokenRsp.Code != 0 {
		log.Error().Int("tokenRspCode", tokenRsp.Code).Msg("token rsp code not 0")
		return fmt.Errorf("token rsp code not 0.%s", tokenRsp.Code)
	}
	log.Trace().Interface("accessToken", tokenRsp).Msg("rsp detail")

	log.Debug().Msg("get access token ok")
	a.accessToken = tokenRsp.TenantAccessToken

	return nil
}

// DumpAccessToken for debug
func (a *Api) DumpAccessToken() string {
	return a.accessToken
}

func (a *Api) FetchUserOpenId(mobile string) (string, error) {
	r := UserOpenIdReq{
		Mobiles: []string{mobile},
	}
	openIdRsp := UserOpenIdRsp{}
	rsp, err := client.R().
		SetQueryParam("user_id_type", "open_id").
		SetBearerAuthToken(a.accessToken).
		SetContentType("application/json; charset=utf-8").
		SetBody(&r).
		SetResult(&openIdRsp).
		Post(UserOpenIdApiAddr)

	if err != nil {
		log.Err(err).Str("mobile", mobile).Msg("failed to fetch open id")
		return "", err
	}
	if !rsp.IsSuccess() {
		log.Error().Int("statusCode", rsp.GetStatusCode()).Msg("fetch open id rsp code err")
		return "", fmt.Errorf("fetch open id rsp code:%v", rsp.GetStatusCode())
	}

	// extract user open id
	if len(openIdRsp.Data.UserList) == 0 {
		log.Error().Interface("openIdRsp", openIdRsp).Msg("rsp user is empty")
		return "", fmt.Errorf("fetch open id rsp user is empty")
	}

	openId := openIdRsp.Data.UserList[0].UserId
	log.Debug().Str("openId", openId).Msg("fetch user open id ok")

	return openId, nil
}

func NewApi(config Config) *Api {
	a := &Api{
		config: config,
	}

	once.Do(func() {
		client = req.C().DevMode().
			EnableDebugLog().
			EnableDumpAll().
			EnableInsecureSkipVerify().
			EnableKeepAlives().
			SetCommonHeader("User-Agent", CustomUserAgent).
			SetTimeout(time.Duration(defaultHttpTimeout)*time.Second).
			SetCommonRetryCount(defaultRetryCount).
			SetCommonRetryBackoffInterval(time.Duration(retryMinTime)*time.Second, time.Duration(retryMaxTime)*time.Second)
	})

	return a
}
