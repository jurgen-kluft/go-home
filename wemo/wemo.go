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
	log     *logpkg.Logger
}

func New() *instance {
	thewemo := &instance{}
	thewemo.devices = map[string]*Switch{}

	thewemo.log = logpkg.New("wemo")
	thewemo.log.AddEntry("emitter")
	thewemo.log.AddEntry("wemo")

	return thewemo
}

func main() {
	thewemo := New()

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/wemo/", "sensor/state/wemo/"}
		subscribe := []string{"config/wemo/", "sensor/state/wemo/"}
		err := client.Connect("wemo", register, subscribe)
		if err == nil {
			thewemo.log.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/wemo/" {
						thewemo.log.LogInfo("wemo", "received configuration")
						thewemo.config, err = config.WemoConfigFromJSON(string(msg.Payload()))
						thewemo.devices = map[string]*Switch{}
						for _, d := range thewemo.config.Devices {
							thewemo.devices[d.Name] = NewSwitch(d.Name, d.IP+":"+d.Port)
						}
					} else if topic == "sensor/state/wemo/" {
						sensor, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							thewemo.log.LogInfo("wemo", "received configuration")
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
							thewemo.log.LogError("wemo", "received bad configuration")
						}
					} else if topic == "client/disconnected/" {
						thewemo.log.LogInfo("emitter", "disconnected")
						connected = false
					}
				}
			}
		}

		if err != nil {
			thewemo.log.LogError("wemo", err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}
