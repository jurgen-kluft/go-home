package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/metrics"
	"github.com/jurgen-kluft/go-home/pubsub"
)

// AWAY    is a state that happens when   NOT_SEEN > N seconds

type state uint8

const (
	home state = iota
	away
)

func (s state) String() string {
	switch s {
	case home:
		return "home"
	case away:
		return "away"
	default:
		return "unknown"
	}
}

// LEFT    is a trigger that happens when state changes from HOME => AWAY
// ARRIVED is a trigger that happens when state changes from AWAY => HOME

// Config is the configuration needed to call presence.New
// Router is the public interface with one member function to obtain devices
type provider interface {
	get(mac *map[string]bool) error
}

func newProvider(name string, host string, username string, password string) provider {
	switch name {
	case "netgear", "Netgear Nighthawk R6900":
		return newNetgearRouter(host, username, password)
	}
	return nil
}

// Detect holds a single router detection of a member
type Detect struct {
	Time  time.Time
	State state
}

// Member holds the name and detection history with a last and current state
type Member struct {
	name    string
	last    Detect
	current Detect
	index   int
	detect  []Detect
}

// Update will set the current presence state based on historical information
func (m *Member) updateState(current time.Time) {
	m.last = m.current
	for _, d := range m.detect {
		if d.State == home {
			m.current.Time = current
			m.current.State = home
			return
		}
	}
	m.current.Time = current
	m.current.State = away
}

// Presence contains multiple device tracking states
type Presence struct {
	config            *config.PresenceConfig
	provider          provider
	macToIndex        map[string]int
	macToPresence     map[string]bool
	members           []*Member
	updateIntervalSec int
	metrics           *metrics.Metrics
}

// New will return an instance of presence.Home
func New(configjson string) (presence *Presence, err error) {
	config, err := config.PresenceConfigFromJSON(configjson)
	if err != nil {
		return
	}

	var metric *metrics.Metrics
	metric, err = metrics.New()
	if err != nil {
		return nil, err
	}

	presence = &Presence{}
	presence.config = config
	presence.provider = newProvider(presence.config.Name, presence.config.Host, presence.config.User, presence.config.Password)
	presence.macToIndex = map[string]int{}
	presence.macToPresence = map[string]bool{}

	presence.metrics = metric
	metricTags := map[string]string{}
	metricTags["presence"] = "router"

	metricFields := map[string]interface{}{}

	updateHist := presence.config.UpdateHistory
	presence.members = []*Member{}
	for i, device := range presence.config.Devices {
		member := &Member{name: device.Name}
		member.last = Detect{Time: time.Now(), State: away}
		member.current = Detect{Time: time.Now(), State: away}
		member.index = 0
		member.detect = make([]Detect, updateHist, updateHist)
		for j := range member.detect {
			member.detect[j] = Detect{Time: time.Now(), State: away}
		}
		presence.macToIndex[device.Mac] = i
		presence.macToPresence[device.Mac] = false
		presence.members = append(presence.members, member)

		metricFields[member.name] = float64(away)
	}

	presence.metrics.Register("presence", metricTags, metricFields)

	presence.updateIntervalSec = presence.config.UpdateIntervalSec

	return
}

func (p *Presence) presence(currentTime time.Time) error {
	// First reset the presence of every entry in the 'macToPresence' map
	for k := range p.macToPresence {
		p.macToPresence[k] = false
	}
	// Ask the router to update the presence
	result := p.provider.get(&p.macToPresence)
	if result == nil {
		// All members initialize detected presence state to 'away'
		for _, m := range p.members {
			m.detect[m.index] = Detect{Time: currentTime, State: away}
		}
		// For any member registered at the Router mark them as 'home'
		for mac, presence := range p.macToPresence {
			mi, exists := p.macToIndex[mac]
			if exists {
				m := p.members[mi]
				if presence {
					m.detect[m.index] = Detect{Time: currentTime, State: home}
				}
			}
		}
	}

	// Update final presence state for all members
	for _, m := range p.members {
		m.index = (m.index + 1) % len(m.detect)
		m.updateState(currentTime)
	}

	// Report metrics
	p.metrics.Begin("presence")
	for _, m := range p.members {
		p.metrics.Set("presence", m.name, float64(m.current.State))
	}
	p.metrics.Send("presence")

	return result
}

func (p *Presence) publish(now time.Time, client *pubsub.Context) {
	sensor := config.NewSensorState("state.presence")
	sensor.Time = now
	for _, m := range p.members {
		//fmt.Printf("member: %s, presence: %v\n", m.name, m.current.State.String())
		sensor.AddStringAttr(m.name, m.current.State.String())
	}
	jsonstr, err := sensor.ToJSON()
	if err == nil {
		client.Publish("state/presence/", jsonstr)
		//fmt.Println(jsonstr)
	} else {
		fmt.Println(err)
	}
}

func main() {
	var presence *Presence

	logger := logpkg.New("presence")
	logger.AddEntry("emitter")
	logger.AddEntry("presence")

	updateIntervalSec := time.Second * time.Duration(10)
	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/presence/", "state/presence/"}
		subscribe := []string{"config/presence/"}
		err := client.Connect("presence", register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")
			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/presence/" {
						logger.LogInfo("presence", "received configuration")
						presence, err := New(string(msg.Payload()))
						if err == nil {
							updateIntervalSec = time.Second * time.Duration(presence.config.UpdateIntervalSec)
						} else {
							logger.LogError("presence", err.Error())
						}
					} else if topic == "client/disconnected/" {
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}

				case <-time.After(updateIntervalSec):
					if presence != nil {
						now := time.Now()
						err := presence.presence(now)
						if err == nil {
							presence.publish(now, client)
						} else {
							logger.LogError("presence", err.Error())
						}
					}
				}
			}
		}

		if err != nil {
			logger.LogError("presence", err.Error())
		}

		fmt.Println("Waiting 5 seconds before re-connecting to pubsub...")
		time.Sleep(5 * time.Second)
	}
}
