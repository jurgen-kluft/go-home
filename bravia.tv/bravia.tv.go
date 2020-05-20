package main

// Currently: https://github.com/czerwe/gobravia
// Better ? : https://github.com/szatmary/bravia

// https://pro-bravia.sony.net/develop/integrate/ssip/overview/index.html

// Features:
// - Sony Bravia TVs: Turn On/Off

import (
	"github.com/czerwe/gobravia"
	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/micro-service"
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
		}
		return true
	})

	m.RegisterHandler("state/bravia.tv/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received state")
		state, err := config.SensorStateFromJSON(msg)
		if err == nil {
			power := state.GetValueAttr("power", "idle")
			if power == "off" {
				c.poweroff(state.Name)
			} else if power == "on" {
				c.poweron(state.Name)
			}
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if tickCount == 29 {
			tickCount = 0
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
