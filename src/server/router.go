package server

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ajagnic/ARImport/src/output"
)

var reinit chan bool

func init() {
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

//contentHandler serves static files based on URL path.
func contentHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+r.URL.Path)
}

//configHandler parses config form on POST, saves to local file, then re-initializes scheduler.
func configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		e1 := r.ParseForm()
		fmt.Println("parse")

		cfg, e2 := output.ReadConfig()
		fmt.Println("read", cfg)

		addr := r.FormValue("addr")
		if addr != "" {
			cfg["Addr"] = addr
		}
		fmt.Println("addr", addr)
		runTime := r.FormValue("runtime")
		if runTime != "" {
			cfg["RunTime"] = runTime
		}
		fmt.Println("rt", runTime)

		e3 := output.WriteConfig(cfg)
		fmt.Println("write")
		output.Check(e1, e2, e3)

		reinit <- true
		fmt.Println("reinit")
	}
	http.ServeFile(w, r, "./static/config.html")
}

//postHandler streams up form data from POST and saves to local file.
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rdr, err := r.MultipartReader()
		output.Pf("MultipartReader: %v", err, false)

		if rdr != nil {
			for { //Loop file parts until EOF and place in write buffer.
				part, e1 := rdr.NextPart()
				if e1 == io.EOF {
					break
				}

				if part.FileName() == "" {
					continue
				}

				file, e2 := os.OpenFile("./static/csv/"+part.FileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				defer file.Close()

				_, e3 := io.CopyBuffer(file, part, nil)
				output.Check(e1, e2, e3)
			}
		}

		http.ServeFile(w, r, "./static/index.html")
	} else {
		http.Error(w, "Invalid Method", 405)
	}
}
