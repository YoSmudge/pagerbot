package pagerduty

import(
  "fmt"
  "time"
  log "github.com/Sirupsen/logrus"
)

type Schedules []Schedule

type Schedule struct{
  Id              string
  Name            string
  Timezone        string    `json:"time_zone"`
  CurrentPeriod   *CallPeriod
  NextPeriod      *CallPeriod
}

type CallPeriod struct{
  Start     time.Time
  User      string
}

// Fetch the main schedule list then the details about specific schedules
func (a *Api) Schedules() (Schedules, error){
  var schdList Schedules

  err := a.requestThing("schedules", &schdList)
  if err != nil {
    return schdList, err
  }

  var today string = time.Now().UTC().Format("2006-01-02")
  var nextWeek string = time.Now().UTC().Add(time.Hour*24*7).Format("2006-01-02")

  for i,_ := range schdList{
    schd := &schdList[i]
    rsp, err := a.request(fmt.Sprintf("schedules/%s/entries?since=%s&until=%s&time_zone=%s&overflow=true", schd.Id, today, nextWeek, a.timezone))
    if err != nil {
      return schdList, err
    }

    var schdInfo struct{
      Entries     []struct{
        User        struct{
          Id          string
        }
        Start       time.Time
        End         time.Time
      }
    }

    rsp.Unmarshal(&schdInfo)

    var activeEntries int
    for _,se := range schdInfo.Entries{
      if se.Start.Before(time.Now().UTC()) && se.End.After(time.Now().UTC()){
        if activeEntries == 0{
          p := CallPeriod{}
          p.Start = se.Start
          p.User = se.User.Id
          schd.CurrentPeriod = &p
        }
        activeEntries += 1
      }

      if se.Start.After(time.Now().UTC()) && (schd.NextPeriod == nil || se.Start.Before(schd.NextPeriod.Start)){
        p := CallPeriod{}
        p.Start = se.Start
        p.User = se.User.Id
        schd.NextPeriod = &p
      }
    }

    lf := log.Fields{
      "id": schd.Id,
    }

    if schd.CurrentPeriod == nil{
      log.WithFields(lf).Warning("No active current period for schedule")
    } else {
      lf["currentCall"] = schd.CurrentPeriod.User
    }

    if schd.NextPeriod == nil{
      log.WithFields(lf).Warning("No active next period for schedule")
    } else {
      lf["nextCall"] = schd.NextPeriod.User
      lf["changeover"] = schd.NextPeriod.Start
    }

    if activeEntries > 1{
      log.WithFields(lf).Warning("Multiple active schedules")
    }
    log.WithFields(lf).Debug("Got schedule entries")
  }

  return schdList, nil
}
