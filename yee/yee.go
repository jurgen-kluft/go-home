package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/pubsub"
	yee "github.com/nunows/goyeelight"
)

// https://github.com/nunows/goyeelight

type instance struct {
	key    string
	lamps  map[string]*yee.Yeelight
	config *config.YeeConfig
}

func main() {
	yeelighting := &instance{}
	yeelighting.lamps = map[string]*yee.Yeelight{}

	for {
		client := pubsub.New()
		err := client.Connect("yee")

		if err == nil {

			fmt.Println("Connected to emitter")

			client.Register("config/yee")
			client.Register("sensor/light/yee")

			client.Subscribe("config/yee")
			client.Subscribe("sensor/light/yee")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/yee" {
						yeelighting.config, err = config.YeeConfigFromJSON(string(msg.Payload()))
						yeelighting.lamps = map[string]*yee.Yeelight{}
						for _, lamp := range yeelighting.config.Lights {
							yeelighting.lamps[lamp.Name] = yee.New(lamp.IP, lamp.Port)
						}
					} else if topic == "sensor/light/yee" {
						yeesensor, _ := config.SensorStateFromJSON(string(msg.Payload()))
						lampname := yeesensor.GetValue("name", "")
						if lampname != "" {
							lamp, exists := yeelighting.lamps[lampname]
							if exists {
								power := yeesensor.GetValue("power", "")
								if power != "" {
									if power == "on" {
										lamp.On()
									} else if power == "off" {
										lamp.Off()
									}
								}
								if power == "on" {
									ct := yeesensor.GetFloatValue("ct", -1.0)
									if ct != -1.0 {
										lamp.SetCtAbx(fmt.Sprintf("%f", ct), "smooth", "500")
									}
									bri := yeesensor.GetFloatValue("bri", -1.0)
									if ct != -1.0 {
										lamp.SetBright(fmt.Sprintf("%f", bri), "smooth", "500")
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
