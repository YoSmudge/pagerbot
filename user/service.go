package user

import (
	"context"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/slack-go/slack"
)

type Service interface {
	GetUserByPagerDutyID(id string) (User, error)
}

type service struct {
	slackClient *slack.Client
	pdClient    *pagerduty.Client
}

func NewService(slackClient *slack.Client, pdClient *pagerduty.Client) Service {
	log.Debugf("initiating user service")
	return service{
		slackClient: slackClient,
		pdClient:    pdClient,
	}
}

func (s service) GetUserByPagerDutyID(id string) (User, error) {
	pdUser, err := s.pdClient.GetUserWithContext(context.Background(), id, pagerduty.GetUserOptions{})
	if err != nil {
		return User{}, err
	}

	slackUser, err := s.slackClient.GetUserByEmailContext(context.Background(), pdUser.Email)
	if err != nil {
		return User{}, err
	}

	log.Debugf("success getting slack user with id %s email %s", slackUser.ID, pdUser.Email)

	return User{
		Email:           pdUser.Email,
		PagerdutyUserID: pdUser.ID,
		SlackUserID:     slackUser.ID,
	}, nil
}
