package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/jurgen-kluft/migateway"
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

			client.Register("config/xiaomi")
			client.Register("state/xiaomi")

			client.Register("state/xiaomi/gateway")
			client.Register("state/xiaomi/magnet")
			client.Register("state/xiaomi/motion")
			client.Register("state/xiaomi/plug")
			client.Register("state/xiaomi/switch")
			client.Register("state/xiaomi/dualwiredwallswitch")

			client.Subscribe("config/xiaomi")
			client.Subscribe("state/xiaomi")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/xiaomi" {
						xiaomi.config, err = config.XiaomiConfigFromJSON(string(msg.Payload()))
					} else if topic == "state/xiaomi" {
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							// TODO: Figure out what state to change on which device
							// Gateway color
							// Dualwiredwallswitch Channel 0/1 On/Off
							// Plug On/Off
							if state.Name == "" {

							}
						}
					} else if topic == "client/disconnected" {
						connected = false
					}

				case msg := <-xiaomi.aqara.StateMessages:

					// Push xiaomi gateway and device state changes onto pubsub channels

					switch msg.(type) {
					case migateway.GatewayStateChange:
						state := msg.(migateway.GatewayStateChange)
						name := "xiaomi.gateway." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddFloatSensor("illumination", state.To.Illumination)
						sensor.AddIntSensor("rgb", int64(state.To.RGB))
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							client.Publish("state/xiaomi/gateway", jsonstr)
						}

					case migateway.MagnetStateChange:
						state := msg.(migateway.MagnetStateChange)
						name := "xiaomi.magnet." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddFloatSensor("battery", state.To.Battery)
						sensor.AddBoolSensor("open", state.To.Opened)
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							client.Publish("state/xiaomi/magnet", jsonstr)
						}

					case migateway.MotionStateChange:
						state := msg.(migateway.MotionStateChange)
						name := "xiaomi.motion." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddTimeSlotSensor("battery", state.To.LastMotion, time.Now())
						sensor.AddBoolSensor("motion", state.To.HasMotion)
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							client.Publish("state/xiaomi/motion", jsonstr)
						}

					case migateway.PlugStateChange:
						state := msg.(migateway.PlugStateChange)
						name := "xiaomi.plug." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddBoolSensor("inuse", state.To.InUse)
						sensor.AddBoolSensor("ison", state.To.IsOn)
						sensor.AddIntSensor("loadvoltage", int64(state.To.LoadVoltage))
						sensor.AddIntSensor("loadpower", int64(state.To.LoadPower))
						sensor.AddIntSensor("powerconsumed", int64(state.To.PowerConsumed))
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							client.Publish("state/xiaomi/plug", jsonstr)
						}

					case migateway.SwitchStateChange:
						state := msg.(migateway.SwitchStateChange)
						name := "xiaomi.switch." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddIntSensor("battery", int64(state.To.Battery))
						sensor.AddValueSensor("click", state.To.Click.String())
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							client.Publish("state/xiaomi/switch", jsonstr)
						}

					case migateway.DualWiredWallSwitchStateChange:
						state := msg.(migateway.DualWiredWallSwitchStateChange)
						name := "xiaomi.dualwiredwallswitch." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddBoolSensor("channel0", state.To.Channel0On)
						sensor.AddBoolSensor("channel1", state.To.Channel1On)
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							client.Publish("state/xiaomi/dualwiredwallswitch", jsonstr)
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
