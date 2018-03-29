// To parse and unparse this JSON data, add this code to your project and do:
//
//    config, err := UnmarshalConfig(bytes)
//    bytes, err = config.Marshal()

package weather

import "encoding/json"

func ConfigFromJSON(jsonstr string) (*Config, error) {
	var r *Config
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

func (r *Config) ConfigFromJSON() ([]byte, error) {
	return json.Marshal(r)
}

type Config struct {
	Location    Location `json:"location"`
	Darksky     Aqi      `json:"darksky"`
	Aqi         Aqi      `json:"aqi"`
	IM          IM       `json:"im"`
	Notify      []Notify `json:"notify"`
	Clouds      []Cloud  `json:"clouds"`
	Rain        []Rain   `json:"rain"`
	Wind        []Wind   `json:"wind"`
	Temperature []Cloud  `json:"temperature"`
}

type Aqi struct {
	Key string `json:"key"`
}

type Cloud struct {
	Name        *string    `json:"name"`
	Description string     `json:"description"`
	Unit        *CloudUnit `json:"unit"`
	Min         float64    `json:"min"`
	Max         float64    `json:"max"`
}

type IM struct {
	Channel string `json:"channel"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Notify struct {
	Type    string  `json:"type"`
	Warning Warning `json:"warning"`
	Alert   Alert   `json:"alert"`
}

type Alert struct {
	Min float64 `json:"min"`
	Max int64   `json:"max"`
}

type Warning struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type Rain struct {
	Name         string  `json:"name"`
	Unit         string  `json:"unit"`
	IntensityMin float64 `json:"intensity_min"`
	IntensityMax float64 `json:"intensity_max"`
}

type Wind struct {
	Unit        WindUnit `json:"unit"`
	Speed       float64  `json:"speed"`
	Description []string `json:"description"`
}

type CloudUnit string

const (
	Celcius CloudUnit = "Celcius"
	Empty   CloudUnit = "%"
)

type WindUnit string

const (
	KMH WindUnit = "km/h"
)
