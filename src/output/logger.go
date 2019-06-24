/*Package output contains interface for logging errors to text file and parsing JSON.
 */
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type config struct { //TODO: possibly change to a map[string]string and parse datetime string.
	Addr    string
	RunTime string
	LastRun time.Time
}

var Log *log.Logger

var file *os.File

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

//ReadJSON parses config.txt and returns config struct.
func ReadJSON() (cfg config, err error) {
	file, err := os.OpenFile("./static/cfg/config.txt", os.O_RDWR, 0644)
	if err != nil {
		Pf("ReadJSON - Opening file: %v", err, false)
		return cfg, err
	}
	defer file.Close()

	fileBytes := make([]byte, 66)
	b, err := file.Read(fileBytes)
	if err != nil && err != io.EOF {
		Pf("ReadJSON - Reading file: %v", err, false)
	}
	fmt.Println(b)

	err = json.Unmarshal(fileBytes, &cfg)
	if err != nil {
		Pf("ReadJSON - Unmarshal: %v", err, false)
	}

	return
}

//WriteJSON serializes config struct to file.
func WriteJSON(cfg config) (err error) {
	file, err := os.OpenFile("./static/cfg/config.txt", os.O_WRONLY, 0644)
	if err != nil {
		Pf("WriteJSON - Opening file: %v", err, false)
		return err
	}
	defer file.Close()

	bytes, err := json.Marshal(&cfg)
	file.Write(bytes)

	return
}
