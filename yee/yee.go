package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
	yee "github.com/nunows/goyeelight"
)

// https://github.com/nunows/goyeelight

type instance struct {
	name   string
	key    string
	lamps  map[string]*yee.Yeelight
	config *config.YeeConfig
	logger *logpkg.Logger
}

func new() *instance {
	c := &instance{}
	c.name = "yee"
	c.lamps = map[string]*yee.Yeelight{}

	c.logger = logpkg.New(c.name)
	c.logger.AddEntry("emitter")
	c.logger.AddEntry(c.name)

	return c
}

func (c *instance) initialize(jsonstr string) error {
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

	for {
		client := pubsub.New(config.EmitterIOCfg)
		register := []string{"config/yee/", "state/sensor/yee/", "state/light/yee/"}
		subscribe := []string{"config/yee/", "state/sensor/yee/", "state/light/yee/"}
		err := client.Connect(c.name, register, subscribe)
		if err == nil {
			c.logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/yee/" {
						c.logger.LogInfo(c.name, "received configuration")
						c.initialize(string(msg.Payload()))
					} else if topic == "state/light/yee/" {
						yeesensor, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							c.logger.LogInfo(c.name, "received state")

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

					} else if topic == "client/disconnected/" {
						c.logger.LogInfo("emitter", "disconnected")
						connected = false
					}

				case <-time.After(time.Second * 10):

				}
			}
		}

		if err != nil {
			c.logger.LogError(c.name, err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}
