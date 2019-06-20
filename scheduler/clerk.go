package scheduler

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ajagnic/ARImport/output"
)

var runTime time.Time

//Config parses config.txt and initiates runTime var.
func Config() {
	cfg, err := os.Open("./scheduler/config.txt")
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
