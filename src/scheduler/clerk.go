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
	cfg, err := os.Open("./static/cfg/config.txt")
	defer cfg.Close()
	output.Pf("Could not open config in scheduler: %v", err, false)

	cfgBytes := []byte{0, 0, 0, 0}
	_, err = cfg.Read(cfgBytes)
	output.Pf("Could not read from config file: %v", err, false)

	cfgString := string(cfgBytes)

	runHour, err := strconv.Atoi(cfgString[:2])
	output.Pf("strconv.Atoi() - runHour: %v", err, false)

	runMin, err := strconv.Atoi(cfgString[2:])
	output.Pf("strconv.Atoi() - runMin: %v", err, false)

	today := time.Now()
	runTime = time.Date(today.Year(), today.Month(), today.Day(), runHour, runMin, 0, 0, today.Location())
	fmt.Println(runTime)

	stopexec = make(chan bool, 1)
	go start()
}

func start() {
	now := time.Now()
	fmt.Println(now)
	if now.Before(runTime) {
		durationUntil := time.Until(runTime)
		fmt.Println(durationUntil)
		exeTimer := time.AfterFunc(durationUntil, func() {
			fmt.Println("RUNNING EXEC")
			output.Log.Println("RUNNING EXEC")
		})
		select {
		case <-stopexec:
			output.Log.Println("STOPPED EXEC")
			exeTimer.Stop()
		}
	} else {
		fmt.Println("not before")
		fmt.Println(now.Location())
		os.Exit(1)
	}
}

//EventListener waits for either re-init or kill events and calls necessary functions.
func EventListener(reinit, kill chan bool) {
	for {
		select {
		case <-kill:
			fmt.Println("stopping scheduler")
			stopexec <- true
			break
		case <-reinit:
			stopexec <- true
			Config()
		}
	}
}
