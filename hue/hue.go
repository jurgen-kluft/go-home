package main

import (
	"fmt"
	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/stefanwichmann/go.hue"
	"time"
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
		client := pubsub.New()
		err := client.Connect("hue")
		if err == nil {

			fmt.Println("Connected to emitter")
			client.Subscribe(config.EmitterSensorLightChannelKey, "sensor/light/+")

			for {
				select {
				case msg := <-client.InMsgs:
					if msg.Topic() == "hue/config" {
						huelighting.config, err = config.HueConfigFromJSON(string(msg.Payload()))
					} else if msg.Topic() == "sensor/light/hue" {
						//huesensor, _ := config.SensorStateFromJSON(string(msg.Payload()))
					} else if msg.Topic() == "sensor/light/yee" {
						//yeesensor, _ := config.SensorStateFromJSON(string(msg.Payload()))
					}
					break
				case <-time.After(time.Second * 10):
					// do something if messages are taking too long
					// or if we haven't received enough state info.

					break
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
