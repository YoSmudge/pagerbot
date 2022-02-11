package schedule

import (
	"github.com/qoharu/pagerbot/user"
)

type Schedule struct {
	ID    string
	Users []user.User
}
