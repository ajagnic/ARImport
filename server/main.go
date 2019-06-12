package main

import (
	"net/http"
	"fmt"

	"github.com/ajagnic/ARImport/output"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

//contentHandler serves static files based on URL path.
func contentHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+r.URL.Path)
}

//postHandler parses form data POSTed to /store.
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// if err := r.ParseForm(); err != nil {
		// 	output.Log.Print("Could not parse form in /store: ", err)
		// }
		// if err := r.ParseMultipartForm(); err != nil {
		// 	output.Log.Print("Could not parse form in /store: ", err)
		// }

		// csvFile, head, err := r.FormFile("csv")
		// if err != nil {
		// 	output.Log.Fatal(err)
		// }
		// output.Log.Print(head)
		// output.Log.Print(csvFile)
	} else {
		http.ServeFile(w, r, "./static/error.html")
		output.Log.Printf("Invalid Method in /store: %v", r.Method)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/static/", contentHandler)
	http.HandleFunc("/store", postHandler)

	addr := ":8001"

	fmt.Printf("Starting server on URL/Port: %v\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		output.Log.Fatal(err)
	}
	defer output.Close()
}
