package main

import (
	"net/http"
	"ARImport/output"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func contentHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "." + r.URL.Path)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/static/", contentHandler)

	err := http.ListenAndServe(":8001", nil)
	output.Log.Fatal(err)
	defer output.Close()
}
