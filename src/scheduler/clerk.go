/*Package scheduler executes binaries at a scheduled datetime.
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
func Config() string {
	today := time.Now()

	cfg, e1 := output.ReadConfig()

	rt := cfg["RunTime"] //format: '0000'
	hour, e2 := strconv.Atoi(rt[:2])
	min, e3 := strconv.Atoi(rt[2:])

	runTime = time.Date(today.Year(), today.Month(), today.Day(), hour, min, 0, 0, today.Location())
	//If runTime is set to early morning hours on previous day, runTime.Day will be incorrect. Add 24 hours.
	if hour >= 0 && hour < 7 {
		runTime = runTime.Add(24 * time.Hour)
	}

	fmt.Println("Scheduled Runtime: ", runTime)
	fmt.Println("Last Run: ", cfg["LastRun"])

	stopexec = make(chan bool, 1)
	go start()

	output.Check(e1, e2, e3)
	return cfg["Addr"]
}

//EventListener waits for either re-init or kill events and calls necessary functions.
func EventListener(reinit, kill chan bool) {
	for {
		select {
		case <-kill:
			fmt.Println("\nStopping scheduler...")
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
			//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
			output.Log.Println("RUNNING EXEC")
			//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
			cfg, err := output.ReadConfig()
			output.Pf("start - ReadConfig: %v", err, false)

			now = time.Now()
			cfg["LastRun"] = now.Format(time.ANSIC)

			err = output.WriteConfig(cfg)
			output.Pf("start - WriteConfig: %v", err, false)
		})
		//Listen for cancel event. (blocking call)
		select {
		case <-stopexec:
			output.Log.Println("STOPPED EXEC")
			exeTimer.Stop()
		}
	} else {
		//Wait 23 hours and restart.
		waitTimer := time.AfterFunc(23*time.Hour, func() {
			Config()
		})
		select {
		case <-stopexec:
			output.Log.Println("STOPPED WAITING")
			waitTimer.Stop()
		}
	}
}

// func runBin() {
// 	cmd := exec.Command("./exe/RunExternally.exe")
// 	out, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(string(out))
// }
