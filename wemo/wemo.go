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
		client := pubsub.New()
		err := client.Connect("wemo")

		if err == nil {

			fmt.Println("Connected to emitter")

			client.Register("config/wemo")
			client.Register("sensor/device/wemo")

			client.Subscribe("config/wemo")
			client.Subscribe("sensor/device/wemo")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/wemo" {
						thewemo.config, err = config.WemoConfigFromJSON(string(msg.Payload()))
						thewemo.devices = map[string]*Switch{}
						for _, d := range thewemo.config.Devices {
							thewemo.devices[d.Name] = NewSwitch(d.Name, d.IP+":"+d.Port)
						}
					} else if topic == "sensor/device/wemo" {
						sensor, _ := config.SensorStateFromJSON(string(msg.Payload()))
						devicename := sensor.GetValue("name", "")
						if devicename != "" {
							device, exists := thewemo.devices[devicename]
							if exists {
								power := sensor.GetValue("power", "")
								if power != "" {
									if power == "on" {
										device.On()
									} else if power == "off" {
										device.Off()
									}
								}
							}
						}

					} else if topic == "client/disconnected" {
						connected = false
					}

				case <-time.After(time.Second * 10):

				}
			}
		} else {
			fmt.Println(err.Error())
		}

		// Wait for 10 seconds before retrying
		fmt.Println("Connecting to emitter (retry every 10 seconds)")
		time.Sleep(10 * time.Second)
	}
}
