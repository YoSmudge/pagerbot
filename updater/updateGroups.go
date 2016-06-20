package updater

import(
  "fmt"
  "time"
  "sort"
  "strings"
  "reflect"
  log "github.com/Sirupsen/logrus"
  "github.com/yosmudge/pagerbot/config"
  "github.com/yosmudge/pagerbot/pagerduty"
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

      var activePeriod *pagerduty.CallPeriod

      if schd.NextPeriod != nil{
        if changeover.IsZero() || schd.NextPeriod.Start.Before(changeover){
          changeover = schd.NextPeriod.Start
        }
      }

      if !changeover.IsZero() && time.Now().UTC().After(changeover){
        activePeriod = schd.NextPeriod
      } else if schd.CurrentPeriod != nil{
        activePeriod = schd.CurrentPeriod
      }

      if activePeriod != nil{
        lf["userId"] = activePeriod.User
        usr := u.Users.ById(activePeriod.User)
        if usr == nil{
          log.WithFields(lf).Warning("Could not find user with ID")
          continue
        }
        currentUsers = append(currentUsers, usr)
      }
    }

    lf["scheduleId"] = nil
    lf["userId"] = nil
    lf["changeover"] = changeover

    var pdUsers []string
    var slackUsers []string
    var userNames []string

    for _,u := range currentUsers{
      pdUsers = append(pdUsers, u.PagerdutyId)
      slackUsers = append(slackUsers, u.SlackId)
      userNames = append(userNames, fmt.Sprintf("@%s", u.SlackName))
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

      var userList string
      if len(userNames) > 1{
        userList = strings.Join(userNames[:len(userNames)-1], ", ")
      }

      if len(userNames) > 1{
        userList = fmt.Sprintf("%s & %s", userList, userNames[len(userNames)-1])
      } else {
        userList = userNames[0]
      }

      msgText := fmt.Sprintf(group.UpdateMessage.Message, userList)
      for _,c := range group.UpdateMessage.Channels{
        u.Slack.PostMessage(c, msgText)
      }
    }
  }
}
