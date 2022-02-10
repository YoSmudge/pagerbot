package updater

import (
	"github.com/qoharu/pagerbot/pagerduty"
)

type ScheduleList struct {
	schedules []*pagerduty.Schedule
}

func (s *ScheduleList) ById(id string) *pagerduty.Schedule {
	var schd *pagerduty.Schedule

	for _, sc := range s.schedules {
		if sc.Id == id {
			schd = sc
			break
		}
	}

	return schd
}
