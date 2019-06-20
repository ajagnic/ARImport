package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/ajagnic/ARImport/output"
	"github.com/ajagnic/ARImport/scheduler"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

//contentHandler serves static files based on URL path.
func contentHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+r.URL.Path)
}

//postHandler streams up form data from POST and saves to local file.
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rdr, _ := r.MultipartReader()
		// output.Pf("File reader: %v", err, false)

		for { //Loop file parts until EOF and place in write buffer.
			part, err := rdr.NextPart()
			if err == io.EOF {
				break
			} else {
				// output.Pf("File parts: %v", err, false)
			}

			if part.FileName() == "" {
				continue
			}

			file, err := os.OpenFile("./static/csv/"+part.FileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			// output.Pf("", err, true)
			defer file.Close()

			_, err = io.CopyBuffer(file, part, nil)
			// output.Pf("CopyBuffer: %v", err, false)
		}

		http.ServeFile(w, r, "./static/index.html") //TODO: implement this route w/ custom handler (middleware), test if process continues while http serves new
	} else {
		http.ServeFile(w, r, "./static/error.html")
		output.Log.Printf("Invalid method in /store: %v", r.Method)
	}
}

func main() {
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~SCHEDULER
	timekill := make(chan bool, 1)
	scheduler.Config()
	go scheduler.Run(timekill)

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~SERVER
	addr := ":8001"
	srv := http.Server{
		Addr:     addr,
		ErrorLog: output.Log,
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/static/", contentHandler)
	http.HandleFunc("/store", postHandler)

	//sigint listens for interrupt signal and calls http.Server.Shutdown.
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		fmt.Printf("Starting server on URL/Port: %v\n", addr)
		_ = srv.ListenAndServe()
		// output.Pf("", err, true)
	}()

	<-sigint
	_ = srv.Shutdown(context.Background())
	// output.Pf("", err, true)

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	timekill <- true
	output.Close()
	fmt.Println("Server gracefully shutdown.")
}
