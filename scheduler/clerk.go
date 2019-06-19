package scheduler

import (
	"time"

	"github.com/ajagnic/ARImport/output"
)

var runTime time.Time
var hourTimer *time.Timer

func Config() {
	today := time.Now()
	loc := time.FixedZone("UTC-7", 0)
	//Run every month, every day at 12:00:00:00
	runTime = time.Date(today.Year(), today.Month(), today.Day(), 2, 4, 0, 0, loc)
}

func Run() {
	//initaite timer ever hour to compare current time
	hourTimer = time.NewTimer(time.Minute)

	for {
		select {
		case tick := <-hourTimer.C:
			//hour event, check if equal to runTime
			output.Log.Printf("Timer: %v", tick)
			if compare(tick) {
				output.Log.Print("RUNNING EXEC")
			} else {
				continue
			}
		}
	}
}

func compare(tick time.Time) bool {
	if tick.Hour() == runTime.Hour() {
		if tick.Minute() == runTime.Minute() { //possibly check if +/- 1, 2, 3 minutes
			return true
		}
		//within hour, reset timer to every minute
		hourTimer.Stop()
		hourTimer.Reset(time.Minute)
	}
	return false
}

// var runTime time.Time
// var currentTime time.Time

// func init() {
// 	loc := time.FixedZone("UTC-7", 0)
// 	runTime = time.Date(2019, time.June, 18, 24, 0, 0, 0, loc) // set time to run (read from config)
// }

// func Run() {
// 	// contains channel receiving Time every 3 sec
// 	ticker := time.NewTicker(3 * time.Second)
// 	defer ticker.Stop()
// 	stop := make(chan bool) // kill pill chan

// 	//func to delay kill pill
// 	go func() {
// 		time.Sleep(12 * time.Second)
// 		stop <- true
// 	}()

// 	// while w/ select listening to chans for Time or kill pill
// 	// on Time, calls function
// 	for {
// 		select {
// 		case <-stop:
// 			fmt.Println("Stopped")
// 			return
// 		case <-ticker.C:
// 			durationUntil()
// 		}
// 	}
// }

// func durationUntil() {
// 	durationTilExec := time.Until(runTime)
// 	fmt.Println(durationTilExec)
// }
