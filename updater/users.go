package updater

type UserList struct {
	users []*User
}

type User struct {
	Name        string
	SlackId     string
	SlackName   string
	PagerdutyId string
	Email       string
}

func (u *UserList) ById(id string) *User {
	var usr *User

	for _, us := range u.users {
		if us.PagerdutyId == id {
			usr = us
			break
		}
	}

	return usr
}
