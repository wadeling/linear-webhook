package config

import (
	"github.com/spf13/viper"
	"path"
	"path/filepath"
)

type LinApiKey struct {
	Workspace string `mapstructure:"workspace"`
	ApiKey    string `mapstructure:"api_key"`
}

type Linear struct {
	ApiAddr string      `mapstructure:"api_addr"`
	ApiKeys []LinApiKey `mapstructure:"api_keys"`
}

type Feishu struct {
	WebhookUrl string `mapstructure:"webhook_url"`
	AppId      string `mapstructure:"app_id"`
	AppSecret  string `mapstructure:"app_secret"`
	StaffFile  string `mapstructure:"staff_file"` // staff file contain username->mobile, use mobile for feishu msg
}

type Tls struct {
	CertFile string `mapstructure:"certificate"`
	KeyFile  string `mapstructure:"private_key"`
}

type Server struct {
	ListenAddr string `mapstructure:"listen_addr"`
	Https      Tls    `mapstructure:"https"`
}

type config struct {
	Server Server `mapstructure:"server"`
	Linear Linear `mapstructure:"linear"`
	Feishu Feishu `mapstructure:"feishu"`
}

var (
	Config      config
	CfgFullPath string
)

func InitConfig(fullPath string) error {

	configFileName := path.Base(fullPath)
	lookPath := filepath.Dir(fullPath)
	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(lookPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		return err
	}
	return nil
}
