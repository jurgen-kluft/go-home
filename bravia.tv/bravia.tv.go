package main

// Currently: https://github.com/czerwe/gobravia
// Better ? : https://github.com/szatmary/bravia

// https://pro-bravia.sony.net/develop/integrate/ssip/overview/index.html

// Features:
// - Sony Bravia TVs: Turn On/Off

import (
	"fmt"
	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/micro-service"
	"github.com/szatmary/bravia"
)

type instance struct {
	name   string
	ccfg   string
	config *config.BraviaTVConfig
	tvs    map[string]*bravia.Bravia
}

// New ...
func new() *instance {
	c := &instance{}
	c.name = "bravia.tv"
	c.ccfg = "config/bravia.tv/"
	c.tvs = make(map[string]*bravia.Bravia)
	return c
}

func (c *instance) Close() {
	for _, tv := range c.tvs {
		tv.Close()
	}
}

func (c *instance) AddTV(host string, mac string, name string) {
	tv := bravia.NewBravia(host + ":20060")
	c.tvs[name] = tv
}

func (c *instance) changePower(name string, power string) {
	if power == "none" {
		return
	}

	tv, exists := c.tvs[name]
	if exists {
		if power == "on" {
			tv.SetPowerStatus(true)
		} else if power == "off" {
			tv.SetPowerStatus(false)
		}
	}
}
func (c *instance) changeInput(name string, input string) {
	if input == "none" {
		return
	}

	tv, exists := c.tvs[name]
	if exists {
		if input == "1" {
			tv.SetInput(bravia.HDMI, 1)
		} else if input == "2" {
			tv.SetInput(bravia.HDMI, 2)
		} else if input == "3" {
			tv.SetInput(bravia.HDMI, 3)
		} else if input == "4" {
			tv.SetInput(bravia.HDMI, 4)
		}
	}
}

func main() {
	c := new()
	defer c.Close()

	register := []string{c.ccfg, "state/bravia.tv/"}
	subscribe := []string{c.ccfg, "state/bravia.tv/", "config/request/"}

	m := microservice.New("bravia.tv")
	m.RegisterAndSubscribe(register, subscribe)

	m.RegisterHandler("config/bravia.tv/", func(m *microservice.Service, topic string, msg []byte) bool {
		var err error
		c.config, err = config.BraviaTVConfigFromJSON(msg)
		m.Logger.LogInfo(m.Name, "received configuration")
		if err != nil {
			m.Logger.LogError(m.Name, err.Error())
		} else {
			for _, tv := range c.config.Devices {
				c.AddTV(tv.IP, tv.MAC, tv.Name)
				m.Logger.LogInfo(m.Name, fmt.Sprintf("Added TV '%s' with IP '%s' (MAC: %s)", tv.Name, tv.IP, tv.MAC))
			}
		}
		return true
	})

	m.RegisterHandler("state/bravia.tv/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received state")
		state, err := config.SensorStateFromJSON(msg)
		if err == nil {
			power := state.GetValueAttr("power", "none")
			c.changePower(state.Name, power)
			if power != "none" {
				m.Logger.LogInfo(m.Name, fmt.Sprintf("TV '%s'; power -> '%s'", state.Name, power))
			}

			input := state.GetValueAttr("input", "none")
			c.changeInput(state.Name, input)
			if input != "none" {
				m.Logger.LogInfo(m.Name, fmt.Sprintf("TV '%s'; input -> '%s'", state.Name, input))
			}
		} else {
			m.Logger.LogError(m.Name, err.Error())
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if tickCount == 29 {
			tickCount = 0
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
