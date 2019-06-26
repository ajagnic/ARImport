package main

import (
	"github.com/ajagnic/ARImport/src/output"
	"github.com/ajagnic/ARImport/src/scheduler"
	"github.com/ajagnic/ARImport/src/server"
)

func main() {
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~SCHEDULER
	reinit := make(chan bool, 1)
	killsched := make(chan bool, 1)
	addr := scheduler.Config()
	go scheduler.EventListener(reinit, killsched)

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~SERVER
	server.Run(addr, reinit) //Blocking call.

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~CLEANUP
	killsched <- true
	server.Shutdown()
	output.Close()
}
