package config

import "encoding/json"

func WeatherConfigFromJSON(jsonstr string) (*WeatherConfig, error) {
	var r *WeatherConfig
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

func (r *WeatherConfig) FromJSON() ([]byte, error) {
	return json.Marshal(r)
}

type WeatherConfig struct {
	Location    Geo      `json:"location"`
	Darksky     Aqi      `json:"darksky"`
	Aqi         Aqi      `json:"aqi"`
	IM          IM       `json:"im"`
	Notify      []Notify `json:"notify"`
	Clouds      []WItem  `json:"clouds"`
	Rain        []WItem  `json:"rain"`
	Wind        []WItem  `json:"wind"`
	Temperature []WItem  `json:"temperature"`
}

type Aqi struct {
	Key string `json:"key"`
}

type WItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        Unit   `json:"unit"`
	Range       MinMax `json:"range"`
}

type IM struct {
	Channel string `json:"channel"`
}

type Notify struct {
	Type    string `json:"type"`
	Warning MinMax `json:"warning"`
	Alert   MinMax `json:"alert"`
}

type Unit string

const (
	Celcius Unit = "Celcius"
	Empty   Unit = "%"
	KMH     Unit = "km/h"
)
