package main

import (
	"time"

	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
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
	for _, dev := range x.config.Motion {
		if dev.ID == ID {
			return dev.Name
		}
	}
	return ID
}
func (x *xiaomi) GetNameOfMagnetSensor(ID string) string {
	for _, dev := range x.config.Magnet {
		if dev.ID == ID {
			return dev.Name
		}
	}
	return ID
}
func (x *xiaomi) GetNameOfPlug(ID string) string {
	for _, dev := range x.config.Plug {
		if dev.ID == ID {
			return dev.Name
		}
	}
	return ID
}
func (x *xiaomi) GetPlugByName(name string) *migateway.Plug {
	for _, dev := range x.config.Plug {
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
	for _, dev := range x.config.Switch {
		if dev.ID == ID {
			return dev.Name
		}
	}
	return ID
}

func (x *xiaomi) GetDualWiredWallSwitchByName(name string) *migateway.DualWiredWallSwitch {
	for _, device := range x.config.Switch {
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

	register := []string{"state/xiaomi/", "config/request/"}
	subscribe := []string{"config/xiaomi/"}

	m := microservice.New("xiaomi")
	m.RegisterAndSubscribe(register, subscribe)

	m.RegisterHandler("config/xiaomi/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo("xiaomi", "received configuration")
		var err error
		xiaomi.config, err = config.XiaomiConfigFromJSON(msg)
		if err == nil {
			err = xiaomi.aqara.Start(nil)
			if err == nil {
				xiaomi.aqara.SetAESKey(xiaomi.config.Key.String)
			} else {
				m.Logger.LogError(m.Name, err.Error())
			}
		} else {
			m.Logger.LogError(m.Name, err.Error())
		}
		return true
	})

	m.RegisterHandler("state/xiaomi/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received state")
		state, err := config.SensorStateFromJSON(msg)
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
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if tickCount%5 == 0 {
			if xiaomi.config == nil {
				m.Pubsub.PublishStr("config/request/", m.Name)
			}
		}
		tickCount++
		return true
	})

	// Xiaomi Aqara Gateway and device state changes onto the process messages channel of the micro-service
	go func() {
		for true {
			select {
			case msg := <-xiaomi.aqara.StateMessages:
				//fmt.Printf("STATE message received %v (type: %s)\n", msg, reflect.TypeOf(msg))

				switch msg.(type) {
				case *migateway.GatewayStateChange:
					state := msg.(*migateway.GatewayStateChange)
					name := "gateway"
					sensor := config.NewSensorState(name)
					sensor.AddFloatAttr("illumination", state.To.Illumination)
					sensor.AddIntAttr("rgb", int64(state.To.RGB))
					jsondata, err := sensor.ToJSON()
					if err == nil {
						//fmt.Println(jsonstr)
						msg := &microservice.Message{Topic: "state/xiaomi/", Payload: jsondata}
						m.ProcessMessages <- msg
					}

				case *migateway.MagnetStateChange:
					state := msg.(*migateway.MagnetStateChange)
					name := xiaomi.GetNameOfMagnetSensor(state.ID)
					sensor := config.NewSensorState(name)
					sensor.AddBoolAttr("open", state.To.Opened)
					sensor.AddFloatAttr("battery", float64(state.To.Battery))
					jsondata, err := sensor.ToJSON()
					if err == nil {
						//fmt.Println(jsonstr)
						msg := &microservice.Message{Topic: "state/xiaomi/", Payload: jsondata}
						m.ProcessMessages <- msg
					}

				case *migateway.MotionStateChange:
					state := msg.(*migateway.MotionStateChange)
					name := xiaomi.GetNameOfMotionSensor(state.ID)
					sensor := config.NewSensorState(name)
					sensor.AddTimeWndAttr("lastmotion", state.To.LastMotion, time.Now())
					sensor.AddBoolAttr("motion", state.To.HasMotion)
					sensor.AddFloatAttr("battery", float64(state.To.Battery))
					jsondata, err := sensor.ToJSON()
					if err == nil {
						//fmt.Println(jsonstr)
						msg := &microservice.Message{Topic: "state/xiaomi/", Payload: jsondata}
						m.ProcessMessages <- msg
					}

				case *migateway.PlugStateChange:
					state := msg.(*migateway.PlugStateChange)
					name := xiaomi.GetNameOfPlug(state.ID)
					sensor := config.NewSensorState(name)
					sensor.AddBoolAttr("inuse", state.To.InUse)
					sensor.AddBoolAttr("ison", state.To.IsOn)
					sensor.AddIntAttr("loadvoltage", int64(state.To.LoadVoltage))
					sensor.AddIntAttr("loadpower", int64(state.To.LoadPower))
					sensor.AddIntAttr("powerconsumed", int64(state.To.PowerConsumed))
					jsondata, err := sensor.ToJSON()
					if err == nil {
						//fmt.Println(jsonstr)
						msg := &microservice.Message{Topic: "state/xiaomi/", Payload: jsondata}
						m.ProcessMessages <- msg
					}

				case *migateway.SwitchStateChange:
					state := msg.(*migateway.SwitchStateChange)
					name := xiaomi.GetNameOfSwitch(state.ID)
					sensor := config.NewSensorState(name)
					sensor.AddStringAttr("click", state.To.Click.String())
					sensor.AddFloatAttr("battery", float64(state.To.Battery))
					jsondata, err := sensor.ToJSON()
					if err == nil {
						//fmt.Println(jsonstr)
						msg := &microservice.Message{Topic: "state/xiaomi/", Payload: jsondata}
						m.ProcessMessages <- msg
					}

				case *migateway.DualWiredWallSwitchStateChange:
					state := msg.(*migateway.DualWiredWallSwitchStateChange)
					name := xiaomi.GetNameOfSwitch(state.ID)
					sensor := config.NewSensorState(name)
					sensor.AddBoolAttr("channel0", state.To.Channel0On)
					sensor.AddBoolAttr("channel1", state.To.Channel1On)
					jsondata, err := sensor.ToJSON()
					if err == nil {
						//fmt.Println(jsonstr)
						msg := &microservice.Message{Topic: "state/xiaomi/", Payload: jsondata}
						m.ProcessMessages <- msg
					}
				}
			}
		}
	}()

	m.Loop()
}
