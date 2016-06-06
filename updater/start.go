package updater

import(
  "time"
)

// Start the updater process
func (u *Updater) Start(){
  u.Wg.Add(1)
  go u.run()
}

// Loop for updater
// Will call for new data then call the update function, both rate limit themselves so this just runs them on every loop
func (u *Updater) run(){
  defer u.Wg.Done()
  t := time.Tick(time.Second*15)

  for {
    u.fetchData()
    u.updateGroups()
    <-t
  }
}
