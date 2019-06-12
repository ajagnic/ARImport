package main

import (
	"log"
	"net/http"
)

func main() {
	logger := log.New(nil, "", 0)
	logger.Fatal(http.ListenAndServe("127.0.0.1:8001", nil))
}
