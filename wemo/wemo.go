// Features:
// - Turn On/Off Wemo Switches
// - Publish state of Switches

package main

import (
	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/micro-service"
)

type instance struct {
	devices map[string]*Switch
	config  *config.WemoConfig
}

func new() *instance {
	c := &instance{}
	c.devices = map[string]*Switch{}
	return c
}

func main() {
	c := new()

	register := []string{"config/wemo/", "sensor/state/wemo/"}
	subscribe := []string{"config/wemo/", "sensor/state/wemo/", "config/request/"}

	m := microservice.New("wemo")
	m.RegisterAndSubscribe(register, subscribe)

	m.RegisterHandler("config/wemo/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		var err error
		c.config, err = config.WemoConfigFromJSON(msg)
		c.devices = map[string]*Switch{}
		for _, d := range c.config.Devices {
			c.devices[d.Name] = NewSwitch(d.Name, d.IP+":"+d.Port)
		}
		return true
	})

	m.RegisterHandler("sensor/state/wemo/", func(m *microservice.Service, topic string, msg []byte) bool {
		sensor, err := config.SensorStateFromJSON(msg)
		if err == nil {
			m.Logger.LogInfo(m.Name, "received state")
			devicename := sensor.GetValueAttr("name", "")
			if devicename != "" {
				device, exists := c.devices[devicename]
				if exists {
					power := sensor.GetValueAttr("power", "")
					if power != "" {
						if power == "on" {
							device.On()
						} else if power == "off" {
							device.Off()
						}
					}
				}
			}
		} else {
			m.Logger.LogError(m.Name, "received bad configuration")
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if tickCount%5 == 0 {
			if c.config == nil {
				m.Pubsub.Publish("config/request/", m.Name)
			}
		}
		tickCount++
		return true
	})

	m.Loop()
}
