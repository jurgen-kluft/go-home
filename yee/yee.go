package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/micro-service"
	yee "github.com/nunows/goyeelight"
)

// https://github.com/nunows/goyeelight

type instance struct {
	lamps  map[string]*yee.Yeelight
	config *config.YeeConfig
}

func new() *instance {
	c := &instance{}
	c.lamps = map[string]*yee.Yeelight{}

	c.logger = logpkg.New(c.name)
	c.logger.AddEntry("emitter")
	c.logger.AddEntry(c.name)

	return c
}

func (c *instance) initialize(jsonstr []byte]) error {
	var err error
	c.config, err = config.YeeConfigFromJSON(jsonstr)
	c.lamps = map[string]*yee.Yeelight{}
	for _, lamp := range c.config.Lights {
		c.lamps[lamp.Name] = yee.New(lamp.IP, lamp.Port)
	}
	return err
}

func (c *instance) poweron(name string) {
	lamp, exists := c.lamps[name]
	if exists {
		lamp.On()
	}
}
func (c *instance) poweroff(name string) {
	lamp, exists := c.lamps[name]
	if exists {
		lamp.Off()
	}
}

func main() {
	c := new()

	register := []string{"config/yee/", "state/sensor/yee/", "state/light/yee/"}
	subscribe := []string{"config/yee/", "state/sensor/yee/", "state/light/yee/", "config/request/"}

	m := microservice.New("yee")
	m.RegisterAndSubscribe(register, subscribe)

	m.RegisterHandler("config/yee/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		c.initialize(msg)
		return true
	})

	m.RegisterHandler("state/light/yee/", func(m *microservice.Service, topic string, msg []byte) bool {
		yeesensor, err := config.SensorStateFromJSON(msg)
		if err == nil {
			m.Logger.LogInfo(m.Name, "received state")
			lampname := yeesensor.GetValueAttr("name", "")
			if lampname != "" {
				lamp, exists := c.lamps[lampname]
				if exists {
					yeesensor.ExecValueAttr("power", func(power string) {
						if power == "on" {
							lamp.On()
						} else if power == "off" {
							lamp.Off()
						}
					})
					yeesensor.ExecFloatAttr("ct", func(ct float64) {
						lamp.SetCtAbx(fmt.Sprintf("%f", ct), "smooth", "500")
					})
					yeesensor.ExecFloatAttr("bri", func(bri float64) {
						lamp.SetBright(fmt.Sprintf("%f", bri), "smooth", "500")
					})
				}
			} else if lampname == "all" {
				for _, lamp := range c.lamps {
					yeesensor.ExecFloatAttr("ct", func(ct float64) {
						lamp.SetCtAbx(fmt.Sprintf("%f", ct), "smooth", "500")
					})
					yeesensor.ExecFloatAttr("bri", func(bri float64) {
						lamp.SetBright(fmt.Sprintf("%f", bri), "smooth", "500")
					})
				}
			}
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if tickCount%30 == 0 {
			if c.config == nil {
				m.Pubsub.Publish("config/request/", m.Name)
			}
		}
		return true
	})

	m.Loop()
}
