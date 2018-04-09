package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	log "github.com/jurgen-kluft/go-home/logging"
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

	aqara, err := migateway.NewAqaraManager(nil)
	if err != nil {
		panic(err)
	}

	xiaomi.aqara = aqara
	xiaomi.aqara.SetAESKey(xiaomi.key)

	logger := log.New("xiaomi")
	logger.AddEntry("emitter")
	logger.AddEntry("xiaomi")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/xiaomi/", "state/xiaomi/", "state/xiaomi/gateway/", "state/xiaomi/magnet/", "state/xiaomi/motion/", "state/xiaomi/plug/", "state/xiaomi/switch/", "state/xiaomi/dualwiredwallswitch/"}
		subscribe := []string{"config/xiaomi/", "state/xiaomi/"}
		err := client.Connect("xiaomi", register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")
			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/xiaomi/" {
						logger.LogInfo("xiaomi", "received configuration")
						xiaomi.config, err = config.XiaomiConfigFromJSON(string(msg.Payload()))
					} else if topic == "state/xiaomi/" {
						logger.LogInfo("xiaomi", "received state")
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							// TODO: Figure out what state to change on which device
							// Gateway color
							// Dualwiredwallswitch Channel 0/1 On/Off
							// Plug On/Off
							if state.Name == "" {

							}
						}
					} else if topic == "client/disconnected/" {
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}

				case msg := <-xiaomi.aqara.StateMessages:

					// Push xiaomi gateway and device state changes onto pubsub channels
					// cfmt.Printf("STATE message received %v (type: %s)\n", msg, reflect.TypeOf(msg))

					switch msg.(type) {
					case *migateway.GatewayStateChange:
						state := msg.(*migateway.GatewayStateChange)
						name := "xiaomi.gateway." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddFloatAttr("illumination", state.To.Illumination)
						sensor.AddIntAttr("rgb", int64(state.To.RGB))
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							fmt.Println(jsonstr)
							client.Publish("state/xiaomi/gateway/", jsonstr)
						}

					case *migateway.MagnetStateChange:
						state := msg.(*migateway.MagnetStateChange)
						name := "xiaomi.magnet." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddBoolAttr("open", state.To.Opened)
						sensor.AddFloatAttr("battery", float64(state.To.Battery))
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							fmt.Println(jsonstr)
							client.Publish("state/xiaomi/magnet/", jsonstr)
						}

					case *migateway.MotionStateChange:
						fmt.Println("STATE motion change message received")

						state := msg.(*migateway.MotionStateChange)
						name := "xiaomi.motion." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddTimeWndAttr("lastmotion", state.To.LastMotion, time.Now())
						sensor.AddBoolAttr("motion", state.To.HasMotion)
						sensor.AddFloatAttr("battery", float64(state.To.Battery))
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							fmt.Println(jsonstr)
							client.Publish("state/xiaomi/motion/", jsonstr)
						}

					case *migateway.PlugStateChange:
						state := msg.(*migateway.PlugStateChange)
						name := "xiaomi.plug." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddBoolAttr("inuse", state.To.InUse)
						sensor.AddBoolAttr("ison", state.To.IsOn)
						sensor.AddIntAttr("loadvoltage", int64(state.To.LoadVoltage))
						sensor.AddIntAttr("loadpower", int64(state.To.LoadPower))
						sensor.AddIntAttr("powerconsumed", int64(state.To.PowerConsumed))
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							fmt.Println(jsonstr)
							client.Publish("state/xiaomi/plug/", jsonstr)
						}

					case *migateway.SwitchStateChange:
						state := msg.(*migateway.SwitchStateChange)
						name := "xiaomi.switch." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddStringAttr("click", state.To.Click.String())
						sensor.AddFloatAttr("battery", float64(state.To.Battery))
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							fmt.Println(jsonstr)
							client.Publish("state/xiaomi/switch/", jsonstr)
						}

					case *migateway.DualWiredWallSwitchStateChange:
						state := msg.(*migateway.DualWiredWallSwitchStateChange)
						name := "xiaomi.dualwiredwallswitch." + state.ID
						sensor := config.NewSensorState(name)
						sensor.AddBoolAttr("channel0", state.To.Channel0On)
						sensor.AddBoolAttr("channel1", state.To.Channel1On)
						jsonstr, err := sensor.ToJSON()
						if err == nil {
							fmt.Println(jsonstr)
							client.Publish("state/xiaomi/dualwiredwallswitch/", jsonstr)
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
		}

		if err != nil {
			logger.LogError("xiaomi", err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}

func (p *instance) handle(object string) {

}
