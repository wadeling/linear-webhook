package webhook

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"regexp"
)

var (
	reg = regexp.MustCompile(`.*linear\.app/(.*)/issue/([A-Za-z]+-\d+).*`)
)

const (
	PayloadTypeComment  = "Comment" // Capital 'c'
	PayloadTypeIssue    = "Issue"   //
	PayloadActionUpdate = "update"
	PayloadActionCreate = "create"
	PayloadActionRemove = "remove"
)

// Data define webhook content
type Data struct {
	Id         string `json:"id"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	ArchivedAt string `json:"archivedAt"`
	Body       string `json:"body"`
	Edited     bool   `json:"edited"`
	IssueId    string `json:"issueId"`
	UserId     string `json:"userId"`
}

// PayLoad define linear webhook payload
type PayLoad struct {
	Action    string `json:"action"`
	Type      string `json:"type"`
	Url       string `json:"url"`
	CreatedAt string `json:"createdAt"`
	Data      Data   `json:"data"`
}

// Issue contain some info that extracted from linear payload
// for payload.url: https://linear.app/my-workspace-name/issue/teamid-12#comment-bc123 ,
// workSpaceName is "my-workspace-name",issueId is "teamid-12"
type Issue struct {
	WorkSpaceName string // eg.my-workspace
	IssueId       string // eg.my
}

// GetInfo get info from linear payload
func GetInfo(payLoad PayLoad) (Issue, error) {
	if len(payLoad.Url) == 0 {
		// when remove a issue, url is empty
		log.Debug().Msg("payload url is nil,return empty msg")
		return Issue{}, nil
	}

	info := Issue{}
	match := reg.FindSubmatch([]byte(payLoad.Url))
	if len(match) < 3 {
		return Issue{}, fmt.Errorf("linear payload not match regexp")
	}
	info.WorkSpaceName = string(match[1])
	info.IssueId = string(match[2])
	log.Debug().Interface("webhookIssue", info).Msg("get info from payload")

	return info, nil
}
