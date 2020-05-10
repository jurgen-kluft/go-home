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
	name   string
	ccfg   string
	config *config.BraviaTVConfig
	tvs    map[string]*gobravia.BraviaTV
}

// New ...
func new() *instance {
	c := &instance{}
	c.name = "bravia.tv"
	c.ccfg = "config/bravia.tv/"
	c.tvs = map[string]*gobravia.BraviaTV{}
	return c
}

func (c *instance) AddTV(host string, mac string, name string) {
	tv := gobravia.GetBravia(host, "0000", mac)
	tv.GetCommands()
	c.tvs[name] = tv
}

func (c *instance) poweron(name string) {
	tv, exists := c.tvs[name]
	if exists {
		tv.Poweron("10.0.0.255")
	}
}
func (c *instance) poweroff(name string) {
	tv, exists := c.tvs[name]
	if exists {
		tv.SendAlias("poweroff")
	}
}

func main() {
	c := new()

	logger := logpkg.New(c.name)
	logger.AddEntry("emitter")
	logger.AddEntry(c.name)

	for {
		client := pubsub.New(config.PubSubCfg)
		register := []string{c.ccfg, "state/bravia.tv/"}
		subscribe := []string{c.ccfg, "state/bravia.tv/", "config/request/"}
		err := client.Connect(c.name, register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == c.ccfg {
						c.config, err = config.BraviaTVConfigFromJSON(msg.Payload())
						logger.LogInfo(c.name, "received configuration")
					} else if topic == "state/bravia.tv/" {
						logger.LogInfo(c.name, "received state")
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							power := state.GetValueAttr("power", "idle")
							if power == "off" {
								c.poweroff(state.Name)
							} else if power == "on" {
								c.poweron(state.Name)
							}
						}
					} else if topic == "client/disconnected/" {
						connected = false
					}

				case <-time.After(time.Minute * 1): // Try and request our configuration
					if c.config == nil {
						client.Publish("config/request/", "bravia.tv")
					}

				}
			}
		}

		if err != nil {
			logger.LogError(c.name, err.Error())
		}
		time.Sleep(5 * time.Second)

	}
}
