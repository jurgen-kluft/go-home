package main

import (
	"time"

	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
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

	m, err := microservice.New("wemo", "state/wemo", time.Second*15)
	if err != nil {
		panic(err)
	}

	peers := []string{"config/request"}
	m.ConnectTo(peers)

	m.RegisterHandler("config/request", func(m *microservice.Service, msg *microservice.Message) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		c.config, _ = config.WemoConfigFromJSON(msg.Payload)
		c.devices = map[string]*Switch{}
		for _, d := range c.config.Devices {
			c.devices[d.Name] = NewSwitch(d.Name, d.IP+":"+d.Port)
		}
		return true
	})

	m.RegisterHandler("state/wemo", func(m *microservice.Service, msg *microservice.Message) bool {
		sensor, err := config.SensorStateFromJSON(msg.Payload)
		if err == nil {
			m.Logger.LogInfo(m.Name, "received state")
			devicename := sensor.Name
			if devicename != "" {
				device, exists := c.devices[devicename]
				if exists {
					power := sensor.GetValueAttr("power", "")
					switch power {
					case "on":
						device.On()
					case "off":
						device.Off()
					}
				}
			}
		} else {
			m.Logger.LogError(m.Name, "received bad configuration")
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick", func(m *microservice.Service, msg *microservice.Message) bool {
		if tickCount%5 == 0 {
			if c.config == nil {
				m.SendJsonTo("config/request/", m.Name)
			}
		}
		tickCount++
		return true
	})

	m.Loop()
}
