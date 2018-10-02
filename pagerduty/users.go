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

	res, err := a.client.ListUsers(pagerduty.ListUsersOptions{})
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

	return usr, nil
}
