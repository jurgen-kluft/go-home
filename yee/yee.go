package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
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

	logger := logpkg.New("yee")
	logger.AddEntry("emitter")
	logger.AddEntry("yee")

	for {
		client := pubsub.New(config.EmitterIOCfg)
		register := []string{"config/yee/", "state/sensor/yee/", "state/light/yee/"}
		subscribe := []string{"config/yee/", "state/sensor/yee/", "state/light/yee/"}
		err := client.Connect("yee", register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/yee/" {
						logger.LogInfo("yee", "received configuration")
						yeelighting.config, err = config.YeeConfigFromJSON(string(msg.Payload()))
						yeelighting.lamps = map[string]*yee.Yeelight{}
						for _, lamp := range yeelighting.config.Lights {
							yeelighting.lamps[lamp.Name] = yee.New(lamp.IP, lamp.Port)
						}
					} else if topic == "state/light/yee/" {
						yeesensor, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							logger.LogInfo("yee", "received state")

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
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}

				case <-time.After(time.Second * 10):

				}
			}
		}

		if err != nil {
			logger.LogError("yee", err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}
