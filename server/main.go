package main

import (
	"net/http"

	"github.com/ajagnic/ARImport/output"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

//contentHandler serves static files based on URL path
func contentHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+r.URL.Path)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			output.Log.Fatal("Could not parse form")
		}
		csvFile, head, err := r.FormFile("csv")
		if err != nil {
			output.Log.Fatal(err)
		}
		output.Log.Print(head)
		output.Log.Print(csvFile)
	} else {
		http.ServeFile(w, r, "./static/error.html")
		output.Log.Printf("Invalid Method in /store/: %v", r.Method)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/static/", contentHandler)
	http.HandleFunc("/store", postHandler)

	if err := http.ListenAndServe(":8001", nil); err != nil {
		output.Log.Fatal(err)
	}
	defer output.Close()
}
