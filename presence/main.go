package main

import "time"

// HOME    is a state that happens when       SEEN > N seconds
// AWAY    is a state that happens when   NOT_SEEN > N seconds

type PresenceState uint32

const (
	Home PresenceState = iota
	Away
)

// LEFT    is a trigger that happens when state changes from HOME => AWAY
// ARRIVED is a trigger that happens when state changes from AWAY => HOME

type PresenceConfig struct {
	Host          string `json:"host"`
	Port          string `json:"port"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Detectionwnd  string `json:"detectionwnd"`
	Detectionfreq string `json:"detectionfreq"`
	Devices       []struct {
		Name string `json:"name"`
		Mac  string `json:"mac"`
	} `json:"devices"`
}

type MemberState struct {
	Name           string
	Presence       PresenceState
	DetectDuration time.Time
	DetectPresence PresenceState
}

type Presence struct {
	State  PresenceState
	Member map[string]MemberState
}

func main() {

	// Open REDIS client
	// Command-Line or Config file should specify the connection details of REDIS

	//presence_config_key := "Go-Home-Presence-Config"
	//presence_state_key := "Go-Home-Presence-State"

	//channel_name := "Go-Home"

	// Load config from key 'presence_config_key'

	// Get password for

	// Config is a JSON configuration like this:
	//
	// {
	//     "host": "192.168.1.3",
	//     "port": "5000",
	//     "user": "admin",
	//     "password": "*********",
	//     "detectionwnd": "200",
	//     "detectionfreq": "120",
	//     "devices": [
	//         {
	//             "name": "Faith",
	//             "mac": "D8:1D:72:97:51:94"
	//         },
	//         {
	//             "name": "Jurgen",
	//             "mac": "90:3C:92:72:0D:C8"
	//         },
	//         {
	//             "name": "GrandPa",
	//             "mac": "BC:44:86:5A:CD:D4"
	//         },
	//         {
	//             "name": "GrandMa",
	//             "mac": "94:94:26:B5:E6:1C"
	//         }
	//     ]
	// }

	// Load state from key 'presence_state_key'

	// Subscribe to channel 'channel_name'

	// router := NewRouter("192.168.1.3", "admin", "postulate518")
	// devices, result := router.GetAttachedDevices()

	// Every N seconds
	//     'collect connected devices'
	//     build list of members ARRIVING/LEAVING/HOME/AWAY
	//     send as JSON to REDIS channel 'channel_name', like:
	//
	// {
	//     "members":
	//		[
	//         {
	//             "name": "Faith",
	//             "state": "HOME"
	//         },
	//         {
	//             "name": "Jurgen",
	//             "state": "LEAVING"
	//         },
	//         {
	//             "name": "GrandPa",
	//             "state": "ARRIVING"
	//         },
	//         {
	//             "name": "GrandMa",
	//             "state": "AWAY"
	//         }
	//     ]
	// }
	// save state in REDIS under key 'presence_state_key'

}
