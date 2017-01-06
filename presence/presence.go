package presence

import (
	"encoding/json"
	"time"
)

// HOME    is a state that happens when       SEEN > N seconds
// AWAY    is a state that happens when   NOT_SEEN > N seconds

type state uint32

const (
	home state = iota
	away
	leaving
	arriving
)

// LEFT    is a trigger that happens when state changes from HOME => AWAY
// ARRIVED is a trigger that happens when state changes from AWAY => HOME

type config struct {
	Host       string  `json:"host"`
	Port       int     `json:"port"`
	User       string  `json:"user"`
	Password   string  `json:"password"`
	UpdateHist int     `json:"uhist"`
	UpdateFreq float64 `json:"ufreq"`
	Devices    []struct {
		Name string `json:"name"`
		Mac  string `json:"mac"`
	} `json:"devices"`
}

func createConfig(jsondata string) (c *config) {
	c = &config{}
	json.Unmarshal([]byte(jsondata), c)
	return
}

type member struct {
	name    string
	current state
	index   int
	detect  []state
}

// Update will set the current presence state based on historical information
func (m *member) update() {

}

func getNameOfState(state state) string {
	switch {
	case state == home:
		return "HOME"
	case state == away:
		return "AWAY"
	case state == leaving:
		return "LEAVING"
	case state == arriving:
		return "ARRIVING"
	default:
		return "HOME"
	}
}

// Home contains multiple device tracking states
type Home struct {
	config          *config
	router          irouter
	macToIndex      map[string]int
	members         []*member
	UpdateFrequency float64
}

// Create will return an instance of presence.Home
func Create(configjson string) *Home {

	presence := &Home{}
	presence.config = createConfig(configjson)
	presence.router = newRouter(presence.config.Host, presence.config.User, presence.config.Password)
	presence.macToIndex = map[string]int{}

	updateHist := presence.config.UpdateHist
	for i, device := range presence.config.Devices {
		member := &member{name: device.Name}
		member.current = home
		member.index = 0
		member.detect = make([]state, updateHist, updateHist)
		for i := range member.detect {
			member.detect[i] = home
		}
		presence.macToIndex[device.Mac] = i
		presence.members = append(presence.members, member)
	}

	presence.UpdateFrequency = presence.config.UpdateFreq

	return presence
}

// MemberState contains the JSON data retrieved from REDIS
type MemberState struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

// Presence is a snapshot state of all detected members
type Presence struct {
	Current time.Time     `json:"datetime"`
	Members []MemberState `json:"members"`
}

// Presence will  a new Presence
func (p *Home) Presence(currentTime time.Time, s *Presence) bool {
	//     'collect connected devices'
	devices, result := p.router.getAttachedDevices()
	if result == nil {
		// All members initialize detected presence state to 'away'
		for _, m := range p.members {
			m.detect[m.index] = away
			m.index = (m.index + 1) % len(m.detect)
		}
		// For any member registered at the Router mark them as 'home'
		for _, device := range devices {
			mi := p.macToIndex[device.mac]
			m := p.members[mi]
			pi := (m.index + len(m.detect) - 1) % len(m.detect)
			m.detect[pi] = home
		}
		// Update final presence state for all members
		for _, m := range p.members {
			m.update()
		}

		// Build JSON structure of members
		// Send as compact JSON to REDIS channel ghChannelName, like:
		// {"datetime":"30/12", "members": [{"name": "Faith", "state": "HOME"},{"name": "Jurgen", "state": "LEAVING"}]}
		s.Current = currentTime
		for _, m := range p.members {
			var member MemberState
			member.Name = m.name
			member.State = getNameOfState(m.current)
			s.Members = append(s.Members, member)
		}
		return true
	}
	return false
}
