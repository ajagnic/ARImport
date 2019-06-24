/*Package scheduler executes binaries at a certain datetime.
 */
package scheduler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ajagnic/ARImport/src/output"
)

var runTime time.Time
var stopexec chan bool

//Config parses config.txt and initiates runTime variable. Returns address for the server.
func Config() (addr string) {
	today := time.Now()

	cfgP, err := output.ReadConfig()
	cfg := *cfgP
	if err != nil {
		output.Pf("Config: %v", err, false)
	}

	rt := cfg["RunTime"]

	hour, _ := strconv.Atoi(rt[:2])
	min, _ := strconv.Atoi(rt[2:])

	runTime = time.Date(today.Year(), today.Month(), today.Day(), hour, min, 0, 0, today.Location())
	//If runTime is set to early morning hours on previous day, runTime.Day will incorrect. Add 24 hours.
	if hour > 0 && hour < 7 {
		runTime.Add(24 * time.Hour)
	}

	fmt.Println(runTime)
	stopexec = make(chan bool, 1)
	go start()

	return cfg["Addr"]
}

//EventListener waits for either re-init or kill events and calls necessary functions.
func EventListener(reinit, kill chan bool) {
	for {
		select {
		case <-kill:
			fmt.Println("Stopping scheduler.")
			stopexec <- true
			break
		case <-reinit:
			stopexec <- true
			Config()
		}
	}
}

func start() {
	now := time.Now()

	if now.Before(runTime) {
		durationUntil := time.Until(runTime)
		fmt.Println(durationUntil)
		//Runs timer in routine that will execute func after the duration.
		exeTimer := time.AfterFunc(durationUntil, func() {
			output.Log.Println("RUNNING EXEC")

			cfgP, err := output.ReadConfig()
			output.Pf("start - ReadConfig: %v", err, false)

			cfg := *cfgP
			now = time.Now()
			cfg["LastRun"] = now.Format("ANSIC") //BUG(r): this doesnt work.

			err = output.WriteConfig(cfgP)
			output.Pf("start - WriteConfig: %v", err, false)
		})
		//Listen for cancel event. (blocking call)
		select {
		case <-stopexec:
			output.Log.Println("STOPPED EXEC")
			exeTimer.Stop()
		}
	} else {
		//time.Now after runTime:
		//Either have a buffer of 1-2 hours where exec still runs, or
		//Wait 24 hours.
	}
}
