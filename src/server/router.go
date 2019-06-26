/*Package server contains the http server and routing for a web interface.
 */
package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ajagnic/ARImport/src/output"
)

var srv http.Server
var reinit chan bool

//Run registers the handler functions and blocks until interrupt.
func Run(addr string, reinitC chan bool) {
	reinit = reinitC
	srv = http.Server{
		Addr:     addr,
		ErrorLog: output.Log,
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/static/", contentHandler)
	http.HandleFunc("/store", postHandler)
	http.HandleFunc("/config", configHandler)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		fmt.Printf("Starting server on URL/Port: %v\n", addr)
		err := srv.ListenAndServe()
		output.Pf("", err, true) //Exit
	}()

	<-sigint
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
		cfg, e2 := output.ReadConfig()

		addr := r.FormValue("addr")
		if addr != "" {
			cfg["Addr"] = addr
		}
		runTime := r.FormValue("runtime")
		if runTime != "" {
			cfg["RunTime"] = runTime
		}

		e3 := output.WriteConfig(cfg)
		output.Check(e1, e2, e3)

		reinit <- true
	}
	http.ServeFile(w, r, "./static/config.html")
}

//postHandler streams up form data from POST and saves to local file.
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rdr, err := r.MultipartReader()
		output.Pf("postHandler - MultipartReader: %v", err, false)

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

//Shutdown calls http.Server.Shutdown with a 10ms timeout.
func Shutdown() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	err = srv.Shutdown(ctx)
	output.Pf("", err, true)
	fmt.Println("Server gracefully shutdown.")
	return
}
