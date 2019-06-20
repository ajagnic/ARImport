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
var reInit chan bool

//Config parses config.txt and initiates runTime var.
func Config() {
	cfg, err := os.Open("./static/cfg/config.txt")
	defer cfg.Close()
	output.Pf("Could not open config in scheduler: %v", err, false)

	runTimeBytes := []byte{0, 0, 0, 0}
	_, err = cfg.Read(runTimeBytes)
	output.Pf("Could not read from config file: %v", err, false)

	runTimeString := string(runTimeBytes)

	runHour, err := strconv.Atoi(runTimeString[:2])
	output.Pf("strconv.Atoi(): %v", err, false)

	runMin, err := strconv.Atoi(runTimeString[2:])
	output.Pf("strconv.Atoi(): %v", err, false)

	today := time.Now()
	loc := time.FixedZone("UTC-7", 0)
	runTime = time.Date(today.Year(), today.Month(), today.Day(), runHour, runMin, 0, 0, loc)
	fmt.Println(runTime)
}

//ReInit listens for new config and re-starts any running scheduler processes.
func ReInit(reInitChan, kill chan bool) {
	reInit = reInitChan
	for {
		select {
		case <-kill:
			return
		case <-reInit:
			Config()
		}
	}
}
