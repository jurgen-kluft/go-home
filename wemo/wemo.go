// Features:
// - Turn On/Off Wemo Switches
// - Publish state of Switches

package main

import (
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
)

type instance struct {
	name    string
	key     string
	devices map[string]*Switch
	config  *config.WemoConfig
	log     *logpkg.Logger
}

func new() *instance {
	c := &instance{}
	c.name = "wemo"
	c.devices = map[string]*Switch{}

	c.log = logpkg.New(c.name)
	c.log.AddEntry("emitter")
	c.log.AddEntry(c.name)

	return c
}

func main() {
	c := new()

	for {
		client := pubsub.New(config.EmitterIOCfg)
		register := []string{"config/wemo/", "sensor/state/wemo/"}
		subscribe := []string{"config/wemo/", "sensor/state/wemo/", "config/request/"}
		err := client.Connect(c.name, register, subscribe)
		if err == nil {
			c.log.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/wemo/" {
						c.log.LogInfo(c.name, "received configuration")
						c.config, err = config.WemoConfigFromJSON(string(msg.Payload()))
						c.devices = map[string]*Switch{}
						for _, d := range c.config.Devices {
							c.devices[d.Name] = NewSwitch(d.Name, d.IP+":"+d.Port)
						}
					} else if topic == "sensor/state/wemo/" {
						sensor, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							c.log.LogInfo(c.name, "received configuration")
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
							c.log.LogError(c.name, "received bad configuration")
						}
					} else if topic == "client/disconnected/" {
						c.log.LogInfo("emitter", "disconnected")
						connected = false
					}

				case <-time.After(time.Minute * 1): // Try and request our configuration
					if c.config == nil {
						client.Publish("config/request/", "wemo")
					}

				}
			}
		}

		if err != nil {
			c.log.LogError(c.name, err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}
