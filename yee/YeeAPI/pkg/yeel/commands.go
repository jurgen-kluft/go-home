package yeel

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Commander interface {
	Command() (Command, error)
}

// Command .
type Command struct {
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	resultC chan Result
	errorC  chan error
}

func newCommand(method string, params ...interface{}) (Command, error) {
	if params == nil {
		params = make([]interface{}, 0)
	}
	return Command{
		ID:      rand.Int() % 99999,
		Method:  method,
		Params:  params,
		resultC: make(chan Result, 0),
		errorC:  make(chan error, 0),
	}, nil

}

func newDurationCommand(method string, dur time.Duration, params ...interface{}) (Command, error) {
	if err := validateDurationShort("duration", dur); err != nil {
		return Command{}, err
	}

	var effect string
	if dur > 0 {
		effect = "smooth"
	} else {
		effect = "sudden"
	}
	durMillis := int(dur / time.Millisecond)
	var allparams []interface{}
	allparams = append(allparams, params...)
	allparams = append(allparams, effect, durMillis)
	return newCommand(method, allparams...)

}

// SetBrightnessCommand .
type SetBrightnessCommand struct {
	Brightness int           `json:"brightness"`
	Duration   time.Duration `json:"duration"`
}

func (s SetBrightnessCommand) Command() (Command, error) {
	b := Brightness(s.Brightness)

	if err := b.Validate(); err != nil {
		return Command{}, err
	}
	return newDurationCommand("set_bright", s.Duration, b.Int())

}

// SetBrightnessNormCommand .
type SetBrightnessNormCommand struct {
	Brightness float64
	Duration   time.Duration
}

func (s SetBrightnessNormCommand) Command() (Command, error) {
	b := NewBrightnessNorm(s.Brightness)
	if err := b.Validate(); err != nil {
		return Command{}, err
	}
	return newDurationCommand("set_bright", s.Duration, b.Int())

}

type SetRGBCommand struct {
	R, G, B  int
	Duration time.Duration
}

func (s SetRGBCommand) Command() (Command, error) {
	rgb := RGB{s.R, s.G, s.B}
	if err := rgb.Validate(); err != nil {
		return Command{}, err
	}

	return newDurationCommand("set_rgb", s.Duration, rgb.Int())

}

type SetRGBNormCommand struct {
	R, G, B  float64
	Duration time.Duration
}

func (s SetRGBNormCommand) Command() (Command, error) {
	rgb := NewRGBNorm(s.R, s.G, s.B)
	if err := rgb.Validate(); err != nil {
		return Command{}, err
	}

	return newDurationCommand("set_rgb", s.Duration, rgb.Int())

}

// SetNameCommand .
type SetNameCommand struct {
	Name string
}

func (c SetNameCommand) Command() (Command, error) {
	return newCommand("set_name", c.Name)
}

// ToggleCommand .
type ToggleCommand struct {
}

func (c ToggleCommand) Command() (Command, error) {
	return newCommand("toggle")
}

// SetDefaultCommand .
type SetDefaultCommand struct {
}

func (c SetDefaultCommand) Command() (Command, error) {
	return newCommand("set_default")
}

// SetColorTempratureCommand .
type SetColorTempratureCommand struct {
	ColorTemprature int
	Duration        time.Duration
}

func (c SetColorTempratureCommand) Command() (Command, error) {
	ct := ColorTemprature(c.ColorTemprature)
	if err := ct.Validate(); err != nil {
		return Command{}, err
	}
	return newDurationCommand("set_ct_abx", c.Duration, ct.Int())
}

// SetColorTempratureNormCommand .
type SetColorTempratureNormCommand struct {
	ColorTemprature float64
	Duration        time.Duration
}

// SetHSVCommand .
type SetHSVCommand struct {
	Hue        int
	Saturation int
	Duration   time.Duration
}

func (s SetHSVCommand) Command() (Command, error) {
	hue := Hue(s.Hue)
	if err := hue.Validate(); err != nil {
		return Command{}, err
	}
	sat := Saturation(s.Saturation)
	if err := sat.Validate(); err != nil {
		return Command{}, err
	}
	return newDurationCommand("set_hsv", s.Duration, hue.Int(), sat.Int())
}

// SetHSVCommand .
type SetHSVNormCommand struct {
	Hue        float64
	Saturation float64
	Duration   time.Duration
}

func (s SetHSVNormCommand) Command() (Command, error) {
	hue := NewHueNorm(s.Hue)
	if err := hue.Validate(); err != nil {
		return Command{}, err
	}
	sat := NewSaturationNorm(s.Saturation)
	if err := sat.Validate(); err != nil {
		return Command{}, err
	}
	return newDurationCommand("set_hsv", s.Duration, hue.Int(), sat.Int())
}

// SetPowerCommand .
type SetPowerCommand struct {
	On       bool
	Duration time.Duration
}

func (s SetPowerCommand) Command() (Command, error) {
	var strvalue string
	if s.On {
		strvalue = "on"
	} else {
		strvalue = "off"
	}
	return newDurationCommand("set_power", s.Duration, strvalue)

}

// StartColorFlowCommand .
type StartColorFlowCommand struct {
	Count     int
	Action    int
	Animation Animation
}

func (s StartColorFlowCommand) Command() (Command, error) {
	expstr, err := s.Animation.Expression()
	if err != nil {
		return Command{}, err
	}

	count := s.Count
	if count < 0 {
		return Command{}, IllegalArgumentError{field: "count", value: fmt.Sprintf("%d", count), supported: "0>"}
	}
	if count > 0 {
		count = len(s.Animation) * count
	}
	return newCommand("start_cf", count, s.Action, expstr)
}

// StopColorFlowCommand .
type StopColorFlowCommand struct{}

func (c StopColorFlowCommand) Command() (Command, error) {
	return newCommand("stop_cf")
}

// GetPropCommand .
type GetPropCommand struct {
	Properties []string
}

func (g GetPropCommand) Command() (Command, error) {
	var errors []string
	for _, v := range g.Properties {
		if _, ok := validProperties[v]; !ok {
			errors = append(errors, v)
		}
	}
	if len(errors) > 0 {
		return Command{}, IllegalArgumentError{"properties", strings.Join(errors, ","), ""}
	}
	var args []interface{}
	for _, v := range g.Properties {
		args = append(args, v)
	}
	return newCommand("get_prop", args...)
}
