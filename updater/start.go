package updater

import (
	"time"
)

// Start the updater process
func (u *Updater) Start() {
	u.Wg.Add(1)
	go u.run()
}

// Loop for updater
// Will call for new data then call the update function
// Runs on each `updateEvery` interval
const updateEvery time.Duration = time.Minute * 5

func (u *Updater) run() {
	defer u.Wg.Done()

	for {
		u.fetchData()
		u.updateGroups()

		nextInterval := time.Unix((time.Now().UTC().Unix()/int64(updateEvery.Seconds())+1)*int64(updateEvery.Seconds()), 0)
		waitTime := nextInterval.Sub(time.Now().UTC())
		time.Sleep(waitTime)
	}
}
