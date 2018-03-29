package flux

import (
	"encoding/json"
	"time"

	"github.com/jurgen-kluft/go-home/suncalc"
)

func UnmarshalFluxState(data []byte) (FluxState, error) {
	var r FluxState
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r FluxState) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FluxState struct {
	Lights      []FluxLightType `json:"lights"`
	Darkorlight Darkorlight     `json:"darkorlight"`
}

type FluxLightType struct {
	Name string  `json:"name"`
	CT   float64 `json:"ct"`
	BRI  float64 `json:"bri"`
}

func NewWeatherState(jsonstr string) (*WeatherState, error) {
	data := []byte(jsonstr)
	r := &WeatherState{}
	err := json.Unmarshal(data, r)
	return r, err
}

type WeatherState struct {
	current WeatherForecast   `json:"current"`
	hourly  []WeatherForecast `json:"hourly"`
}

type TimeWindow struct {
	From  time.Time
	Until time.Time
}

func (t TimeWindow) IsInside(start, end time.Time) bool {
	return t.From.After(start) && t.Until.Before(end)
}

type WeatherForecast struct {
	Window      TimeWindow `json:"time_window"`
	Clouds      float64    `json:"clouds"`
	Rain        float64    `json:"rain"`
	Wind        float64    `json:"wind"`
	Temperature float64    `json:"temperature"`
}

func NewSuncalc(jsonstr string) (*suncalc.State, error) {
	r := &suncalc.State{}
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

func UnmarshalConfig(data []byte) (*Config, error) {
	r := &Config{}
	err := json.Unmarshal(data, r)
	return r, err
}

func (r *Config) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Config struct {
	Seasons        []Season        `json:"seasons"`
	Weather        []Weather       `json:"weather"`
	SuncalcMoments []SuncalcMoment `json:"suncalc_moments"`
	Lighttype      []Lighttype     `json:"lighttype"`
	Lighttime      []Lighttime     `json:"lighttime"`
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

type SuncalcMoment struct {
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
