package usergroup

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/qoharu/pagerbot/user"
	"github.com/slack-go/slack"
	"strings"
)

type Service interface {
	UpdateUserGroup(group string, users []user.User) error
}

type service struct {
	slackClient *slack.Client

	userGroupsMap map[string]slack.UserGroup
}

func NewService(slackClient *slack.Client) Service {
	// hacky way to cache user groups list
	userGroupsMap := make(map[string]slack.UserGroup)
	log.Debugf("initiating user group service")
	userGroups, err := slackClient.GetUserGroups()
	if err != nil {
		panic(err)
	}
	for _, group := range userGroups {
		userGroupsMap[group.Handle] = group
	}

	log.Debugf("user groups loaded, total: %d", len(userGroups))

	return service{
		slackClient:   slackClient,
		userGroupsMap: userGroupsMap,
	}
}

func (s service) UpdateUserGroup(group string, users []user.User) error {
	userGroup, err := s.GetUserGroup(group)
	if err != nil {
		return err
	}
	log.Debugf("success get slack user group %s|%s", userGroup.Handle, userGroup.Name)

	usersCSV := strings.Join(convertUserListToSlackUserIds(users), ",")

	log.Debugf("new users for on-call: %v", usersCSV)
	userGroup, err = s.slackClient.UpdateUserGroupMembers(
		userGroup.ID,
		usersCSV,
	)

	log.Debugf("updated user group: %v", userGroup.Users)
	return err
}

func (s service) GetUserGroup(group string) (slack.UserGroup, error) {
	userGroup, found := s.userGroupsMap[group]
	if !found {
		return slack.UserGroup{}, fmt.Errorf("usergroup with handle %s is not found", group)
	}

	return userGroup, nil
}

func convertUserListToSlackUserIds(users []user.User) []string {
	var slackUserIds []string
	for _, u := range users {
		slackUserIds = append(slackUserIds, u.SlackUserID)
	}

	return slackUserIds
}
