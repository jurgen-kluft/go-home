// Features:
// - Turn On/Off Wemo Switches
// - Publish state of Switches

package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
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

	for {
		client := pubsub.New("tcp://10.0.0.22:8080")
		register := []string{"config/wemo/", "sensor/device/wemo/"}
		subscribe := []string{"config/wemo/", "sensor/device/wemo/"}
		err := client.Connect("wemo", register, subscribe)
		if err == nil {
			fmt.Println("Connected to emitter")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/wemo/" {
						thewemo.config, err = config.WemoConfigFromJSON(string(msg.Payload()))
						thewemo.devices = map[string]*Switch{}
						for _, d := range thewemo.config.Devices {
							thewemo.devices[d.Name] = NewSwitch(d.Name, d.IP+":"+d.Port)
						}
					} else if topic == "sensor/device/wemo/" {
						sensor, _ := config.SensorStateFromJSON(string(msg.Payload()))
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

					} else if topic == "client/disconnected/" {
						connected = false
					}

				case <-time.After(time.Second * 10):

				}
			}
		}

		if err != nil {
			fmt.Println("Error: " + err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}
