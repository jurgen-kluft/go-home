package main

func main() {

	// Open REDIS client
	// Command-Line or Config file should specify the connection details of REDIS

	logging_config_key := "Go-Home-EventLog-Config"
	logging_state_key := "Go-Home-EventLog-State"

	channel_name := "Go-Home-EventLog"

	// Load config from key 'logging_config_key'

	// Config is a YAML configuration like this:
	//
	// --- EventLog Configuration
	// Filename : "go-home-{n}.log"
	// Folder : "~/go-home/log"
	// Rolling : Yes
	// Maximum-Size : 10 MiB
	// Maximum-Logs : 10

	// Scan log folder and identify state

	// Subscribe to channel 'channel_name'

	// Block for message on channel
	//
	//     save event to log
	//     check if we need to roll the log
	//

}
