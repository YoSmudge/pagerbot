package updater

import(
  "time"
  "sort"
  "reflect"
  log "github.com/Sirupsen/logrus"
  "github.com/yosmudge/pagerbot/config"
)

// Ensure all the slack groups are up to date
func (u *Updater) updateGroups(){
  for _,group := range config.Config.Groups{
    lf := log.Fields{
      "group": group.Name,
    }

    var currentUsers []*User
    var changeover time.Time
    for _,s := range group.Schedules{
      lf["scheduleId"] = s
      schd := u.Schedules.ById(s)
      if schd == nil{
        log.WithFields(lf).Warning("Could not find schedule with ID")
        continue
      }

      if schd.CurrentPeriod != nil{
        lf["userId"] = schd.CurrentPeriod.User
        usr := u.Users.ById(schd.CurrentPeriod.User)
        if usr == nil{
          log.WithFields(lf).Warning("Could not find user with ID")
          continue
        }
        currentUsers = append(currentUsers, usr)

        if schd.NextPeriod != nil{
          if changeover.IsZero() || schd.NextPeriod.Start.Before(changeover){
            changeover = schd.NextPeriod.Start
          }
        }
      }
    }

    lf["scheduleId"] = nil
    lf["userId"] = nil
    lf["changeover"] = changeover

    var pdUsers []string
    var slackUsers []string

    for _,u := range currentUsers{
      pdUsers = append(pdUsers, u.PagerdutyId)
      slackUsers = append(slackUsers, u.SlackId)
    }

    lf["pdUsers"] = pdUsers
    lf["slackUsers"] = slackUsers

    currentMembers, err := u.Slack.GroupMembers(group.Name)
    if err != nil {
      lf["err"] = err
      log.WithFields(lf).Warning("Could not get Slack group members")
      continue
    }

    lf["currentMembers"] = currentMembers
    log.WithFields(lf).Debug("Group status")
    sort.Strings(currentMembers)
    sort.Strings(slackUsers)
    if !reflect.DeepEqual(currentMembers, slackUsers){
      err := u.Slack.UpdateMembers(group.Name, slackUsers)
      if err != nil {
        lf["err"] = err
        log.WithFields(lf).Warning("Could not update Slack group members")
        continue
      }
      log.WithFields(lf).Info("Updating group members")
    }
  }
}
