package yeel

import (
	"fmt"

	"github.com/thomasf/lg"
)

// Result .
type Result struct {
	ID     int
	Result []string
	Err    error
	Error  map[string]interface{}
}

// Notification .
type Notification struct {
	DeviceID string
	Method   string
	Params   NotificationParams
}

// NotificationParams .
type NotificationParams struct {
	Bright     *int    `json:"bright"`
	CT         *int    `json:"ct"`
	RGB        *int    `json:"rgb"`
	Power      *string `json:"power"`
	Hue        *int    `json:"hue"`
	Sat        *int    `json:"sat"`
	ColorMode  *int    `json:"color_mode"`
	Flowing    *int    `json:"flowing"`
	DelayOff   *int    `json:"delayoff"`
	FlowParams *string `json:"flow_params"`
	MusicOn    *int    `json:"music_on"`
	Name       *string `json:"name"`
}

type resultOrNotification struct {
	*Result
	*Notification
}

// RGB .
type RGB struct {
	R, G, B int
}

func NewRGBNorm(r, g, b float64) RGB {
	return RGB{
		R: scale(0, 255, r),
		G: scale(0, 255, g),
		B: scale(0, 255, b),
	}
}

func (R RGB) Int() int {
	r := R.R
	g := R.G
	b := R.B

	// Adjusting to 1 as minimal value because turning off a whole color looks stranger than the benifits
	if r < 1 {
		r = 1
	}
	if r > 255 {
		r = 255
	}
	if g < 1 {
		g = 1
	}
	if g > 255 {
		g = 255
	}
	if b < 1 {
		b = 1
	}
	if b > 255 {
		b = 255
	}
	rgb := r
	rgb = (rgb << 8) + g
	rgb = (rgb << 8) + b
	return rgb
}

func (R RGB) String() string {
	return fmt.Sprintf("%d,%d,%d", R.R, R.G, R.B)
}

func (R RGB) HexRGB() string {
	return fmt.Sprintf("#%.6x", R.Int())
}

func (rgb RGB) Validate() error {
	if rgb.R < 0 || rgb.R > 255 {
		return newIntError("rgb.r", rgb.R, "0-255")
	}
	if rgb.G < 0 || rgb.G > 255 {
		return newIntError("rgb.g", rgb.G, "0-255")
	}
	if rgb.B < 0 || rgb.B > 255 {
		return newIntError("rgb.B", rgb.B, "0-255")
	}
	return nil
}

func IntToRGB(value int) RGB {
	return RGB{
		R: (value >> 16) & 255,
		G: (value >> 8) & 255,
		B: value & 255,
	}
}

type ColorTemprature int

func (c ColorTemprature) Int() int {
	if err := c.Validate(); err != nil {
		lg.Fatal(err)
	}
	return int(c)
}
func (c ColorTemprature) Validate() error {
	if c < 1700 || c > 6500 {
		return IllegalArgumentError{field: "color_temprature", value: fmt.Sprintf("%d", c), supported: "1700-6500"}
	}
	return nil
}

type Brightness int

func NewBrightnessNorm(v float64) Brightness {
	scaled := scale(1, 100, v)
	return Brightness(scaled)
}

func (b Brightness) Int() int {
	if err := b.Validate(); err != nil {
		lg.Fatal(err)
	}

	return int(b)
}
func (b Brightness) Validate() error {
	if b < 1 || b > 100 {
		return IllegalArgumentError{field: "brightness", value: fmt.Sprintf("%d", b), supported: "1-100"}
	}
	return nil
}

type Hue int

func NewHueNorm(v float64) Hue {
	scaled := scale(1, 359, v)
	return Hue(scaled)
}

func (h Hue) Int() int {
	if err := h.Validate(); err != nil {
		lg.Fatal(err)
	}

	return int(h)
}
func (h Hue) Validate() error {
	if h < 1 || h > 359 {
		return IllegalArgumentError{field: "hue", value: fmt.Sprintf("%d", h), supported: "0-359"}
	}
	return nil
}

type Saturation int

func NewSaturationNorm(v float64) Saturation {
	scaled := scale(0, 100, v)
	return Saturation(scaled)
}
func (s Saturation) Int() int {
	if err := s.Validate(); err != nil {
		lg.Fatal(err)
	}

	return int(s)
}
func (s Saturation) Validate() error {
	if s < 0 || s > 100 {
		return IllegalArgumentError{field: "saturation", value: fmt.Sprintf("%d", s), supported: "0-100"}
	}
	return nil
}
