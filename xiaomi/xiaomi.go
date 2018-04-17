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

type xiaomi struct {
	config *config.XiaomiConfig
	aqara  *migateway.AqaraManager
}

func (x *xiaomi) GetNameOfMotionSensor(ID string) string {
	for _, dev := range x.config.Motions {
		if dev.ID == ID {
			return dev.Name
		}
	}
	return ID
}
func (x *xiaomi) GetNameOfMagnetSensor(ID string) string {
	for _, dev := range x.config.Magnets {
		if dev.ID == ID {
			return dev.Name
		}
	}
	return ID
}
func (x *xiaomi) GetNameOfPlug(ID string) string {
	for _, dev := range x.config.Plugs {
		if dev.ID == ID {
			return dev.Name
		}
	}
	return ID
}
func (x *xiaomi) GetPlugByName(name string) *migateway.Plug {
	for _, dev := range x.config.Plugs {
		if dev.Name == name {
			for ID, plug := range x.aqara.Plugs {
				if dev.ID == ID {
					return plug
				}
			}
		}
	}
	return nil
}

func (x *xiaomi) GetNameOfSwitch(ID string) string {
	for _, dev := range x.config.Switches {
		if dev.ID == ID {
			return dev.Name
		}
	}
	return ID
}

func (x *xiaomi) GetDualWiredWallSwitchByName(name string) *migateway.DualWiredWallSwitch {
	for _, device := range x.config.Switches {
		if device.Name == name {
			for ID, hwdev := range x.aqara.DualWiredSwitches {
				if device.ID == ID {
					return hwdev
				}
			}
		}
	}
	return nil
}

func main() {
	xiaomi := &xiaomi{}
	xiaomi.aqara = migateway.NewAqaraManager()

	logger := log.New("xiaomi")
	logger.AddEntry("emitter")
	logger.AddEntry("xiaomi")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/xiaomi/", "state/xiaomi/"}
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
						if err == nil {
							err = xiaomi.aqara.Start(nil)
							if err == nil {
								xiaomi.aqara.SetAESKey(xiaomi.config.Key.String)
							} else {
								logger.LogError("xiaomi", err.Error())
							}
						} else {
							logger.LogError("xiaomi", err.Error())
						}
					} else if topic == "state/xiaomi/" {
						logger.LogInfo("xiaomi", "received state")
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							if state.Name == "gateway" {
								if state.GetValueAttr("light", "none") == "on" {
									xiaomi.aqara.GateWay.TurnOn()
								} else if state.GetValueAttr("light", "none") == "off" {
									xiaomi.aqara.GateWay.TurnOff()
								}
							}
							plug := xiaomi.GetPlugByName(state.Name)
							if plug != nil {
								if state.GetValueAttr("power", "none") == "on" {
									plug.TurnOn()
								} else if state.GetValueAttr("power", "none") == "off" {
									plug.TurnOff()
								} else if state.GetValueAttr("power", "none") == "toggle" {
									plug.Toggle()
								}
							} else {
								dwwswitch := xiaomi.GetDualWiredWallSwitchByName(state.Name)
								if dwwswitch != nil {
									if state.GetValueAttr("switch0", "none") == "on" {
										dwwswitch.TurnOnChannel0()
									} else if state.GetValueAttr("switch0", "none") == "off" {
										dwwswitch.TurnOffChannel0()
									}
									if state.GetValueAttr("switch1", "none") == "on" {
										dwwswitch.TurnOnChannel1()
									} else if state.GetValueAttr("switch1", "none") == "off" {
										dwwswitch.TurnOffChannel1()
									}
								}
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
							client.Publish("state/xiaomi/", jsonstr)
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
							client.Publish("state/xiaomi/", jsonstr)
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
							client.Publish("state/xiaomi/", jsonstr)
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
							client.Publish("state/xiaomi/", jsonstr)
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
							client.Publish("state/xiaomi/", jsonstr)
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
							client.Publish("state/xiaomi/", jsonstr)
						}

					}

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
