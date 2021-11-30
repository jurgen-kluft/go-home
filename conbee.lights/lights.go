package main

import (
	"sync/atomic"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
)

/*
This process will scan for events from Conbee, mainly sensors and will send those as
sensor states over NATS.
*/

type lightState struct {
	Name      string
	IDs       []string
	LastSeen  time.Time
	CT        float32
	BRI       float32
	Reachable bool
	OnOff     bool
	Conbee    config.ConbeeLightGroup
}

type fullState struct {
	lights map[string]*lightState
}

func fullStateFromConfig(c *config.ConbeeLightsConfig) fullState {
	full := fullState{}
	full.lights = make(map[string]*lightState)

	for _, e := range c.Lights {
		state := &lightState{Name: e.Name, IDs: e.IDS, LastSeen: time.Now(), Reachable: false, OnOff: false, Conbee: e}
		for _, id := range state.IDs {
			full.lights[id] = state
		}
	}

	return full
}

type signal_t int32

func (b *signal_t) set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32((*int32)(b), int32(i))
}
func (b *signal_t) is_true() bool {
	return atomic.LoadInt32((*int32)(b)) != 0
}
func (b *signal_t) is_not_true() bool {
	return atomic.LoadInt32((*int32)(b)) == 0
}

func main() {
	var cc *config.ConbeeLightsConfig = nil
	var nc *config.ConbeeLightsConfig = nil

	var alive signal_t

	for {
		var err error

		cc = nc

		register := []string{"config/request/", "config/light/conbee/"}
		subscribe := []string{"config/light/conbee/"}

		if cc != nil {
			subscribe = append(subscribe, cc.LightsIn...)
		}

		m := microservice.New("conbee")
		m.RegisterAndSubscribe(register, subscribe)

		m.RegisterHandler("config/light/conbee/", func(m *microservice.Service, topic string, msg []byte) bool {
			m.Logger.LogInfo(m.Name, "Received configuration, schedule restart")
			nc, err = config.ConbeeLightsConfigFromJSON(msg)
			if err != nil {
				m.Logger.LogError(m.Name, err.Error())
			} else {
				cc = nil
				return false
			}
			return true
		})

		tickCount := 0
		m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
			if (tickCount % 30) == 0 {
				if nc == nil {
					m.Logger.LogInfo(m.Name, "Requesting configuration..")
					m.Pubsub.PublishStr("config/request/", m.Name)
				}
			} else if (tickCount % 9) == 0 {
				if cc != nil {
					if alive.is_not_true() {
						m.Logger.LogInfo(m.Name, "Conbee routine is not running, schedule restart..")
						// seems that async_conbee go routine is not running
						// micro-service exit and restart
						return false
					}
				}
			}
			tickCount++
			return true
		})

		m.Loop()

		// Sleep for a while before restarting
		m.Logger.LogInfo(m.Name, "Waiting 10 seconds before re-starting..")
		time.Sleep(10 * time.Second)

		m.Logger.LogInfo(m.Name, "Re-start..")
	}

}
