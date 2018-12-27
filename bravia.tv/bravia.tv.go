package main

// https://github.com/czerwe/gobravia

// Features:
// - Sony Bravia TVs: Turn On/Off

import (
	"time"

	"github.com/czerwe/gobravia"
	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
)

type instance struct {
	config *config.BraviaTVConfig
	tvs    map[string]*gobravia.BraviaTV
}

// New ...
func New() *instance {
	x := &instance{}
	x.tvs = map[string]*gobravia.BraviaTV{}
	return x
}

func (x *instance) AddTV(host string, mac string, name string) {
	tv := gobravia.GetBravia(host, "0000", mac)
	tv.GetCommands()
	x.tvs[name] = tv
}

func (x *instance) poweron(name string) {
	tv, exists := x.tvs[name]
	if exists {
		tv.Poweron("10.0.0.255")
	}
}
func (x *instance) poweroff(name string) {
	tv, exists := x.tvs[name]
	if exists {
		tv.SendAlias("poweroff")
	}
}

func main() {
	bravia := New()

	logger := logpkg.New("bravia.tv")
	logger.AddEntry("emitter")
	logger.AddEntry("bravia.tv")

	for {
		client := pubsub.New(config.EmitterIOCfg)
		register := []string{"config/bravia.tv/", "state/bravia.tv/"}
		subscribe := []string{"config/bravia.tv/", "state/bravia.tv/"}
		err := client.Connect("bravia.tv", register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/bravia.tv/" {
						bravia.config, err = config.BraviaTVConfigFromJSON(string(msg.Payload()))
						logger.LogInfo("bravia.tv", "received configuration")
					} else if topic == "state/bravia.tv/" {
						logger.LogInfo("bravia.tv", "received state")
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							power := state.GetValueAttr("power", "idle")
							if power == "off" {
								bravia.poweroff(state.Name)
							} else if power == "on" {
								bravia.poweron(state.Name)
							}
						}
					} else if topic == "client/disconnected/" {
						connected = false
					}

				case <-time.After(time.Second * 10):

				}
			}
		}

		if err != nil {
			logger.LogError("bravia.tv", err.Error())
		}
		time.Sleep(5 * time.Second)

	}
}
