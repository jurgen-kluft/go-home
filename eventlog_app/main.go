package main

import (
	"github.com/jurgen-kluft/go-home/com"
	"github.com/jurgen-kluft/go-home/eventlog"
)

func main() {
	ghCom := com.New()
	ghCom.Open()

	// Create and initialize event log
	msgEventLogConfigJSON, _ := ghCom.GetKV("GO-HOME-EVENTLOG-CONFIG")
	eventLog := eventlog.Create([]byte(msgEventLogConfigJSON))

	// Initialize EVENT LOGGING
	// Scan log folder and identify state
	eventLog.Initialize()

	// Config is a JSON configuration like this:
	//
	// EventLog Configuration
	// {
	//   "Filename": "go-home-%.log",
	//   "Folder": "~/go-home/log/",
	//   "Rolling": true,
	//   "MaximumFileSize": 10485760,
	//   "MaximumLogs": 10
	// }

	// Subscribe to event-logging channel
	ghChannel, err := ghCom.Subscribe("Go-Home-EventLog")

	if err == nil {
		// Block for message on channel
		for true {
			m, err := ghCom.SubRecv(ghChannel)
			if err == nil {
				//     save event to log
				eventLog.SaveEvent("", []byte(m))
			}
		}
	}

	ghCom.Close()
}
