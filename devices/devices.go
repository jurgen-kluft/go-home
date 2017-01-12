package devices

import (
	"time"
)

/*

We will abstract lights and switches and have the API work by name

Lights
Switches

Examples:
    alllights := devices.GetAllLights()

    livingroom := devices.GetLights("Livingroom")
    livingroom.SetScene("Winter:Evening")
    livingroom.SetBrightnessRange(0.1, 0.8)
    livingroom.SetFadeDuration(5 * time.Minutes)
    livingroom.TurnOn()
    livingroom.TurnOff()
    livingroom.SetScene("Winter:Noon")

    christmastree := devices.GetSwitches("Christmas")
    christmastree.TurnOff()
    christmastree.TurnOn()

*/
type Light interface {
	TurnOn()
	TurnOff()
	SetScene(name string)
	SetBrightnessRange(from, to float64)
	SetFadeDuration(duration time.Duration)
}

type Switch interface {
	TurnOn()
	TurnOff()
}

type Devices interface {
	GetAllLights() []string
	GetLight(id string) Light

	GetAllSwitches() []string
	GetSwitch(id string) Switch
}

func New() Devices {
	return NewHueAndWemo()
}
