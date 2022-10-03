package linear

import (
	"encoding/json"
	"github.com/wadeling/linear-webhook/api/linear"
	"github.com/wadeling/linear-webhook/pkg/config"
	"os"
	"testing"
)

type TestConfig struct {
	Host      string `json:"host"`
	ApiKey    string `json:"api_key"`
	IssueId   string `json:"issue_id"`
	Workspace string `json:"workspace"`
}

func LoadConfigFromFile(configFile string) (TestConfig, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return TestConfig{}, err
	}
	c := TestConfig{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return TestConfig{}, err
	}
	return c, nil
}

func TestFetchIssue(t *testing.T) {
	c, err := LoadConfigFromFile("config.json")
	if err != nil {
		t.Fatalf("load config file err:%v", err)
	}
	t.Logf("config:%+v", c)

	apiKeys := make([]config.LinApiKey, 0)
	apiKeys = append(apiKeys, config.LinApiKey{Workspace: c.Workspace, ApiKey: c.ApiKey})
	api := linear.NewLinearApi(linear.Config{Host: c.Host, ApiKeys: apiKeys})
	rsp, err := api.FetchIssueWithId(c.Workspace, c.IssueId)
	if err != nil {
		t.Fatalf("fetch issue failed:%v", err)
	}
	t.Logf("issue info:%v", rsp)
}
