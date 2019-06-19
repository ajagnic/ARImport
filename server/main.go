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

		http.ServeFile(w, r, "./static/index.html") //TODO: implement this route w/ custom handler (middleware), test if process continues while http serves new
	} else {
		http.ServeFile(w, r, "./static/error.html")
		output.Log.Printf("Invalid method in /store: %v", r.Method)
	}
}

func main() {
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~SCHEDULER
	scheduler.Config()
	go scheduler.Run()

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~SERVER
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/static/", contentHandler)
	http.HandleFunc("/store", postHandler)

	addr := ":8001"
	srv := http.Server{Addr: addr, ErrorLog: output.Log}

	serverKill := make(chan os.Signal)
	signal.Notify(serverKill, os.Interrupt)

	go func() {
		fmt.Printf("Starting server on URL/Port: %v\n", addr)
		err := srv.ListenAndServe()
		output.Pf("", err, true)
	}()

	<-serverKill
	err := srv.Shutdown(context.Background())
	output.Pf("", err, true)
	fmt.Println("Server gracefully shutdown.")

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	output.Close()
}
