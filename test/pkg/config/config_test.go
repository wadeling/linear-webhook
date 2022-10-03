package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/wadeling/linear-webhook/pkg/config"
	"testing"
)

func TestReadConfig(t *testing.T) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		t.Fatalf("load config err:%v", err)
	}
	err = viper.Unmarshal(&config.Config)
	if err != nil {
		t.Fatalf("unmarshall config err:%v", err)
	}
	t.Logf("%+v", config.Config)
	assert.NotEqual(t, len(config.Config.Linear.ApiKeys), 0)
}
