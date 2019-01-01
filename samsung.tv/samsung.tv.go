package main

// Features:
// - Samsung TV: Turn On/Off

import (
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/saljam/samote"
)

type instance struct {
	name   string
	config *config.SamsungTVConfig
	tvs    map[string]samote.Remote
}

// New ...
func new() *instance {
	c := &instance{}
	c.name = "samsung.tv"
	c.config = nil
	c.tvs = map[string]samote.Remote{}
	return c
}

func (c *instance) Add(host string, name string, id string) error {
	var err error
	var remote samote.Remote
	remote, err = samote.Dial(host, name, id)
	if err == nil {
		c.tvs[name] = remote
	}
	return err
}

func (c *instance) poweron(name string) {
	remote, exists := c.tvs[name]
	if exists {
		remote.SendKey(samote.KEY_POWERON)
	}
}
func (c *instance) poweroff(name string) {
	remote, exists := c.tvs[name]
	if exists {
		remote.SendKey(samote.KEY_POWEROFF)
	}
}

func main() {
	c := new()

	logger := logpkg.New(c.name)
	logger.AddEntry("emitter")
	logger.AddEntry(c.name)

	for {
		client := pubsub.New(config.EmitterIOCfg)
		register := []string{"config/samsung.tv/", "state/samsung.tv/"}
		subscribe := []string{"config/samsung.tv/", "state/samsung.tv/"}
		err := client.Connect(c.name, register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/samsung.tv/" {
						logger.LogInfo(c.name, "received configuration")
						c.config, err = config.SamsungTVConfigFromJSON(string(msg.Payload()))
						for _, tv := range c.config.Devices {
							err = c.Add(tv.IP, tv.Name, tv.ID)
							if err == nil {
								logger.LogInfo(c.name, "registered TV with name "+tv.Name)
							} else {
								logger.LogError(c.name, err.Error())
							}
						}
					} else if topic == "state/samsung.tv/" {
						logger.LogInfo(c.name, "received configuration")
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
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}

				case <-time.After(time.Second * 10):

				}
			}
		}

		if err != nil {
			logger.LogError(c.name, err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}
