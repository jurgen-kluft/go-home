package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
	huegroups "github.com/jurgen-kluft/go-hue/groups"
	huelights "github.com/jurgen-kluft/go-hue/lights"
)

// hue holds all necessary information to control the HUE bridge and lights
type huecontext struct {
	key        string
	config     *config.HueConfig
	lights     *huelights.Lights
	groups     *huegroups.Groups
	group      huegroups.Group
	light      map[string]huelights.Light
	log        *logpkg.Logger
	CT         float64
	BRI        float64
	lightState huelights.State
}

// New creates a new instance of hue instance
func New() *huecontext {
	hue := &huecontext{}
	hue.light = map[string]huelights.Light{}
	hue.lightState.CT = new(uint16)
	hue.lightState.Bri = new(uint8)
	return hue
}

func (hue *huecontext) initialize() (err error) {
	bridgeIP := hue.config.Host
	bridgeKey := hue.config.Key

	if bridgeIP == "" || bridgeKey == "" {
		return fmt.Errorf("please specify a host and key in the hue configuration")
	}

	// Obtain information of all lights
	hue.lights = huelights.New(bridgeIP, bridgeKey)
	lights, err := hue.lights.GetAllLights()
	if err == nil {
		for _, light := range lights {
			hue.light[light.Name] = light
		}
	} else {
		hue.log.LogError("hue", err.Error())
	}

	// Construct group interface and create a new
	// group called 'All' so that we can update
	// CT and BRI with one call.
	hue.groups = huegroups.New(bridgeIP, bridgeKey)
	hue.group = huegroups.Group{}
	hue.group.Name = "All"
	response, err := hue.groups.CreateGroup(hue.group)
	if err == nil && len(response) > 0 {
		idvar, exists := response[0].Success["id"]
		if exists {
			id := idvar.(int)
			hue.group.ID = id
		}
	} else {
		hue.log.LogError("hue", err.Error())
	}
	return err
}

func (hue *huecontext) flux() {
	*hue.lightState.CT = uint16(hue.CT)
	*hue.lightState.Bri = uint8(hue.BRI)
	hue.groups.SetGroupState(hue.group.ID, hue.lightState)
}

func main() {
	hue := New()

	hue.log = logpkg.New("hue")
	hue.log.AddEntry("emitter")
	hue.log.AddEntry("hue")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/hue/", "state/sensor/hue/", "sensor/light/hue/"}
		subscribe := []string{"config/hue/", "state/sensor/hue/", "sensor/light/hue/"}
		err := client.Connect("hue", register, subscribe)
		if err == nil {
			hue.log.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/hue/" {
						config, err := config.HueConfigFromJSON(string(msg.Payload()))
						if err == nil {
							hue.log.LogInfo("hue", "received configuration")
							hue.config = config
							err = hue.initialize()
							if err != nil {
								hue.log.LogError("hue", err.Error())
							}
						} else {
							hue.log.LogError("hue", err.Error())
						}
					} else if topic == "client/disconnected/" {
						hue.log.LogInfo("emitter", "disconnected")
						connected = false
					} else if topic == "state/sensor/hue/" {
						huesensor, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							hue.log.LogInfo("hue", "received flux")
							hue.CT = huesensor.GetFloatAttr("CT", 325.0)
							hue.BRI = huesensor.GetFloatAttr("BRI", 128.0)
						} else {
							hue.log.LogError("hue", err.Error())
						}
					} else if topic == "state/light/hue/" {
						if hue.config != nil {
							huesensor, err := config.SensorStateFromJSON(string(msg.Payload()))
							if err == nil {
								hue.log.LogInfo("hue", "received state")
								lightname := huesensor.GetValueAttr("name", "")
								if lightname != "" {
									light, exists := hue.light[lightname]
									if exists {
										huesensor.ExecValueAttr("power", func(power string) {
											if power == "on" {
												hue.lightState.On = new(bool)
												*hue.lightState.On = true
												*hue.lightState.CT = uint16(hue.CT)
												*hue.lightState.Bri = uint8(hue.BRI)
												hue.lights.SetLightState(light.ID, hue.lightState)
												hue.lightState.On = nil
											} else if power == "off" {
												*hue.lightState.On = false
												hue.lights.SetLightState(light.ID, hue.lightState)
												hue.lightState.On = nil
											}
										})
									}
								}
							} else {
								hue.log.LogError("hue", err.Error())
							}
						} else {
							hue.log.LogError("hue", fmt.Sprintf("error, receiving message on channel %s but we haven't received a configuration", topic))
						}
					}

				case <-time.After(time.Second * 60):
					hue.flux()
				}
			}
		}

		if err != nil {
			hue.log.LogError("hue", err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}
