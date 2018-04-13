package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/stefanwichmann/go.hue"
)

// HueLighting holds all necessary information to control the HUE bridge and lights
type HueLighting struct {
	key    string
	config *config.HueConfig
	bridge *hue.Bridge
	lights map[string]*hue.Light
	log    *logpkg.Logger
}

// New creates a new instance of huelighting instance
func New() *HueLighting {
	huelighting := &HueLighting{}
	huelighting.lights = map[string]*hue.Light{}
	return huelighting
}

func (h *HueLighting) initializeBridgeConnection() (err error) {
	bridgeIP := h.config.Host
	bridgeKey := h.config.Key

	if bridgeIP == "" || bridgeKey == "" {
		bridgeFound := false
		for !bridgeFound {
			bridges, err := hue.DiscoverBridges(false)
			if err == nil && len(bridges) > 0 {
				bridge := bridges[0] // Use the first bridge found

				hueBridgeKeyPressRetry := 60
				for hueBridgeKeyPressRetry > 0 {
					err := bridge.CreateUser("go-home")
					if err != nil {
						h.log.LogError("hue", fmt.Sprintf("HUE bridge connection failed: %v", err))
						time.Sleep(5 * time.Second)
						hueBridgeKeyPressRetry--
					} else {
						bridgeIP = bridge.IpAddr
						bridgeKey = bridge.Username
						bridgeFound = true
						break
					}
				}

				if err == nil && bridgeFound {
					h.log.LogInfo("hue", fmt.Sprintf("HUE bridge connection succeeded => %+v", bridge))
				}
			}

			if err != nil {
				h.log.LogInfo("hue", "HUE bridge scanning ... (retry every 5 seconds)")
				time.Sleep(5 * time.Second)
			}
		}
	}

	h.bridge = hue.NewBridge(bridgeIP, bridgeKey)

	var lights []*hue.Light
	lights, err = h.bridge.GetAllLights()
	if err == nil {
		for _, light := range lights {
			h.lights[light.Name] = light
		}
	}

	return err
}

func main() {
	huelighting := New()

	huelighting.log = logpkg.New("hue")
	huelighting.log.AddEntry("emitter")
	huelighting.log.AddEntry("hue")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/hue/", "state/sensor/hue/", "state/light/hue/"}
		subscribe := []string{"config/hue/", "state/light/hue/", "sensor/light/hue/"}
		err := client.Connect("hue", register, subscribe)
		if err == nil {
			huelighting.log.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/hue/" {
						config, err := config.HueConfigFromJSON(string(msg.Payload()))
						if err == nil {
							huelighting.log.LogInfo("hue", "received configuration")
							huelighting.config = config
							err = huelighting.initializeBridgeConnection()
							if err != nil {
								huelighting.log.LogError("hue", err.Error())
							}
						} else {
							huelighting.log.LogError("hue", err.Error())
						}
					} else if topic == "sensor/light/hue/" {
						//huesensor, _ := config.SensorStateFromJSON(string(msg.Payload()))
					} else if topic == "client/disconnected/" {
						huelighting.log.LogInfo("emitter", "disconnected")
						connected = false
					} else if topic == "state/light/hue/" {
						if huelighting.config != nil {
							huesensor, err := config.SensorStateFromJSON(string(msg.Payload()))
							if err == nil {
								huelighting.log.LogInfo("hue", "received state")
							}
							lightname := huesensor.GetValueAttr("name", "")
							if lightname != "" {
								light, exists := huelighting.lights[lightname]
								if exists {
									huesensor.ExecValueAttr("power", func(power string) {
										if power == "on" {
											light.On()
										} else if power == "off" {
											light.Off()
										}
									})
								}
							}
						} else {
							huelighting.log.LogError("hue", fmt.Sprintf("error, receiving message on channel %s but we haven't received a configuration", topic))
						}
					}

				case <-time.After(time.Second * 10):

				}
			}
		}

		if err != nil {
			huelighting.log.LogError("hue", err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}
