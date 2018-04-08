package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/stefanwichmann/go.hue"
)

type instance struct {
	key    string
	config *config.HueConfig
}

func main() {
	huelighting := &instance{}

	bridgeIP := "10.0.0.72"
	bridgeKey := "4bWTg3R8JKwNDhMwZ8d6G4O9symsSi4iASzIJNgQ"

	untilBridgeFound := false
	for untilBridgeFound {
		bridges, _ := hue.DiscoverBridges(false)
		if len(bridges) > 0 {
			bridge := bridges[0] // Use the first bridge found

			waitUntilHueBridgeKeyPress := true
			for waitUntilHueBridgeKeyPress {
				err := bridge.CreateUser("go-home")
				if err != nil {
					fmt.Printf("HUE bridge connection failed: %v\n", err)
					time.Sleep(5 * time.Second)
				} else {
					bridgeIP = bridge.IpAddr
					bridgeKey = bridge.Username
					waitUntilHueBridgeKeyPress = false
					break
				}
			}
			fmt.Printf("HUE bridge connection succeeded => %+v\n", bridge)
			untilBridgeFound = false
		} else {
			fmt.Println("HUE bridge scanning ... (retry every 5 seconds)")
			time.Sleep(5 * time.Second)
		}
	}

	bridge := hue.NewBridge(bridgeIP, bridgeKey)

	lights, _ := bridge.GetAllLights()
	for _, l := range lights {
		fmt.Printf("Found hue light with name: %s\n", l.Name)
	}

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/hue/", "sensor/light/hue/"}
		subscribe := []string{"config/hue/", "sensor/light/hue/"}
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
