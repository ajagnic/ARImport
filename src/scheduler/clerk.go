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

//Config parses config.txt and initiates runTime variable.
func Config() (addr string) {
	today := time.Now()

	cfg, err := output.ReadJSON()
	addr = ""
	if err == nil {
		rt := cfg.RunTime

		hour, _ := strconv.Atoi(rt[:2])
		min, _ := strconv.Atoi(rt[2:])
		addr = cfg.Addr

		runTime = time.Date(today.Year(), today.Month(), today.Day(), hour, min, 0, 0, today.Location())
		//If runTime is set to early morning hours on previous day, runTime.Day will incorrect. Add 24 hours.
		if hour > 0 && hour < 7 {
			runTime.Add(24 * time.Hour)
		}
	} else {
		//Config not read, default to 11:45pm.
		output.Pf("Config: %v", err, false)
		runTime = time.Date(today.Year(), today.Month(), today.Day(), 23, 45, 0, 0, today.Location())
	}

	fmt.Println(runTime)
	stopexec = make(chan bool, 1)
	go start()

	return
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
			cfg, _ := output.ReadJSON()
			cfg.LastRun = time.Now()
			_ = output.WriteJSON(cfg)
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
