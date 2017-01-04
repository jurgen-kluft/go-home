package eventlog

import (
	"encoding/json"
)

// EventLogConfig is a structure that holds information to manage an instance of an Event Log
type EventLogConfig struct {
	Filename        string `json:"Filename"`
	Folder          string `json:"Folder"`
	Rolling         bool   `json:"Rolling"`
	MaximumFileSize int    `json:"MaximumFileSize"`
	MaximumLogs     int    `json:"MaximumLogs"`
}

// CreateEventLogConfig creates an instance of EventLogConfig from a stream of bytes containing JSON
func CreateEventLogConfig(jsondata []byte) (log *EventLogConfig) {
	json.Unmarshal(jsondata, log)
	return
}

func (log *EventLogConfig) Initialize() {
	// See if the folder exists
	// Glob all *.log files
	// The number of log files determines the count
	// Sort by date and open the latest
}

func (log *EventLogConfig) SaveEvent(subject string, data []byte) {
	// Save the event in the current open log file
}

func (log *EventLogConfig) rollLog() {
	// See if the current log file is reaching the MaximumFileSize.
	// If it is then close this log file and open a new one.
	// If we have reached the MaximumFileCount then delete the oldest log file.
}
