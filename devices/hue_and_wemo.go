package devices

import (
	"time"
)

func NewHueAndWemo() Devices {
	return &HueWemoDevices{}
}

type HueWemoDevices struct {
}

func (d *HueWemoDevices) GetAllLights() []string {
	return nil
}

func (d *HueWemoDevices) GetLight(id string) Light {
	return nil
}

func (d *HueWemoDevices) GetAllSwitches() []string {
	return nil
}

func (d *HueWemoDevices) GetSwitch(id string) Switch {
	return nil
}

type HueLight struct {
}

func (hue *HueLight) TurnOn() {

}

func (hue *HueLight) TurnOff() {

}

func (hue *HueLight) SetScene(name string) {

}

func (hue *HueLight) SetBrightnessRange(from, to float64) {

}

func (hue *HueLight) SetFadeDuration(duration time.Duration) {

}

type WemoSwitch struct {
}

func (wemo *WemoSwitch) TurnOn() {

}

func (wemo *WemoSwitch) TurnOff() {

}
