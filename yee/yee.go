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

func (x *instance) poweron(name string) {
	lamp, exists := x.lamps[name]
	if exists {
		lamp.On()
	}
}
func (x *instance) poweroff(name string) {
	lamp, exists := x.lamps[name]
	if exists {
		lamp.Off()
	}
}

func main() {
	yeelighting := &instance{}
	yeelighting.lamps = map[string]*yee.Yeelight{}

	// yeelighting.lamps["Front door hall light"] = yee.New("10.0.0.113", "55443")
	// yeelighting.poweroff("Front door hall light")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
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
						yeesensor, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							lampname := yeesensor.GetValueAttr("name", "")
							if lampname != "" {
								lamp, exists := yeelighting.lamps[lampname]
								if exists {
									yeesensor.ExecValueAttr("power", func(power string) {
										if power == "on" {
											lamp.On()
										} else if power == "off" {
											lamp.Off()
										}
									})
									yeesensor.ExecFloatAttr("ct", func(ct float64) {
										lamp.SetCtAbx(fmt.Sprintf("%f", ct), "smooth", "500")
									})
									yeesensor.ExecFloatAttr("bri", func(bri float64) {
										lamp.SetBright(fmt.Sprintf("%f", bri), "smooth", "500")
									})
								}
							} else if lampname == "all" {
								for _, lamp := range yeelighting.lamps {
									yeesensor.ExecFloatAttr("ct", func(ct float64) {
										lamp.SetCtAbx(fmt.Sprintf("%f", ct), "smooth", "500")
									})
									yeesensor.ExecFloatAttr("bri", func(bri float64) {
										lamp.SetBright(fmt.Sprintf("%f", bri), "smooth", "500")
									})
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
