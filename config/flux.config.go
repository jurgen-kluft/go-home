package config

import (
	"encoding/json"
	"time"
)

// FluxConfigFromJSON converts a json string to a FluxConfig instance
func FluxConfigFromJSON(data []byte) (*FluxConfig, error) {
	r := &FluxConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

// FromJSON converts a json string to a FluxConfig instance
func (r *FluxConfig) FromJSON(data []byte) error {
	c := FluxConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

// ToJSON converts a FluxConfig to a JSON string
func (r *FluxConfig) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type FluxConfig struct {
	Seasons        []Season           `json:"seasons"`
	Weather        []Weather          `json:"weather"`
	SuncalcMoments []AddSuncalcMoment `json:"suncalc_moments"`
	Lighttype      []Lighttype        `json:"lighttype"`
	Lighttime      []Lighttime        `json:"lighttime"`
}

type Lighttime struct {
	CT          FromTo         `json:"ct"`
	Bri         FromTo         `json:"bri"`
	Darkorlight Darkorlight    `json:"darkorlight"`
	TimeSlot    TaggedTimeSlot `json:"timeslot"`
}

type TaggedTimeSlot struct {
	StartMoment string `json:"start"`
	StartTime   time.Time
	EndMoment   string `json:"end"`
	EndTime     time.Time
}

type Lighttype struct {
	MetricsName string `json:"metricsname"`
	LightName   string `json:"lightname"`
	LightType   string `json:"lighttype"`
	Channel     string `json:"channel"`
	CT          MinMax `json:"ct"`
	BRI         MinMax `json:"bri"`
}

type FromTo struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
}

// LinearInterpolated returns interpolated value between From-To
func (f FromTo) LinearInterpolated(x float64) float64 {
	return f.From + x*(f.To-f.From)
}

type MinMax struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// LinearInterpolated returns interpolated value between Min-Max
func (m MinMax) LinearInterpolated(x float64) float64 {
	return m.Min + x*(m.Max-m.Min)
}

type Season struct {
	Name string `json:"name"`
	CT   MinMax `json:"ct"`
	BRI  MinMax `json:"bri"`
}

type AddSuncalcMoment struct {
	Name  string `json:"name"`
	Tag   string `json:"tag"`
	Shift int64  `json:"shift"` // Shift in minutes +/-
}

type Weather struct {
	Clouds MinMax  `json:"clouds"`
	CTPct  float64 `json:"ct_pct"`
	BriPct float64 `json:"bri_pct"`
}

type Darkorlight string

const (
	Dark     Darkorlight = "dark"
	Light    Darkorlight = "light"
	Twilight Darkorlight = "twilight"
)
