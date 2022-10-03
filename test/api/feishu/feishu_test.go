package feishu

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/wadeling/linear-webhook/api/feishu"
	"os"
	"testing"
)

type TestConfig struct {
	Host      string `json:"host"`
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	Mobile    string `json:"mobile"`
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

func TestFetchAccessToken(t *testing.T) {
	c, err := LoadConfigFromFile("config.json")
	if err != nil {
		t.Fatalf("load config file err:%v", err)
	}
	t.Logf("config:%+v", c)

	api := feishu.NewApi(feishu.Config{Host: c.Host, AppId: c.AppId, AppSecret: c.AppSecret})
	err = api.FetchAccessToken()
	if err != nil {
		t.Fatalf("fetch access token failed:%v", err)
	}
	t.Logf("atoken %v", api.DumpAccessToken())
}

func TestFetchUserOpenId(t *testing.T) {
	c, err := LoadConfigFromFile("config.json")
	if err != nil {
		t.Fatalf("load config file err:%v", err)
	}
	t.Logf("config:%+v", c)

	api := feishu.NewApi(feishu.Config{Host: c.Host, AppId: c.AppId, AppSecret: c.AppSecret})
	err = api.FetchAccessToken()
	if err != nil {
		t.Fatalf("fetch access token failed:%v", err)
	}
	t.Logf("atoken %v", api.DumpAccessToken())

	openid, err := api.FetchUserOpenId(c.Mobile)
	if err != nil {
		t.Fatalf("fetch user open id err:%v", err)
	}
	t.Logf("openid %v", openid)
	assert.NotEqual(t, len(openid), 0)
}
