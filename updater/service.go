package updater

import (
	log "github.com/Sirupsen/logrus"
	"github.com/qoharu/pagerbot/config"
	"github.com/qoharu/pagerbot/schedule"
	"github.com/qoharu/pagerbot/user"
	"github.com/qoharu/pagerbot/usergroup"
)

type Service interface {
	SyncPagerDutySlackUserGroups()
}

type service struct {
	scheduleService  schedule.Service
	userService      user.Service
	userGroupService usergroup.Service
}

func NewService(scheduleService schedule.Service, userService user.Service, userGroupService usergroup.Service) Service {
	log.Debugf("initiating updater service")

	return service{
		scheduleService:  scheduleService,
		userService:      userService,
		userGroupService: userGroupService,
	}
}

func (s service) SyncPagerDutySlackUserGroups() {
	log.Debugf("starting sync pagerduty schedule to slack")

	configuration := config.Config
	for _, group := range configuration.Groups {
		var onCallUsers []user.User
		for _, scheduleID := range group.Schedules {
			sched, err := s.scheduleService.GetActiveSchedule(scheduleID)
			if err != nil {
				log.WithFields(log.Fields{
					"scheduleId": scheduleID,
					"err":        err,
				}).Error("Could not find schedule")
				continue
			}

			onCallUsers = append(onCallUsers, sched.Users...)
		}
		err := s.userGroupService.UpdateUserGroup(group.Name, onCallUsers)
		if err != nil {
			log.WithFields(log.Fields{
				"group": group.Name,
				"err":   err,
			}).Error("Could not update user groups")
		}
	}
}
