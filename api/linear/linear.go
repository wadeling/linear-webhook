package linear

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"github.com/rs/zerolog/log"
	"github.com/wadeling/linear-webhook/pkg/config"
	"time"
)

type Config struct {
	Host    string
	ApiKeys []config.LinApiKey
}

type Api struct {
	config Config
}

func (a *Api) GetApiKeyByWorkspace(workspaceName string) (string, error) {
	for _, v := range a.config.ApiKeys {
		if v.Workspace == workspaceName {
			return v.ApiKey, nil
		}
	}
	return "", fmt.Errorf("not found api key of workspace:%s", workspaceName)
}

// FetchIssueWithId get issue info with id,id eg."IVAN-10" or "ce064237-bb8c-4098-beda-xxxxx"
func (a *Api) FetchIssueWithId(workspaceName, issueId string) (IssueInfo, error) {
	client := graphql.NewClient(a.config.Host)
	client.Log = func(s string) {
		log.Debug().Msgf("%s", s)
	}
	// get api key
	apikey, err := a.GetApiKeyByWorkspace(workspaceName)
	if err != nil {
		log.Err(err).Msg("failed to get api key")
		return IssueInfo{}, err
	}

	// make a request
	req := graphql.NewRequest(QueryIssue)

	// set any variables
	req.Var("key", issueId)

	// set header fields
	req.Header.Set("Authorization", apikey)

	// define a Context for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(20)*time.Second)
	defer cancel()

	var respData IssueInfo
	if err := client.Run(ctx, req, &respData); err != nil {
		return IssueInfo{}, fmt.Errorf("clinet run err:%v", err)
	}

	log.Debug().Interface("issue", respData).Msg("fetch issue info ok")
	return respData, nil
}

func NewLinearApi(config Config) *Api {
	a := Api{
		config: config,
	}
	return &a
}
