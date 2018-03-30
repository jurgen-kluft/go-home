package config

import (
	"encoding/json"
	"time"
)

func FluxConfigFromJSON(jsonstr string) (*FluxConfig, error) {
	r := &FluxConfig{}
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

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
	CT          []float64   `json:"ct"`
	Bri         []float64   `json:"bri"`
	Darkorlight Darkorlight `json:"darkorlight"`
	StartMoment string      `json:"startMoment"`
	EndMoment   string      `json:"endMoment"`
	Start       time.Time
	End         time.Time
}

type Lighttype struct {
	Name string `json:"name"`
	CT   MinMax `json:"ct"`
	BRI  MinMax `json:"bri"`
}

type MinMax struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type Season struct {
	Name string `json:"name"`
	CT   MinMax `json:"ct"`
	BRI  MinMax `json:"bri"`
}

type AddSuncalcMoment struct {
	Name  string `json:"name"`
	Tag   string `json:"tag"`
	Shift int64  `json:"shift"`
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
