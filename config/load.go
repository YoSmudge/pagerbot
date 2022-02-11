package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var Config AppConfig

func Load(filePath string) error {
	Config = AppConfig{}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("Config file not found")
	}

	configContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configContent, &Config)
	if err != nil {
		return fmt.Errorf("Error parsing AppConfig file: %s", err)
	}

	if pdToken := os.Getenv("PAGERDUTY_TOKEN"); pdToken != "" {
		Config.ApiKeys.Pagerduty.Key = pdToken
	}

	if slackToken := os.Getenv("SLACK_TOKEN"); slackToken != "" {
		Config.ApiKeys.Slack = slackToken
	}

	return Config.Validate()
}
