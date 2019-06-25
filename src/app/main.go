package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/ajagnic/ARImport/src/output"
	"github.com/ajagnic/ARImport/src/scheduler"
	"github.com/ajagnic/ARImport/src/server"
)

func main() {
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~SCHEDULER
	reinit := make(chan bool, 1)
	killsched := make(chan bool, 1)
	_ = scheduler.Config()
	go scheduler.EventListener(reinit, killsched)

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~SERVER
	server.Run()
	// srv := http.Server{
	// 	Addr:     addr,
	// 	ErrorLog: output.Log,
	// }

	// http.HandleFunc("/", indexHandler)
	// http.HandleFunc("/static/", contentHandler)
	// http.HandleFunc("/store", postHandler)
	// http.HandleFunc("/config", configHandler)

	//sigint listens for interrupt signal then cleans up main()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	// go func() {
	// 	fmt.Printf("Starting server on URL/Port: %v\n", addr)
	// 	err := srv.ListenAndServe()
	// 	output.Pf("", err, true)
	// }()

	<-sigint

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~CLEANUP
	killsched <- true

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	// defer cancel()
	// err := srv.Shutdown(ctx)
	// output.Pf("", err, true)

	output.Close()
	fmt.Println("Server gracefully shutdown.")
}
