package main

import (
	"fmt"
	"strings"

	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
	yee "github.com/nunows/goyeelight"
)

// https://github.com/nunows/goyeelight

type instance struct {
	lights map[string]light
	config *config.YeeConfig
}

type light struct {
	name string
	yee  *yee.Yeelight
}

func (l light) On() {
	if l.yee != nil {
		l.yee.On()
	}
}
func (l light) Off() {
	if l.yee != nil {
		l.yee.Off()
	}
}
func (l light) SetCtAbx(value string, effect string, duration string) {
	if l.yee != nil {
		l.yee.SetCtAbx(value, effect, duration)
	}
}
func (l light) SetBright(value string, effect string, duration string) {
	if l.yee != nil {
		l.yee.SetBright(value, effect, duration)
	}
}

func (l light) Valid() bool {
	return l.yee != nil
}

func new() *instance {
	c := &instance{}
	c.lights = make(map[string]light)
	return c
}

func (c *instance) initialize(jsonstr []byte) error {
	var err error
	c.config, err = config.YeeConfigFromJSON(jsonstr)
	c.lights = make(map[string]light)
	for _, cl := range c.config.Lights {
		l := light{name: cl.Name, yee: yee.New(cl.IP, cl.Port)}
		c.lights[strings.ToLower(cl.Name)] = l
	}
	return err
}

func (c *instance) getLightByName(name string) light {
	l, exists := c.lights[strings.ToLower(name)]
	if exists {
		return l
	}
	return light{}
}

func (c *instance) poweron(name string) {
	light := c.getLightByName(name)
	light.On()
}
func (c *instance) poweroff(name string) {
	light := c.getLightByName(name)
	light.Off()
}

func main() {
	c := new()

	register := []string{"state/light/yee/", "config/request/"}
	subscribe := []string{"config/yee/", "state/light/yee/automation/", "state/light/yee/ahk/", "state/light/yee/flux/"}

	m := microservice.New("yee")
	m.RegisterAndSubscribe(register, subscribe)

	m.RegisterHandler("config/yee/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		c.initialize(msg)
		return true
	})

	m.RegisterHandler("state/light/yee/*/", func(m *microservice.Service, topic string, msg []byte) bool {
		sensor, err := config.SensorStateFromJSON(msg)
		if err == nil {
			m.Logger.LogInfo(m.Name, "received state")
			lightname := sensor.Name
			if lightname != "" {
				light := c.getLightByName(lightname)
				if light.Valid() {
					sensor.ExecValueAttr("power", func(power string) {
						if power == "on" {
							fmt.Println(lightname + " turning On")
							light.On()
						} else if power == "off" {
							fmt.Println(lightname + " turning Off")
							light.Off()
						}
					})
					sensor.ExecFloatAttr("ct", func(ct float64) {
						light.SetCtAbx(fmt.Sprintf("%f", ct), "smooth", "500")
					})
					sensor.ExecFloatAttr("bri", func(bri float64) {
						light.SetBright(fmt.Sprintf("%f", bri), "smooth", "500")
					})
				} else {
					fmt.Println(lightname + " doesn't exist")
				}
			} else if lightname == "all" {
				for _, light := range c.lights {
					sensor.ExecFloatAttr("ct", func(ct float64) {
						light.SetCtAbx(fmt.Sprintf("%f", ct), "smooth", "500")
					})
					sensor.ExecFloatAttr("bri", func(bri float64) {
						light.SetBright(fmt.Sprintf("%f", bri), "smooth", "500")
					})
				}
			} else {
				fmt.Println(lightname + " doesn't exist")
			}
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if tickCount%30 == 0 {
			if c.config == nil {
				m.Pubsub.PublishStr("config/request/", m.Name)
			}
		}
		tickCount++
		return true
	})

	m.Loop()
}
