/*Package scheduler executes binaries at a certain datetime.
 */
package scheduler

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ajagnic/ARImport/src/output"
)

var runTime time.Time
var stopexec chan bool

//Config parses config.txt and initiates runTime variable.
func Config() {
	today := time.Now()

	cfg, err := os.Open("./static/cfg/config.txt")
	defer cfg.Close()

	if err == nil {
		cfgBytes := []byte{0, 0, 0, 0}
		_, e1 := cfg.Read(cfgBytes)

		cfgString := string(cfgBytes)

		runHour, e2 := strconv.Atoi(cfgString[:2])
		runMin, e3 := strconv.Atoi(cfgString[2:])
		output.Check(e1, e2, e3)

		runTime = time.Date(today.Year(), today.Month(), today.Day(), runHour, runMin, 0, 0, today.Location())
		manageDays(runHour, runMin)
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
	} else {
		// output.Pf("", fmt.Errorf("err"), true)
	}
}

func manageDays(hour, min int) {
	return
}
