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
func Config() {
	today := time.Now()

	cfg, err := output.ReadJSON()
	if err == nil {
		rt := cfg.RunTime
		hour, _ := strconv.Atoi(rt[:2])
		min, _ := strconv.Atoi(rt[2:])
		runTime = time.Date(today.Year(), today.Month(), today.Day(), hour, min, 0, 0, today.Location())
		manageDays(hour, min)
	} else {
		//Config not read, default to 11:45pm.
		output.Log.Printf("Config: %v", err)
		runTime = time.Date(today.Year(), today.Month(), today.Day(), 23, 45, 0, 0, today.Location())
	}

	fmt.Println(runTime)
	stopexec = make(chan bool, 1)
	go start()
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
		exeTimer := time.AfterFunc(durationUntil, func() {
			output.Log.Println("RUNNING EXEC")
		})

		select {
		case <-stopexec:
			output.Log.Println("STOPPED EXEC")
			exeTimer.Stop()
		}

		cfg, _ := output.ReadJSON()
		cfg.LastRun = time.Now()
		_ = output.WriteJSON(cfg)
	} else {
		// output.Pf("", fmt.Errorf("err"), true)
	}
}

func manageDays(hour, min int) {
	return
}
