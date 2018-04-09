package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/stefanwichmann/go.hue"
)

type instance struct {
	key    string
	config *config.HueConfig
	lights map[string]*hue.Light
}

func main() {
	huelighting := &instance{}

	bridgeIP := "10.0.0.72"
	bridgeKey := "4bWTg3R8JKwNDhMwZ8d6G4O9symsSi4iASzIJNgQ"

	logger := logpkg.New("hue")
	logger.AddEntry("emitter")
	logger.AddEntry("hue")

	for {
		untilBridgeFound := false
		for untilBridgeFound {
			bridges, _ := hue.DiscoverBridges(false)
			if len(bridges) > 0 {
				bridge := bridges[0] // Use the first bridge found

				waitUntilHueBridgeKeyPress := true
				for waitUntilHueBridgeKeyPress {
					err := bridge.CreateUser("go-home")
					if err != nil {
						logger.LogError("hue", fmt.Sprintf("HUE bridge connection failed: %v", err))
						time.Sleep(5 * time.Second)
					} else {
						bridgeIP = bridge.IpAddr
						bridgeKey = bridge.Username
						waitUntilHueBridgeKeyPress = false
						break
					}
				}
				logger.LogInfo("hue", fmt.Sprintf("HUE bridge connection succeeded => %+v", bridge))
				untilBridgeFound = false
			} else {
				logger.LogInfo("hue", "HUE bridge scanning ... (retry every 5 seconds)")
				time.Sleep(5 * time.Second)
			}
		}

		bridge := hue.NewBridge(bridgeIP, bridgeKey)
		lights, err := bridge.GetAllLights()
		if err == nil {
			for _, light := range lights {
				huelighting.lights[light.Name] = light
			}
			break
		}

		logger.LogError("hue", fmt.Sprintf("HUE bridge does not have any lights: %v", err))
	}

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/hue/", "sensor/light/hue/", "state/light/hue/"}
		subscribe := []string{"config/hue/", "sensor/light/hue/", "state/light/hue/"}
		err := client.Connect("hue", register, subscribe)
		if err == nil {
			fmt.Println("Connected to emitter")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/hue/" {
						huelighting.config, err = config.HueConfigFromJSON(string(msg.Payload()))
					} else if topic == "sensor/light/hue/" {
						//huesensor, _ := config.SensorStateFromJSON(string(msg.Payload()))
					} else if topic == "client/disconnected/" {
						connected = false
					} else if topic == "state/light/hue/" {
						huesensor, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							logger.LogInfo("hue", "received state")
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
					}

				case <-time.After(time.Second * 10):

				}
			}
		}

		if err != nil {
			fmt.Println("Error: " + err.Error())
			time.Sleep(5 * time.Second)
		}

	}
}
