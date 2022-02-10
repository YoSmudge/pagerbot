package config

type config struct {
	ApiKeys struct {
		Slack     string
		Pagerduty struct {
			Key string
			Org string
		}
	} `yaml:"api_keys"`
	Groups []struct {
		Name          string
		Schedules     []string
		UpdateMessage struct {
			Message  string
			Channels []string
		} `yaml:"update_message"`
	}
}
