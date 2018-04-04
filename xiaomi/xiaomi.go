package main

import (
	"encoding/json"
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
	config *config.XiaomiConfig
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
			client.Subscribe("config/xiaomi/+")
			client.Subscribe("state/xiaomi/+")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/xiaomi" {
						xiaomi.config, err = config.XiaomiConfigFromJSON(string(msg.Payload()))
					} else if strings.HasPrefix(topic, "state/xiaomi") {
						var object string
						fmt.Sscanf(topic, "state/xiaomi/%s", object)

						// TODO: Figure out what state to change on which device
						// Gateway color
						// Dualwiredwallswitch Channel 0/1 On/Off
						// Plug On/Off

					} else if topic == "client/disconnected" {
						connected = false
					}

				case msg := <-xiaomi.aqara.StateMessages:

					// Push xiaomi gateway and device state changes onto pubsub channels

					switch msg.(type) {
					case migateway.GatewayStateChange:
						state := msg.(migateway.GatewayStateChange)
						name := "xiaomi.gateway." + state.ID
						var jsonmsg struct {
							Name         string  `json:"name"`
							Illumination float64 `json:"illumination"`
							RGB          uint32  `json:"rgb"`
						}
						jsonmsg.Name = name
						jsonmsg.Illumination = state.To.Illumination
						jsonmsg.RGB = state.To.RGB
						jsondata, err := json.Marshal(jsonmsg)
						if err == nil {
							client.Publish("state/xiaomi/gateway", string(jsondata))
						}

					case migateway.MagnetStateChange:
						state := msg.(migateway.MagnetStateChange)
						name := "xiaomi.magnet." + state.ID
						var jsonmsg struct {
							Name    string  `json:"name"`
							Battery float64 `json:"battery"`
							Open    bool    `json:"open"`
						}
						jsonmsg.Name = name
						jsonmsg.Battery = state.To.Battery
						jsonmsg.Open = state.To.Opened
						jsondata, err := json.Marshal(jsonmsg)
						if err == nil {
							client.Publish("state/xiaomi/magnet", string(jsondata))
						}

					case migateway.MotionStateChange:
						state := msg.(migateway.MotionStateChange)
						name := "xiaomi.motion." + state.ID
						var jsonmsg struct {
							Name   string    `json:"name"`
							Motion bool      `json:"motion"`
							Last   time.Time `json:"last"`
						}
						jsonmsg.Name = name
						jsonmsg.Motion = state.To.HasMotion
						jsonmsg.Last = state.To.LastMotion
						jsondata, err := json.Marshal(jsonmsg)
						if err == nil {
							client.Publish("state/xiaomi/motion", string(jsondata))
						}

					case migateway.PlugStateChange:
						state := msg.(migateway.PlugStateChange)
						name := "xiaomi.plug." + state.ID
						var jsonmsg struct {
							Name          string `json:"name"`
							InUse         bool   `json:"inuse"`
							IsOn          bool   `json:"ison"`
							LoadVoltage   uint32 `json:"loadvoltage"`
							LoadPower     uint32 `json:"loadpower"`
							PowerConsumed uint32 `json:"powerconsumed"`
						}
						jsonmsg.Name = name
						jsonmsg.InUse = state.To.InUse
						jsonmsg.IsOn = state.To.IsOn
						jsonmsg.LoadVoltage = state.To.LoadVoltage
						jsonmsg.LoadPower = state.To.LoadPower
						jsonmsg.PowerConsumed = state.To.PowerConsumed
						jsondata, err := json.Marshal(jsonmsg)
						if err == nil {
							client.Publish("state/xiaomi/plug", string(jsondata))
						}

					case migateway.SwitchStateChange:
						state := msg.(migateway.SwitchStateChange)
						name := "xiaomi.switch." + state.ID
						var jsonmsg struct {
							Name    string  `json:"name"`
							Battery float64 `json:"battery"`
							Click   string  `json:"click"`
						}
						jsonmsg.Name = name
						jsonmsg.Battery = state.To.Battery
						jsonmsg.Click = state.To.Click.String()
						jsondata, err := json.Marshal(jsonmsg)
						if err == nil {
							client.Publish("state/xiaomi/switch", string(jsondata))
						}

					case migateway.DualWiredWallSwitchStateChange:
						state := msg.(migateway.DualWiredWallSwitchStateChange)
						name := "xiaomi.dualwiredwallswitch." + state.ID
						var jsonmsg struct {
							Name     string `json:"name"`
							Channel0 bool   `json:"channel0"`
							Channel1 bool   `json:"channel1"`
						}
						jsonmsg.Name = name
						jsonmsg.Channel0 = state.To.Channel0On
						jsonmsg.Channel1 = state.To.Channel1On
						jsondata, err := json.Marshal(jsonmsg)
						if err == nil {
							client.Publish("state/xiaomi/dualwiredwallswitch", string(jsondata))
						}

					}

				// so that we do not have to poll anything and just push it on a emitter channel.
				// SensorState
				// {
				//   "name": "xiaomi.motion/switch/plug."A98C84E"
				//   "time": "Tue Apr 15 18:00:15 2014"
				//   ...
				// }

				case <-time.After(time.Second * 10):

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
