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

//Config parses /static/cfg/config.txt and initiates runTime variable.
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
	loc := time.FixedZone("UTC-7", 0)
	runTime = time.Date(today.Year(), today.Month(), today.Day(), runHour, runMin, 0, 0, loc)
	fmt.Println(runTime)
}

//EventListener waits for either re-init or kill events and calls necessary functions.
func EventListener(reinit, kill chan bool) {
	for {
		select {
		case <-kill:
			fmt.Println("stopping scheduler")
			break
		case <-reinit:
			Config()
		}
	}
}
