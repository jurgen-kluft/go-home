package devices

import (
	"time"
)

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
