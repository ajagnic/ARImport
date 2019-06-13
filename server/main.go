package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ajagnic/ARImport/output"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

//contentHandler serves static files based on URL path.
func contentHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+r.URL.Path)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rdr, err := r.MultipartReader()
		output.Pf("File reader: %v", err, false)

		for {
			part, err := rdr.NextPart()
			if err == io.EOF {
				break
			} else {
				output.Pf("File parts: %v", err, false)
			}

			if part.FileName() == "" {
				continue
			}

			file, err := os.OpenFile("./server/csv/"+part.FileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			output.Pf("", err, true)
			defer file.Close()

			_, err = io.CopyBuffer(file, part, nil)
			output.Pf("CopyBuffer: %v", err, false)
		}
	} else {
		http.ServeFile(w, r, "./static/error.html")
		output.Log.Printf("Invalid method in /store: %v", r.Method)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/static/", contentHandler)
	http.HandleFunc("/store", postHandler)

	addr := ":8001"

	fmt.Printf("Starting server on URL/Port: %v\n", addr)

	err := http.ListenAndServe(addr, nil)
	output.Pf("", err, true)

	output.Close()
}
