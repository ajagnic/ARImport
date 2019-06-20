package scheduler

import (
	"fmt"
	"time"

	"github.com/ajagnic/ARImport/output"
)

var runTime time.Time

func Config() {
	today := time.Now()
	loc := time.FixedZone("UTC-7", 0)
	//Run every month, every day at specified time.
	runTime = time.Date(today.Year(), today.Month(), today.Day(), 9, 42, 0, 0, loc)
	output.Log.Printf("RunTime: %v", runTime)
}

func Run(timekill chan bool) {
	hourTicker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-timekill:
			hourTicker.Stop()
			fmt.Println("stopped")
			break
		case bzz := <-hourTicker.C:
			fmt.Println(bzz)
			output.Log.Println(bzz)
			if compare(bzz) {
				hourTicker.Stop()
				minTicker := time.NewTicker(time.Second)
				minkill := make(chan bool, 1)
				go func() {
					time.Sleep(5 * time.Second)
					minkill <- true
				}()
				for {
					select {
					case <-minkill:
						minTicker.Stop()
						break
					case t := <-minTicker.C:
						fmt.Println(t)
						output.Log.Println(t)
					}
				}
			}
		}
	}
}

func compare(buzz time.Time) bool {
	buzzMins := buzz.Minute()
	runTimeMins := runTime.Minute()
	if buzzMins == runTimeMins || buzzMins == runTimeMins+1 || buzzMins == runTimeMins-1 {
		return true
	}
	return false
}

// if buzz.Hour() == runTime.Hour() {
// 	if compare(buzz) {
// 		output.Log.Printf("RUNNNG EXEC @ %v", buzz)
// 	} else { //runTime within the hour.
// 		hourTimer.Stop()
// 		minTicker := time.NewTicker(time.Second)
// 		//~~~~~~~~~~~~~~~~~~~~~~~
// 		for {
// 			buzz = <-minTicker.C
// 			output.Log.Printf("Timer: %v", buzz)
// 			if compare(buzz) {
// 				output.Log.Printf("RUNNNG EXEC @ %v", buzz)
// 				minTicker.Stop()
// 				break
// 			}
// 		}
// 		//~~~~~~~~~~~~~~~~~~~~~~~
// 		hourTimer.Reset(time.Minute)
// 	}
// }
