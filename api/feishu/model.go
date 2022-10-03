package feishu

// AccessTokenReq define request body when fetch access token
type AccessTokenReq struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

// AccessTokenRsp define response body when fetch access token
type AccessTokenRsp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

// UserOpenIdReq define request body when fetch user open id
type UserOpenIdReq struct {
	Emails  []string `json:"emails"`
	Mobiles []string `json:"mobiles"`
}

type User struct {
	UserId string `json:"user_id"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`
}

type UserData struct {
	UserList []User `json:"user_list"`
}

// UserOpenIdRsp define response body when fetch user open id
type UserOpenIdRsp struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data UserData `json:"data"`
}
