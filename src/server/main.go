package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/ajagnic/ARImport/src/output"
	"github.com/ajagnic/ARImport/src/scheduler"
)

var reinit chan bool

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

//configHandler parses config form on POST, saves to local file, then re-initializes scheduler.
func configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		runTime := r.FormValue("runtime")

		file, err := os.OpenFile("./static/cfg/config.txt", os.O_WRONLY, 0644)
		output.Pf("Could not open config file: %v", err, false)

		if err == nil {
			file.WriteString(runTime)
			file.Sync()
			file.Close()
			reinit <- true
		}

	}
	http.ServeFile(w, r, "./static/config.html")
}

//contentHandler serves static files based on URL path.
func contentHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+r.URL.Path)
}

//postHandler streams up form data from POST and saves to local file.
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rdr, err := r.MultipartReader()
		output.Pf("File reader: %v", err, false)

		for { //Loop file parts until EOF and place in write buffer.
			part, err := rdr.NextPart()
			if err == io.EOF {
				break
			} else {
				output.Pf("File parts: %v", err, false)
			}

			if part.FileName() == "" {
				continue
			}

			file, err := os.OpenFile("./static/csv/"+part.FileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			output.Pf("", err, true)
			defer file.Close()

			_, err = io.CopyBuffer(file, part, nil)
			output.Pf("CopyBuffer: %v", err, false)
		}

		http.ServeFile(w, r, "./static/index.html")
	} else {
		http.ServeFile(w, r, "./static/error.html")
		output.Log.Printf("Invalid method in /store: %v", r.Method)
	}
}

func main() {
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~SCHEDULER
	scheduler.Config()
	reinit = make(chan bool)
	killsched := make(chan bool, 1)
	go scheduler.EventListener(reinit, killsched)

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~SERVER
	addr := ":8001"
	srv := http.Server{
		Addr:     addr,
		ErrorLog: output.Log,
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/static/", contentHandler)
	http.HandleFunc("/store", postHandler)
	http.HandleFunc("/config", configHandler)

	//sigint listens for interrupt signal then cleans up main()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		fmt.Printf("Starting server on URL/Port: %v\n", addr)
		err := srv.ListenAndServe()
		output.Pf("", err, true)
	}()

	<-sigint

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~CLEANUP
	killsched <- true
	err := srv.Shutdown(context.Background())
	output.Pf("", err, true)
	output.Close()
	fmt.Println("Server gracefully shutdown.")
}
