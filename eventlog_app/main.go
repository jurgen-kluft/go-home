package main

import (
	"github.com/jurgen-kluft/go-home/eventlog"
	"gopkg.in/redis.v5"
)

func main() {
	// Open REDIS and read all the configurations
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Create and initialize event log
	msgEventLogConfigJSON, _ := redisClient.Get("GO-HOME-EVENTLOG-CONFIG").Result()
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
	ghEventLoggingChannel := "Go-Home-EventLog"
	redisPubSub, err := redisClient.Subscribe(ghEventLoggingChannel)
	if err == nil {
		// Block for message on channel
		for true {
			m, err := redisPubSub.ReceiveMessage()
			if err == nil {
				//     save event to log
				eventLog.SaveEvent("", []byte(m.Payload))
			}
		}
	}
}
