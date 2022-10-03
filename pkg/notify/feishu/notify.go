package feishu

import (
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
	"github.com/wadeling/linear-webhook/api/feishu"
	"github.com/wadeling/linear-webhook/api/linear"
	"github.com/wadeling/linear-webhook/pkg/config"
	"github.com/wadeling/linear-webhook/pkg/notify"
	"github.com/wadeling/linear-webhook/pkg/staff"
	"github.com/wadeling/linear-webhook/pkg/webhook"
	"sync"
	"time"
)

// client a global client to avoid overhead in prod env
var (
	client *req.Client // client to send feishu msg
	once   sync.Once
)

type Notify struct {
	url       string        // feishu webhook address
	linConfig linear.Config // linear config
}

func init() {
	err := notify.Register(notify.TypeFeishu, NewNotify)
	if err != nil {
		log.Err(err).Msg("failed to register feishu notifier")
	} else {
		log.Info().Msg("success to register feishu notifier")
	}

	InitClient()
}

func InitClient() {
	once.Do(func() {
		client = req.C().
			EnableDebugLog().
			EnableDumpAll().
			EnableInsecureSkipVerify().
			SetCommonHeader("User-Agent", "ts-webhook").
			SetTimeout(time.Duration(defaultHttpTimeout)*time.Second).
			SetCommonRetryCount(defaultRetryCount).
			SetCommonRetryBackoffInterval(time.Duration(retryMinTime)*time.Second, time.Duration(retryMaxTime)*time.Second)
	})
}

func (n *Notify) Deliver(payload webhook.PayLoad) error {
	if len(n.url) == 0 {
		log.Error().Msg("not set notify url,ignore")
		return fmt.Errorf("miss notify url")
	}

	// get issue id from payload
	webhookIssue, err := webhook.GetInfo(payload)
	if err != nil {
		log.Err(err).Msg("failed to get issue info from payload")
		return err
	}

	// fetch linear issue detail
	issueInfo := linear.IssueInfo{}
	if len(webhookIssue.IssueId) != 0 {
		// when remove issue, issue id ni nil, not fetch detail
		if len(n.linConfig.Host) != 0 {
			linApi := linear.NewLinearApi(n.linConfig)
			issueInfo, err = linApi.FetchIssueWithId(webhookIssue.WorkSpaceName, webhookIssue.IssueId)
			if err != nil {
				log.Err(err).Msg("failed to fetch issue info")
			}
		} else {
			log.Warn().Msg("linear host is empty,not fetch issue detail")
		}
	}

	// transform msg
	msg := LinearMsgToFeishu(payload, webhookIssue, issueInfo)

	// notify
	resp, err := client.R().
		SetBody(&msg).
		Post(n.url)

	if err != nil {
		log.Err(err).Msg("failed to send notify to feishu")
		return err
	}

	if !resp.IsSuccess() {
		log.Error().Int("statusCode", resp.GetStatusCode()).Msg("feishu return err")
		return fmt.Errorf("status code %d", resp.GetStatusCode())
	}

	log.Debug().Msg("send notify to feishu success")
	return nil
}

func LinearMsgToFeishu(payLoad webhook.PayLoad, webhookIssue webhook.Issue, issueInfo linear.IssueInfo) Msg {
	m := Msg{}
	segment := make([]Element, 0)

	// add text msg
	et := TextElement(payLoad, webhookIssue, issueInfo)
	segment = append(segment, et)

	// add href
	e := Element{Tag: "a", Text: "issue 链接", Href: payLoad.Url}
	segment = append(segment, e)

	// add at user
	ea, err := MakeAtInfo(issueInfo.Issue.Assignee.DisplayName)
	if err != nil {
		log.Err(err).Msg("fetch user open id err, not append to msg")
	} else {
		segment = append(segment, ea)
	}

	content := make([][]Element, 0)
	content = append(content, segment)

	m.MsgType = "post"
	m.Content.Post.ZhCn.Title = "您有一条新的linear通知: "
	m.Content.Post.ZhCn.Content = content

	return m
}

func IssueDescription(payload webhook.PayLoad) string {
	description := ""
	if payload.Action == webhook.PayloadActionCreate {
		if payload.Type == webhook.PayloadTypeComment {
			description = "issue有新评论"
		} else if payload.Type == webhook.PayloadTypeIssue {
			description = "创建了新issue"
		} else {
			log.Warn().Str("type", payload.Type).Msg("action:create, unsupported type")
		}
	} else if payload.Action == webhook.PayloadActionUpdate {
		if payload.Type == webhook.PayloadTypeIssue {
			description = "更新了issue"
		} else if payload.Type == webhook.PayloadTypeComment {
			description = "更新了评论"
		} else {
			log.Warn().Str("type", payload.Type).Msg("action:update, unsupported type")
		}
	} else if payload.Action == webhook.PayloadActionRemove {
		description = "删除了issue"
	} else {
		log.Warn().Str("action", payload.Action).Msg("unsupported action")
	}

	return description
}

func TextElement(payload webhook.PayLoad, webhookIssue webhook.Issue, issueInfo linear.IssueInfo) Element {

	// feishu msg field
	description := IssueDescription(payload)
	workspaceName := webhookIssue.WorkSpaceName
	title := issueInfo.Issue.Title
	assignee := issueInfo.Issue.Assignee.DisplayName
	status := issueInfo.Issue.State.Name

	// format msg
	var data string
	type el struct {
		Name  string
		Value string
	}
	content := []el{
		{Name: "标    题:", Value: title},
		{Name: "描    述:", Value: description},
		{Name: "工作区:", Value: workspaceName},
		{Name: "负责人:", Value: assignee},
		{Name: "状    态:", Value: status},
	}
	for _, v := range content {
		data = data + fmt.Sprintf("%s    %s\n", v.Name, v.Value)
	}

	e := Element{
		Tag:  "text", //type
		Text: data,
	}

	return e
}

func MakeAtInfo(username string) (Element, error) {
	if len(username) == 0 {
		return Element{}, fmt.Errorf("empty username,ignore make at info")
	}

	// get user mobile
	mobile, err := staff.Instance().GetMobileByUserName(username)
	if err != nil {
		return Element{}, fmt.Errorf("get mobile err.%v", err)
	}

	// fetch access token
	api := feishu.NewApi(feishu.Config{AppId: config.Config.Feishu.AppId, AppSecret: config.Config.Feishu.AppSecret})
	err = api.FetchAccessToken()
	if err != nil {
		return Element{}, err
	}

	// fetch user info
	openId, err := api.FetchUserOpenId(mobile)
	if err != nil {
		return Element{}, err
	}

	e := Element{
		Tag:    "at",
		UserId: openId,
	}
	return e, nil
}

func NewNotify(config notify.Config) (notify.Notifier, error) {
	n := Notify{
		url:       config.Url,
		linConfig: config.LinConfig,
	}
	return &n, nil
}
