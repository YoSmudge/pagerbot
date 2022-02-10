package slack

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/slack-go/slack"
)

type Api struct {
	api *slack.Client
}

func New(key string) (*Api, error) {
	a := Api{}

	a.api = slack.New(key)
	auth, err := a.api.AuthTest()
	if err != nil {
		return &a, fmt.Errorf("Error authenticating with Slack: %s", err)
	}

	log.WithFields(log.Fields{
		"teamName": auth.Team,
		"userId":   auth.UserID,
		"teamUrl":  auth.URL,
	}).Info("Authenticated with Slack")

	return &a, nil
}
