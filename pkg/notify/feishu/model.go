package feishu

type Element struct {
	Tag      string `json:"tag"`
	Text     string `json:"text"`
	Href     string `json:"href"`
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
}

type RichMsg struct {
	Title   string      `json:"title"`
	Content [][]Element `json:"content"`
}

type PostCnContent struct {
	ZhCn RichMsg `json:"zh_cn"`
}

type MsgContent struct {
	Post PostCnContent `json:"post"`
}

// Msg define feishu msg content
type Msg struct {
	MsgType string     `json:"msg_type"`
	Content MsgContent `json:"content"`
}
