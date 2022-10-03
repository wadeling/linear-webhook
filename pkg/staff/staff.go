package staff

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

var (
	instance *Staff
)

type Config struct {
	UserName string `json:"user_name"`
	Mobile   string `json:"mobile"`
}

type Staff struct {
	userInfo map[string]string
	lock     sync.Mutex
}

func init() {
	instance = &Staff{
		userInfo: make(map[string]string),
	}
}

func Instance() *Staff {
	return instance
}

func (s *Staff) LoadStaffInfoFromFile(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Err(err).Str("config", configFile).Msg("failed to load staff config file")
		return err
	}

	var configs []Config
	err = json.Unmarshal(data, &configs)
	if err != nil {
		log.Err(err).Msg("parse config file err")
		return err
	}

	for _, v := range configs {
		s.AddUserInfo(v)
	}
	return nil
}

func (s *Staff) AddUserInfo(config Config) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.userInfo[config.UserName] = config.Mobile
}

// DumpUserInfo for debug
func (s *Staff) DumpUserInfo() {
	for k, v := range s.userInfo {
		log.Debug().Str("userName", k).Str("mobile", v).Msgf("dump user")
	}
}

func (s *Staff) GetMobileByUserName(username string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	m, ok := s.userInfo[username]
	if !ok {
		return "", fmt.Errorf("not found user:%v", username)
	}
	return m, nil
}
