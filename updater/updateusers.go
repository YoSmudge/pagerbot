package updater

import (
	log "github.com/Sirupsen/logrus"
	"github.com/slack-go/slack"
	"sync"
)

// Fetch the users from Pagerduty and slack, and make sure we can match them
// all up. We match Pagerduty users to Slack users based on their email address
func (u *Updater) updateUsers(w *sync.WaitGroup) {
	defer w.Done()

	var err error
	pdUsers, err := u.Pagerduty.Users()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warning("Could not fetch users from Pagerduty")
		return
	}

	slackUsers, err := u.Slack.Users()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warning("Could not fetch users from Slack")
		return
	}

	// Create a map of slack email -> user for searching
	slackUserMap := make(map[string]*slack.User)
	for i, _ := range slackUsers {
		u := slackUsers[i]
		slackUserMap[u.Profile.Email] = &u
	}

	var users UserList
	for _, u := range pdUsers {
		su, found := slackUserMap[u.Email]
		if !found {
			log.WithFields(log.Fields{
				"email":       u.Email,
				"pagerdutyId": u.Id,
			}).Warning("Could not find Slack account for Pagerduty user")
			continue
		}

		usr := User{}
		usr.Name = u.Name
		usr.SlackId = su.ID
		usr.PagerdutyId = u.Id
		usr.SlackName = su.Name
		usr.Email = u.Email
		users.users = append(users.users, &usr)
	}

	u.Users = &users
}
