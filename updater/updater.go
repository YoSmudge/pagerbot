package updater

import (
	"github.com/qoharu/pagerbot/config"
	"github.com/qoharu/pagerbot/pagerduty"
	"github.com/qoharu/pagerbot/slack"
	"sync"
	"time"
)

type Updater struct {
	Wg        *sync.WaitGroup
	Slack     *slack.Api
	Pagerduty *pagerduty.Api
	Users     *UserList
	Schedules *ScheduleList
	LastFetch time.Time
}

func New() (*Updater, error) {
	u := Updater{}
	u.Wg = &sync.WaitGroup{}

	var err error
	u.Slack, err = slack.New(config.Config.ApiKeys.Slack)
	if err != nil {
		return &u, err
	}

	u.Pagerduty, err = pagerduty.New(config.Config.ApiKeys.Pagerduty.Key, config.Config.ApiKeys.Pagerduty.Org)
	if err != nil {
		return &u, err
	}

	u.Users = &UserList{}
	u.Schedules = &ScheduleList{}

	return &u, nil
}
