package config

import (
	"encoding/json"
)

func SuncalcConfigFromJSON(jsonstr string) (*SuncalcConfig, error) {
	data := []byte(jsonstr)
	r := &SuncalcConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

func (r *SuncalcConfig) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type SuncalcConfig struct {
	Geo      Geo       `json:"config"`
	Anglecfg []CAngles `json:"anglecfg"`
	Moments  []CMoment `json:"moments"`
}

type CAngles struct {
	Angle float64 `json:"angle"`
	Rise  string  `json:"rise"`
	Set   string  `json:"set"`
}

type Geo struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CMoment struct {
	Title string `json:"title"`
	Descr string `json:"descr"`
	Begin string `json:"begin"`
	End   string `json:"end"`
}
