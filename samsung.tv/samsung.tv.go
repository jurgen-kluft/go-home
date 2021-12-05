package main

// Features:
// - Samsung TV: Turn On/Off

import (
	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
	"github.com/saljam/samote"
)

type tv struct {
	name   string
	host   string
	id     string
	remote samote.Remote
}

type instance struct {
	name   string
	config *config.SamsungTVConfig
	tvs    map[string]*tv
}

// New ...
func new() *instance {
	c := &instance{}
	c.name = "tv/samsung/"
	c.config = nil
	c.tvs = make(map[string]*tv)
	return c
}

func (c *instance) Add(host string, name string, id string) error {
	var err error
	tv := &tv{name: name, host: host, id: id}
	tv.remote, err = samote.Dial(host, name, id)
	if err == nil {
		c.tvs[name] = tv
	} else {
		return err
	}
	return err
}

func (c *instance) poweron(name string) error {
	tv, exists := c.tvs[name]
	if exists {
		remote, err := samote.Dial(tv.host, tv.name, tv.id)
		if err == nil {
			remote.SendKey(samote.KEY_POWERON)
		} else {
			return err
		}
	}
	return nil
}
func (c *instance) poweroff(name string) error {
	tv, exists := c.tvs[name]
	if exists {
		remote, err := samote.Dial(tv.host, tv.name, tv.id)
		if err == nil {
			remote.SendKey(samote.KEY_POWEROFF)
		} else {
			return err
		}
	}
	return nil
}

func main() {
	register := []string{"config/tv/samsung/", "state/tv/samsung/", "config/request/"}
	subscribe := []string{"config/tv/samsung/", "state/tv/samsung/"}

	c := new()

	m := microservice.New("tv/samsung/")
	m.RegisterAndSubscribe(register, subscribe)

	m.RegisterHandler("config/tv/samsung/", func(m *microservice.Service, topic string, msg []byte) bool {
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

	m.RegisterHandler("state/tv/samsung/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		state, err := config.SensorStateFromJSON(msg)
		if err == nil {
			power := state.GetValueAttr("power", "idle")
			if power == "off" {
				err = c.poweroff(state.Name)
			} else if power == "on" {
				err = c.poweron(state.Name)
			}
			if err != nil {
				m.Logger.LogError(m.Name, err.Error())
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
				m.Pubsub.PublishStr("config/request/", m.Name)
			}
		} else {
			tickCount++
		}
		return true
	})

	m.Loop()

}
