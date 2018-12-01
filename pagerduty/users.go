package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
)

type Users []User

type User struct {
	Id    string
	Name  string
	Email string
}

func (a *Api) Users() (Users, error) {
	var usr Users
	var opts pagerduty.ListUsersOptions

	for {
		res, err := a.client.ListUsers(opts)
		if err != nil {
			return usr, err
		}

		for _, user := range res.Users {
			usr = append(usr, User{
				Id:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			})
		}

		if !res.More {
			break
		}

		opts.Offset = res.Offset + res.Limit
	}

	return usr, nil
}
