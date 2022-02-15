package schedule

import (
	"context"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/qoharu/pagerbot/config"
	"github.com/qoharu/pagerbot/user"
)

type Service interface {
	GetActiveSchedule(id string) (Schedule, error)
}

type service struct {
	userService user.Service
	client      *pagerduty.Client
}

func NewService(userService user.Service, pdClient *pagerduty.Client) Service {
	log.Debugf("initiating schedule service")

	return service{
		userService: userService,
		client:      pdClient,
	}
}

func (s service) GetActiveSchedule(id string) (Schedule, error) {
	tz, err := time.LoadLocation(config.Config.TZ)
	if err != nil {
		return Schedule{}, err
	}

	currentTime := time.Now().In(tz)
	now := currentTime.Format("2006-01-02")
	oneHourLater := currentTime.Add(time.Hour * 1).Format("2006-01-02")

	pdSchedule, err := s.client.GetScheduleWithContext(
		context.Background(),
		id,
		pagerduty.GetScheduleOptions{
			TimeZone: "UTC",
			Since:    now,
			Until:    oneHourLater,
		},
	)

	log.Debugf("found schedule: %s user in schedule: %d", pdSchedule.FinalSchedule.Name, len(pdSchedule.FinalSchedule.RenderedScheduleEntries))

	if err != nil {
		return Schedule{}, err
	}

	var users []user.User
	for _, userReference := range pdSchedule.FinalSchedule.RenderedScheduleEntries {
		log.Debugf("getting user with pd user id: %s", userReference.User.ID)

		onCallUser, err := s.userService.GetUserByPagerDutyID(userReference.User.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"scheduleId": id,
				"pdUserID":   userReference.User.ID,
				"err":        err,
			}).Warning("Could not find user with this email")
		}

		log.Debugf("adding user for schedule id: %s | %s", pdSchedule.Name, onCallUser.Email)
		users = append(users, onCallUser)
	}

	return Schedule{
		ID:    pdSchedule.FinalSchedule.ID,
		Users: users,
	}, nil
}
