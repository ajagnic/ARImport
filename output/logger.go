/*Package 'output' contains interface for logging errors to text file.
 */
package output

import (
	"fmt"
	"log"
	"os"
)

var Log *log.Logger

var file *os.File

func init() {
	file, err := os.OpenFile("./output/output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("!!! Log File Error: %v !!!", err)
	}

	Log = log.New(file, "Error: ", log.LstdFlags)
}

//Close flushes data and releases file resource.
func Close() {
	file.Sync()
	file.Close()
}
