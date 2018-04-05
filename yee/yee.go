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
		client := pubsub.New("tcp://10.0.0.22:8080")
		register := []string{"config/yee/", "sensor/light/yee/"}
		subscribe := []string{"config/yee/", "sensor/light/yee/"}
		err := client.Connect("yee", register, subscribe)
		if err == nil {
			fmt.Println("Connected to emitter")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/yee/" {
						yeelighting.config, err = config.YeeConfigFromJSON(string(msg.Payload()))
						yeelighting.lamps = map[string]*yee.Yeelight{}
						for _, lamp := range yeelighting.config.Lights {
							yeelighting.lamps[lamp.Name] = yee.New(lamp.IP, lamp.Port)
						}
					} else if topic == "sensor/light/yee/" {
						yeesensor, _ := config.SensorStateFromJSON(string(msg.Payload()))
						lampname := yeesensor.GetValueAttr("name", "")
						if lampname != "" {
							lamp, exists := yeelighting.lamps[lampname]
							if exists {
								power := yeesensor.GetValueAttr("power", "")
								if power != "" {
									if power == "on" {
										lamp.On()
									} else if power == "off" {
										lamp.Off()
									}
								}
								if power == "on" {
									ct := yeesensor.GetFloatAttr("ct", -1.0)
									if ct != -1.0 {
										lamp.SetCtAbx(fmt.Sprintf("%f", ct), "smooth", "500")
									}
									bri := yeesensor.GetFloatAttr("bri", -1.0)
									if ct != -1.0 {
										lamp.SetBright(fmt.Sprintf("%f", bri), "smooth", "500")
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
