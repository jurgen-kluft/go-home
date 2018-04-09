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
	key     string
	devices map[string]*Switch
	config  *config.WemoConfig
}

func main() {
	thewemo := &instance{}
	thewemo.devices = map[string]*Switch{}

	logger := logpkg.New("wemo")
	logger.AddEntry("emitter")
	logger.AddEntry("wemo")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/wemo/", "sensor/device/wemo/"}
		subscribe := []string{"config/wemo/", "sensor/device/wemo/"}
		err := client.Connect("wemo", register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/wemo/" {
						logger.LogInfo("wemo", "received configuration")
						thewemo.config, err = config.WemoConfigFromJSON(string(msg.Payload()))
						thewemo.devices = map[string]*Switch{}
						for _, d := range thewemo.config.Devices {
							thewemo.devices[d.Name] = NewSwitch(d.Name, d.IP+":"+d.Port)
						}
					} else if topic == "sensor/device/wemo/" {
						sensor, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							logger.LogInfo("wemo", "received configuration")
							devicename := sensor.GetValueAttr("name", "")
							if devicename != "" {
								device, exists := thewemo.devices[devicename]
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
							logger.LogError("wemo", "received bad configuration")
						}
					} else if topic == "client/disconnected/" {
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}

				case <-time.After(time.Second * 10):

				}
			}
		}

		if err != nil {
			logger.LogError("wemo", err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}
