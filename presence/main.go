package main

import (
	"encoding/json"
	"gopkg.in/redis.v5"
	"time"
)

// HOME    is a state that happens when       SEEN > N seconds
// AWAY    is a state that happens when   NOT_SEEN > N seconds

type presenceState uint32

const (
	home presenceState = iota
	away
	leaving
	arriving
)

// LEFT    is a trigger that happens when state changes from HOME => AWAY
// ARRIVED is a trigger that happens when state changes from AWAY => HOME

// PresenceConfig contains the JSON data retrieved from REDIS
type PresenceConfig struct {
	Host          string `json:"host"`
	Port          int    `json:"port"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Detectionwnd  int    `json:"detectionwnd"`
	Detectionfreq int    `json:"detectionfreq"`
	Devices       []struct {
		Name string `json:"name"`
		Mac  string `json:"mac"`
	} `json:"devices"`
}

// CreateConfig creates and instance of PresenceConfig from provided JSON string
func CreateConfig(jsondata string) (config *PresenceConfig) {
	config = &PresenceConfig{}
	json.Unmarshal([]byte(jsondata), config)
	return
}

// MemberPresence tracks the state of a device on the LAN
type MemberPresence struct {
	Name           string
	Presence       presenceState
	IndexPresence  int
	DetectPresence []presenceState
}

func (m *MemberPresence) UpdatePresence() {

}

func getNameOfState(state presenceState) string {
	switch state {
	case home:
		return "HOME"
	case away:
		return "AWAY"
	case leaving:
		return "LEAVING"
	case arriving:
		return "ARRIVING"
	default:
		return "HOME"
	}
}

// HomePresence contains multiple device tracking states
type HomePresence struct {
	MacToIndex map[string]int
	Member     []*MemberPresence
}

func main() {
	ghChannelName := "Go-Home"
	ghConfigKey := "GO-HOME-PRESENCE-CONFIG"

	// Open REDIS and read all the configurations
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Create and initialize
	ghConfigJSON, err := redisClient.Get(ghConfigKey).Result()
	if err != nil {
		panic(err)
	}

	presence := &HomePresence{}
	presence.MacToIndex = map[string]int{}
	ghConfig := CreateConfig(ghConfigJSON)
	for i, device := range ghConfig.Devices {
		member := &MemberPresence{Name: device.Name}
		member.Presence = home
		member.IndexPresence = 0
		member.DetectPresence = make([]presenceState, 6, 6)
		for i := range member.DetectPresence {
			member.DetectPresence[i] = away
		}
		presence.MacToIndex[device.Mac] = i
		presence.Member[i] = member
	}

	// Config is a JSON configuration like this:
	//
	//  {
	//  	"host": "192.168.1.3",
	//  	"port": "5000",
	//  	"user": "admin",
	//  	"password": "*********",
	//  	"detectionwnd": "200",
	//  	"detectionfreq": "60",
	//  	"devices": [{
	//  		"name": "Faith",
	//  		"mac": "D8:1D:72:97:51:94"
	//  	}, {
	//  		"name": "Jurgen",
	//  		"mac": "90:3C:92:72:0D:C8"
	//  	}, {
	//  		"name": "GrandPa",
	//  		"mac": "BC:44:86:5A:CD:D4"
	//  	}, {
	//  		"name": "GrandMa",
	//  		"mac": "94:94:26:B5:E6:1C"
	//  	}]
	//  }

	router := NewRouter("192.168.1.3", ghConfig.User, ghConfig.Password)

	// Every N seconds
	ticker := time.NewTicker(time.Second * time.Duration(ghConfig.Detectionfreq))
	for currentTime := range ticker.C {
		//     'collect connected devices'
		devices, result := router.GetAttachedDevices()
		if result == nil {
			for _, device := range devices {
				mi := presence.MacToIndex[device.Mac]
				m := presence.Member[mi]
				m.DetectPresence[m.IndexPresence] = home
				m.IndexPresence = (m.IndexPresence + 1) % len(m.DetectPresence)
				m.UpdatePresence()
			}

			// Build JSON structure of members
			// Send as compact JSON to REDIS channel ghChannelName, like:
			// {"datetime":"30/12", "members": [{"name": "Faith", "state": "HOME"},{"name": "Jurgen", "state": "LEAVING"}]}
			presenceInfo := "{ \"datetime\": "
			presenceInfo += currentTime.String()
			presenceInfo += ", \"members\": ["
			for i, m := range presence.Member {
				if i > 0 {
					presenceInfo += ",{ name:\"" + m.Name + "\"" + "state:\"" + getNameOfState(m.Presence) + "\"}"
				} else {
					presenceInfo += "{ name:\"" + m.Name + "\"" + "state:\"" + getNameOfState(m.Presence) + "\"}"
				}
			}
			presenceInfo += "]}"

			redisClient.Publish(ghChannelName, presenceInfo).Result()
		}
	}

}
