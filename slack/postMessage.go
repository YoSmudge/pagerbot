package slack

import(
  "github.com/nlopes/slack"
)

func (a *Api) PostMessage(channel string, message string) error{
  mp := slack.NewPostMessageParameters()
  mp.Username = "Pagerduty Bot"
  mp.LinkNames = 1

  _, _, err := a.api.PostMessage(channel, message, mp)
  return err
}
