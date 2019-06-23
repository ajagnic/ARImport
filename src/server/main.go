package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

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
		e1 := r.ParseForm()
		runTime := r.FormValue("runtime")

		file, e2 := os.OpenFile("./static/cfg/config.txt", os.O_WRONLY, 0644)
		output.Check(e1, e2)

		if e2 == nil && runTime != "" {
			_, e3 := file.WriteString(runTime)
			e4 := file.Sync()
			e5 := file.Close()
			output.Check(e3, e4, e5)

			reinit <- true
		} else {
			http.Error(w, "Error with Config", 500)
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	err := srv.Shutdown(ctx)
	output.Pf("", err, true)
	output.Close()
	fmt.Println("Server gracefully shutdown.")
}
