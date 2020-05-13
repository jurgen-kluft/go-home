package main

// Features:
// - Samsung TV: Turn On/Off

import (
	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/micro-service"
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
	register := []string{"config/samsung.tv/", "state/samsung.tv/"}
	subscribe := []string{"config/samsung.tv/", "state/samsung.tv/", "config/request/"}

	c := new()

	m := microservice.New("samsung.tv")
	m.RegisterAndSubscribe(register, subscribe)

	m.RegisterHandler("config/samsung.tv/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		var err error
		c.config, err = config.SamsungTVConfigFromJSON(msg)
		for _, tv := range c.config.Devices {
			err = c.Add(tv.IP, tv.Name, tv.ID)
			if err == nil {
				m.Logger.LogInfo(m.Name, "registered TV with name "+tv.Name)
			} else {
				m.Logger.LogError(m.Name, err.Error())
			}
		}
		return true
	})

	m.RegisterHandler("state/samsung.tv/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		state, err := config.SensorStateFromJSON(msg)
		if err == nil {
			power := state.GetValueAttr("power", "idle")
			if power == "off" {
				c.poweroff(state.Name)
			} else if power == "on" {
				c.poweron(state.Name)
			}
		} else {
			m.Logger.LogError(m.Name, err.Error())
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if tickCount%30 == 0 {
			if c.config == nil {
				m.Pubsub.Publish("config/request/", m.Name)
			}
		} else {
			tickCount++
		}
		return true
	})

	m.Loop()

}
