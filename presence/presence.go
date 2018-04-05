package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/pubsub"
)

// HOME    is a state that happens when   NOT_AWAY > X seconds
// AWAY    is a state that happens when   NOT_SEEN > N seconds

type state uint32

const (
	home state = iota
	away
	leaving
	arriving
)

func (s state) String() string {
	switch s {
	case home:
		return "home"
	case away:
		return "away"
	case leaving:
		return "leaving"
	case arriving:
		return "arriving"
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

type Detect struct {
	Time  time.Time
	State state
}

type Member struct {
	name    string
	last    Detect
	current Detect
	index   int
	detect  []Detect
}

type PresenceState struct {
	time    time.Time
	members []Member
}

// Update will set the current presence state based on historical information
func (m *Member) updateCurrent(current time.Time) {
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

// Presence contains multiple device tracking states
type Presence struct {
	config            *config.PresenceConfig
	provider          provider
	macToIndex        map[string]int
	macToPresence     *map[string]bool
	members           []*Member
	updateIntervalSec int
}

// New will return an instance of presence.Home
func New(configjson string) *Presence {
	config, err := config.PresenceConfigFromJSON(configjson)
	if err != nil {
		return nil
	}
	presence := &Presence{}
	presence.config = config
	presence.provider = newProvider(presence.config.Name, presence.config.Host, presence.config.User, presence.config.Password)
	presence.macToIndex = map[string]int{}
	presence.macToPresence = &map[string]bool{}

	updateHist := presence.config.UpdateHistory
	for i, device := range presence.config.Devices {
		member := &Member{name: device.Name}
		member.last = Detect{Time: time.Now(), State: away}
		member.current = Detect{Time: time.Now(), State: home}
		member.index = 0
		member.detect = make([]Detect, updateHist, updateHist)
		for i := range member.detect {
			member.detect[i] = Detect{Time: time.Now(), State: away}
		}
		presence.macToIndex[device.Mac] = i
		(*presence.macToPresence)[device.Mac] = false
		presence.members = append(presence.members, member)
	}

	presence.updateIntervalSec = presence.config.UpdateIntervalSec

	return presence
}

// Presence  ...
func (p *Presence) Presence(currentTime time.Time) bool {
	result := p.provider.get(p.macToPresence)
	if result == nil {
		// All members initialize detected presence state to 'away'
		for _, m := range p.members {
			m.detect[m.index] = Detect{Time: currentTime, State: away}
			m.index = (m.index + 1) % len(m.detect)
		}
		// For any member registered at the Router mark them as 'home'
		for mac, presence := range *p.macToPresence {
			mi, exists := p.macToIndex[mac]
			if exists {
				m := p.members[mi]
				pi := (m.index + len(m.detect) - 1) % len(m.detect)
				if presence {
					m.detect[pi] = Detect{Time: currentTime, State: home}
				} else {
					m.detect[pi] = Detect{Time: currentTime, State: away}
				}
			}
		}
		// Update final presence state for all members
		for _, m := range p.members {
			m.updateCurrent(currentTime)
		}

		return true
	}
	fmt.Println(result)
	return false
}

func (p *Presence) publish(now time.Time, client *pubsub.Context) {
	sensor := config.NewSensorState("state.sensor.presence")
	sensor.Time = now
	for _, m := range p.members {
		fmt.Printf("member: %s, presence: %v\n", m.name, m.current.State.String())
		sensor.AddValueSensor(m.name, m.current.State.String())
	}
	jsonstr, err := sensor.ToJSON()
	if err == nil {
		client.Publish("state/sensor/presence", jsonstr)
		fmt.Println(jsonstr)
	} else {
		fmt.Println(err)
	}
}

func main() {

	var presence *Presence

	updateIntervalSec := time.Second * 15

	for {
		connected := true
		for connected {
			client := pubsub.New("tcp://10.0.0.22:8080")

			err := client.Connect("presence")
			if err == nil {

				err = client.Register("config/presence/")
				if err == nil {

					err = client.Subscribe("config/presence/")
					if err == nil {

						for connected {
							select {
							case msg := <-client.InMsgs:
								fmt.Printf("Emitter message received, topic:'%s', msg:'%s'\n", msg.Topic(), string(msg.Payload()))

								topic := msg.Topic()
								if topic == "config/presence/" {
									fmt.Println("Received configuration ...")
									presence = New(string(msg.Payload()))
									updateIntervalSec = time.Second * time.Duration(presence.config.UpdateIntervalSec)
								} else if topic == "client/disconnected/" {
									connected = false
								}

							case <-time.After(updateIntervalSec):
								if presence != nil {
									now := time.Now()
									if presence.Presence(now) {
										presence.publish(now, client)
									}
								}

							}
						}
					}
				}
			}

			if err != nil {
				fmt.Println("Error: " + err.Error())
			}
		}

		fmt.Println("Waiting 10 seconds before re-connecting to pubsub...")
		time.Sleep(10 * time.Second)
	}
}
