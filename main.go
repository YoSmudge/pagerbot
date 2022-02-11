package main

import (
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/qoharu/pagerbot/config"
	"github.com/qoharu/pagerbot/schedule"
	"github.com/qoharu/pagerbot/updater"
	"github.com/qoharu/pagerbot/user"
	"github.com/qoharu/pagerbot/usergroup"
	"github.com/slack-go/slack"
	"github.com/voxelbrain/goptions"
	"os"
)

type options struct {
	Verbose bool          `goptions:"-v, --verbose, description='Log verbosely'"`
	Help    goptions.Help `goptions:"-h, --help, description='Show help'"`
	Config  string        `goptions:"-c, --config, description='Config Yaml file to use'"`
}

func main() {
	parsedOptions := options{}

	parsedOptions.Config = "./config.yaml"

	goptions.ParseAndFail(&parsedOptions)

	if parsedOptions.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	log.Debug("Logging verbosely!")

	err := config.Load(parsedOptions.Config)
	if err != nil {
		log.WithFields(log.Fields{
			"configFile": parsedOptions.Config,
			"error":      err,
		}).Error("Could not load config file")
		os.Exit(1)
	}

	slackClient := slack.New(config.Config.ApiKeys.Slack)
	pdClient := pagerduty.NewClient(config.Config.ApiKeys.Pagerduty.Key)
	userService := user.NewService(slackClient, pdClient)
	scheduleService := schedule.NewService(userService, pdClient)
	userGroupService := usergroup.NewService(slackClient)
	updaterService := updater.NewService(scheduleService, userService, userGroupService)

	updaterService.SyncPagerDutySlackUserGroups()
}
