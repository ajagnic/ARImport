/*Package output contains interface for logging errors to text file and parsing JSON.
 */
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

var config = map[string]string{
	"Addr":    "127.0.0.1:8001",
	"RunTime": "2345",
	"LastRun": "0001-01-01T00:00:00Z",
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

//ReadJSON parses config.txt and returns config pointer.
func ReadJSON() (cfg *map[string]string, err error) {
	cfg = &config

	file, err := os.OpenFile("./static/cfg/config.txt", os.O_RDWR, 0644)
	if err != nil {
		Pf("ReadJSON - Opening file: %v", err, false)
		return
	}
	defer file.Close()

	size := 54
	fileBytes := make([]byte, size) //BUG(r): something happening with the amount of writing/reading. Extra dupped chars end up at end.
	b, err := file.Read(fileBytes)
	if err != nil && err != io.EOF {
		Pf("ReadJSON - Reading file: %v", err, false)
	} else if b == 0 { //Config file empty, write default.
		WriteJSON(&config)
	}
	fmt.Println(b)

	err = json.Unmarshal(fileBytes, cfg)
	if err != nil {
		Pf("ReadJSON - Unmarshal: %v", err, false)
	}

	return
}

//WriteJSON serializes config map to file.
func WriteJSON(cfg *map[string]string) (err error) {
	file, err := os.OpenFile("./static/cfg/config.txt", os.O_WRONLY, 0644)
	if err != nil {
		Pf("WriteJSON - Opening file: %v", err, false)
		return
	}
	defer file.Close()

	bytes, err := json.Marshal(cfg)
	file.Write(bytes)

	return
}
