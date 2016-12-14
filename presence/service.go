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

type Member struct {
	Name string
	Mac string
}

type MemberState struct {
	Self Member
	Presence PresenceState
	DetectDuration Time
	DetectPresence PresenceState
}

type PresenceState struct {
	Host string
	Port int
	User string
	Password string
	DetectionWnd int
	Members []Member
}

type Presence struct {
	State PresenceState
	MemberStates []MemberState
}


func main() {

	// Open REDIS client
	// Command-Line or Config file should specify the connection details of REDIS

	presence_config_key := "Go-Home-Presence-Config"
	presence_state_key := "Go-Home-Presence-State"

	channel_name := "Go-Home"

	// Load config from key 'presence_config_key'

	// Get password for

	// Config is a YAML configuration like this:
	//
	// --- Presence Configuration
	//   host : 192.168.1.3
	//   port : 5000
	//   user : admin
	//   password  : pwd_key
	//   detectionwnd : 200
    //   members:
    //       - name  : Faith
    //         mac   : D8:1D:72:97:51:94
    //       - name  : Jurgen
    //         mac   : 90:3C:92:72:0D:C8
    //       - name  : GrandPa
    //         mac   : BC:44:86:5A:CD:D4
    //       - name  : GrandMa
    //         mac   : 94:94:26:B5:E6:1C
    //           

	// Load state from key 'presence_state_key'

	// Subscribe to channel 'channel_name'

	// Every N seconds 
	//     'collect connected devices'
	//     build list of members ARRIVING/LEAVING/HOME/AWAY
	//     send as YAML to REDIS channel 'channel_name', like:
	//     
	//          ---- PRESENCE
	//          members:
	//              - name  : Faith
	//                state : HOME
	//              - name  : Jurgen
	//                state : LEAVING
	//              - name  : GrandPa
	//                state : ARRIVING
	//              - name  : GrandMa
	//                state : AWAY

	//     save state as YAML in REDIS under key 'presence_state_key'

}