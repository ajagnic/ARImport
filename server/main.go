package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ajagnic/ARImport/exe"
	"github.com/ajagnic/ARImport/output"
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

		http.ServeFile(w, r, "./static/index.html")
	} else {
		http.ServeFile(w, r, "./static/error.html")
		output.Log.Printf("Invalid method in /store: %v", r.Method)
	}
}

func main() {
	//TODO: thread off exe src
	exe.RunBin()
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~SERVER
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/static/", contentHandler)
	http.HandleFunc("/store", postHandler)

	addr := ":8001"

	fmt.Printf("Starting server on URL/Port: %v\n", addr)

	err := http.ListenAndServe(addr, nil) //NOTE: Blocking function call
	output.Pf("", err, true)
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	output.Close()
}
