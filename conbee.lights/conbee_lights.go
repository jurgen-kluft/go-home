package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jurgen-kluft/go-home/conbee.lights/deconz"
	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
)

/*
This process will scan for events from Conbee, mainly sensors and will send those as
sensor states over NATS.
*/

type lightState struct {
	Name   string
	Conbee config.ConbeeLightGroup
}

type fullState struct {
	lights map[string]*lightState
}

func fullStateFromConfig(c *config.ConbeeLightsConfig) *fullState {
	full := &fullState{}
	full.lights = make(map[string]*lightState)
	for _, e := range c.Lights {
		full.lights[e.Name] = &lightState{Name: e.Name, Conbee: e}
	}
	return full
}

func main() {
	var cc *config.ConbeeLightsConfig = nil
	var nc *config.ConbeeLightsConfig = nil

	var fullState *fullState = nil
	var conbee *deconz.Client = nil

	var ctx context.Context

	for {
		var err error

		cc = nc

		register := []string{"config/request/", "config/conbee/lights/"}
		subscribe := []string{"config/conbee/lights/"}

		if cc != nil {
			subscribe = append(subscribe, cc.LightsIn...)
			fullState = fullStateFromConfig(cc)
			conbee = deconz.NewClient(&http.Client{}, cc.Host, cc.Port, cc.APIKey)
		}

		m := microservice.New("conbee/lights")
		m.RegisterAndSubscribe(register, subscribe)

		m.RegisterHandler("config/conbee/lights/", func(m *microservice.Service, topic string, msg []byte) bool {
			m.Logger.LogInfo(m.Name, "Received configuration, schedule restart")
			nc, err = config.ConbeeLightsConfigFromJSON(msg)
			if err != nil {
				m.Logger.LogError(m.Name, err.Error())
			} else {
				cc = nil
				return false
			}
			return true
		})

		m.RegisterHandler("state/light/automation/", func(m *microservice.Service, topic string, msg []byte) bool {
			sensor, err := config.SensorStateFromJSON(msg)
			if err == nil {
				m.Logger.LogInfo(m.Name, "received state")
				lightname := sensor.Name
				if lightname != "" && conbee != nil {
					if topic == "state/light/ahk/" || topic == "state/light/automation/" {
						lstate, exist := fullState.lights[lightname]
						if exist {
							sensor.ExecValueAttr("power", func(power string) {
								if power == "on" {
									fmt.Println(lightname + " turning On")
									conbee.SetGroupStateFromJSON(ctx, lstate.Conbee.Group, lstate.Conbee.On)
								} else if power == "off" {
									fmt.Println(lightname + " turning Off")
									conbee.SetGroupStateFromJSON(ctx, lstate.Conbee.Group, lstate.Conbee.Off)
								}
							})
						} else {
							fmt.Println(lightname + " doesn't exist")
						}
					} else if topic == "state/light/flux/" {
						lstate, exist := fullState.lights[lightname]
						if exist {
							ct := sensor.GetFloatAttr("ct", 500)
							bri := sensor.GetFloatAttr("bri", 500)
							conbee.SetGroupStateFromJSON(ctx, lstate.Conbee.Group, fmt.Sprintf(lstate.Conbee.CT, ct, bri))
						} else {
							fmt.Println(lightname + " doesn't exist")
						}
					} else {
						fmt.Println(lightname + " doesn't exist")
					}
				}
			}
			return true
		})

		tickCount := 0
		m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
			if (tickCount % 30) == 0 {
				if nc == nil {
					m.Logger.LogInfo(m.Name, "Requesting configuration..")
					m.Pubsub.PublishStr("config/request/", m.Name)
				}
			}
			tickCount++
			return true
		})

		m.Loop()

		// Sleep for a while before restarting
		m.Logger.LogInfo(m.Name, "Waiting 10 seconds before re-starting..")
		time.Sleep(10 * time.Second)

		m.Logger.LogInfo(m.Name, "Re-start..")
	}

}
