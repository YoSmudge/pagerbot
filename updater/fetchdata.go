package updater

import(
  "time"
  "sync"
  log "github.com/Sirupsen/logrus"
)

func (u *Updater) fetchData(){
  if u.LastFetch.After(time.Now().UTC().Add(time.Minute*-15)){
    return
  }
  log.Debug("Fetching data")

  w := sync.WaitGroup{}
  w.Add(2)
  go u.updateUsers(&w)
  go u.updateSchedules(&w)

  w.Wait()
  log.WithFields(log.Fields{
    "users": len(u.Users.users),
    "schedules": len(u.Schedules.schedules),
  }).Debug("Update done")
  u.LastFetch = time.Now().UTC()
}
