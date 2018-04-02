package main

import (
	"fmt"
	"strings"
	"time"

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
	aqara  *migateway.AqaraManager
}

func main() {
	xiaomi := &instance{}

	//gatewayIP := "10.0.0.78"
	xiaomi.key = "3C8FA0275CAF4567"

	aqara, err := migateway.NewAqaraManager(nil)
	if err != nil {
		panic(err)
	}

	xiaomi.aqara = aqara
	xiaomi.aqara.SetAESKey(xiaomi.key)

	for {
		client := pubsub.New()
		err := client.Connect("xiaomi")
		if err == nil {

			fmt.Println("Connected to emitter")
			client.Subscribe(config.XiaomiStateChannelKey, "xiaomi/+")

			for {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "xiaomi/config" {
						xiaomi.config, err = config.XiaomiConfigFromJSON(string(msg.Payload()))
					} else if strings.HasPrefix(topic, "xiaomi/state") {
						// state object, json object
						var object string
						fmt.Sscanf(topic, "xiaomi/state/%s", object)

					}
					break

				// We would like to receive state messages from the Gateway on a channel here
				// so that we do not have to poll anything and just push it on a emitter channel.
				// SensorState
				// {
				//   "domain": "xiaomi"
				//   "product": "motion" / "switch" / "plug"
				//   "name": "A98C84E"
				//   "type": "string"
				//   "value": "on/off"
				//   "time": "Tue Apr 15 18:00:15 2014"
				// }

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

func (p *instance) handle(object string) {

}
