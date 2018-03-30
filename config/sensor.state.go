package config

import (
	"encoding/json"
	"time"
)

// SensorState holds all information of a sensor
// e.g. sensor/weather/aqi
type SensorState struct {
	Domain  string    `json:"domain"`
	Product string    `json:"product"`
	Name    string    `json:"name"`
	Type    string    `json:"type"`
	Value   string    `json:"value"`
	Time    time.Time `json:"time"`
}

func SensorStateFromJSON(jsonstr string) (*SensorState, error) {
	r := &SensorState{}
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

func (r *SensorState) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
