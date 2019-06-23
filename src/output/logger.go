/*Package output contains interface for logging errors to text file.
 */
package output

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type config struct {
	Addr    string
	RunTime string
	LastRun string
}

// Log is a pointer to the log.Logger struct.
var Log *log.Logger

var file *os.File
var cfg config

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

//ReadJSON is a test.
func ReadJSON() {
	var cfgg config
	file, _ := os.OpenFile("./static/cfg/config.txt", os.O_RDWR, 0644)
	fileBytes := make([]byte, 49)
	b, _ := file.Read(fileBytes)
	fmt.Println(b)
	err := json.Unmarshal(fileBytes, &cfgg)
	fmt.Println(cfg.Addr, err)
}

func WriteJSON() {
	cfg = config{"8001", "1234", "1234"}
	file, _ := os.OpenFile("./static/cfg/config.txt", os.O_WRONLY, 0644)
	bytes, _ := json.Marshal(cfg)
	file.Write(bytes)
	file.Close()
}
