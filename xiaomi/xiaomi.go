package main

import (
	"fmt"
	"time"

	"github.com/bingbaba/tool/color"
	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/xuebing1110/migateway"
)

// Features:
// - Gateway Light: Turn On/Off, Change color/brightness
// - Gateway Sound: Play sound

// Turn On/Off:
// - WiredDualWallSwitch(es)
// - Electric Power Plug(s)

// Publish state of:
// - Motion Sensor(s)
// - Wireless Switch(es)
// - WiredDualWallSwitch(es)
// - Electric Power Plug(s)

type instance struct {
	key    string
	config *config.HueConfig
}

func main() {
	huelighting := &instance{}

	//gatewayIP := "10.0.0.78"
	gatewayKey := "3C8FA0275CAF4567"

	manager, err := migateway.NewAqaraManager(nil)
	if err != nil {
		panic(err)
	}

	manager.SetAESKey(gatewayKey)

	gateway := manager.GateWay
	for _, color := range color.COLOR_ALL {
		err = gateway.ChangeColor(color)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}

	err = gateway.Flashing(color.COLOR_RED)
	if err != nil {
		panic(err)
	}

	for {
		client := pubsub.New()
		err := client.Connect("xiaomi")
		if err == nil {

			fmt.Println("Connected to emitter")
			client.Subscribe(config.XiaomiStateChannelKey, "xiaomi/+")

			for {
				select {
				case msg := <-client.InMsgs:
					if msg.Topic() == "xiaomi/config" {
						huelighting.config, err = config.HueConfigFromJSON(string(msg.Payload()))
					} else if msg.Topic() == "xiaomi/state" {
						// state object, json object
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
