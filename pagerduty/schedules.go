package pagerduty

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
)

type Schedules []Schedule

type Schedule struct {
	Id            string
	Name          string
	Timezone      string `json:"time_zone"`
	CurrentPeriod *CallPeriod
	NextPeriod    *CallPeriod
}

type CallPeriod struct {
	Start time.Time
	User  string
}

// Fetch the main schedule list then the details about specific schedules
func (a *Api) Schedules() (Schedules, error) {
	var schdList Schedules

	res, err := a.client.ListSchedules(pagerduty.ListSchedulesOptions{})
	if err != nil {
		return schdList, err
	}

	var today string = time.Now().UTC().Format("2006-01-02")
	var nextWeek string = time.Now().UTC().Add(time.Hour * 24 * 7).Format("2006-01-02")

	for _, bareSchedule := range res.Schedules {
		schd := &Schedule{}

		schd.Id = bareSchedule.ID
		schd.Name = bareSchedule.Name
		schd.Timezone = bareSchedule.TimeZone

		res, err := a.client.GetSchedule(bareSchedule.ID, pagerduty.GetScheduleOptions{
			TimeZone: a.timezone,
			Since:    today,
			Until:    nextWeek,
		})
		if err != nil {
			return schdList, err
		}

		var activeEntries int
		for _, se := range res.FinalSchedule.RenderedScheduleEntries {
			start, err := time.Parse(time.RFC3339Nano, se.Start)
			if err != nil {
				return schdList, err
			}

			end, err := time.Parse(time.RFC3339Nano, se.End)
			if err != nil {
				return schdList, err
			}

			if start.Before(time.Now().UTC()) && end.After(time.Now().UTC()) {
				if activeEntries == 0 {
					p := CallPeriod{}
					p.Start = start
					p.User = se.User.ID
					schd.CurrentPeriod = &p
				}
				activeEntries += 1
			}

			if start.After(time.Now().UTC()) && (schd.NextPeriod == nil || start.Before(schd.NextPeriod.Start)) {
				p := CallPeriod{}
				p.Start = start
				p.User = se.User.ID
				schd.NextPeriod = &p
			}

			schdList = append(schdList, *schd)
		}

		lf := log.Fields{
			"id": schd.Id,
		}

		if schd.CurrentPeriod == nil {
			log.WithFields(lf).Warning("No active current period for schedule")
		} else {
			lf["currentCall"] = schd.CurrentPeriod.User
		}

		if schd.NextPeriod == nil {
			log.WithFields(lf).Warning("No active next period for schedule")
		} else {
			lf["nextCall"] = schd.NextPeriod.User
			lf["changeover"] = schd.NextPeriod.Start
		}

		if activeEntries > 1 {
			log.WithFields(lf).Warning("Multiple active schedules")
		}
		log.WithFields(lf).Debug("Got schedule entries")
	}

	return schdList, nil
}
