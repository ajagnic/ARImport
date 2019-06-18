package scheduler

import (
	"fmt"
	"time"
)

var runTime time.Time
var currentTime time.Time

func init() {
	loc := time.FixedZone("UTC-7", 0)
	runTime = time.Date(2019, time.June, 18, 24, 0, 0, 0, loc) // set time to run (read from config)
}

func Run() {
	// contains channel receiving Time every 3 sec
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	stop := make(chan bool) // kill pill chan

	//func to delay kill pill
	go func() {
		time.Sleep(12 * time.Second)
		stop <- true
	}()

	// while w/ select listening to chans for Time or kill pill
	// on Time, calls function
	for {
		select {
		case <-stop:
			fmt.Println("Stopped")
			return
		case <-ticker.C:
			durationUntil()
		}
	}
}

func durationUntil() {
	durationTilExec := time.Until(runTime)
	fmt.Println(durationTilExec)
}
