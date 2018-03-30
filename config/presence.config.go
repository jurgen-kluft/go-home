package config

import (
	"encoding/json"
)

func PresenceConfigFromJSON(jsonstr string) (*PresenceConfig, error) {
	r := &PresenceConfig{}
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

func (r *PresenceConfig) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type PresenceConfig struct {
	Name       string  `json:"name"`
	Host       string  `json:"host"`
	Port       int     `json:"port"`
	User       string  `json:"user"`
	Password   string  `json:"password"`
	UpdateHist int     `json:"uhist"`
	UpdateFreq float64 `json:"ufreq"`
	Devices    []struct {
		Name string `json:"name"`
		Mac  string `json:"mac"`
	} `json:"devices"`
}
