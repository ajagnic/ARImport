/*Package output contains interface for logging errors to text file and parsing JSON.
 */
package output

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//Log is a pointer to the log.Logger struct.
var Log *log.Logger

var file *os.File
var config = map[string]string{
	"Addr":    "127.0.0.1:8001",
	"RunTime": "2345",
	"LastRun": "0001-01-01T00:00:00Z",
}

func init() {
	file, err := os.OpenFile("./static/cfg/output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("!!! Log File Error: %v !!!", err)
	}

	Log = log.New(file, "Error: ", log.LstdFlags)
}

//Check can take multiple errors and log them.
func Check(errs ...error) {
	for _, err := range errs {
		if err != nil {
			Log.Println(err)
		}
	}
}

//Pf is a wrapper around Logger.Printf and Logger.Fatal.
func Pf(fS string, err error, fatal bool) {
	if err != nil {
		if fatal {
			Log.Fatal(err)
		} else {
			Log.Printf(fS, err)
		}
	}
}

//Close flushes data and releases log file resource.
func Close() {
	file.Sync()
	file.Close()
}

//ReadConfig parses config.txt and returns config map.
func ReadConfig() (map[string]string, error) {
	cfg := &config

	cfgBytes, err := ioutil.ReadFile("./static/cfg/config.txt")
	if err != nil { //Could not read config, return default.
		Pf("ReadConfig - ReadFile: %v", err, false)
		return config, err
	}

	err = json.Unmarshal(cfgBytes, cfg)
	if err != nil {
		Pf("ReadConfig - Unmarshal: %v", err, false)
	}

	return config, err
}

//WriteConfig serializes config map to file.
func WriteConfig(cfg map[string]string) (err error) {
	cfgP := &cfg
	file, err := os.OpenFile("./static/cfg/config.txt", os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		Pf("WriteConfig - OpenFile: %v", err, false)
		return
	}

	bytes, err := json.Marshal(cfgP)
	file.Write(bytes)

	return
}
