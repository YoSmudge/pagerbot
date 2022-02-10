package slack

import (
	"github.com/slack-go/slack"
)

func (a *Api) Users() ([]slack.User, error) {
	var usr []slack.User
	usr, err := a.api.GetUsers()
	if err != nil {
		return usr, err
	}

	return usr, nil
}
