package config

import (
	"encoding/json"
)

type HueConfig struct {
	Name   string `json:"name"`
	Host   string `json:"host"`
	Key    string `json:"key"`
	Lights []struct {
		Name string `json:"name"`
	} `json:"lights"`
	Groups []struct {
		Name string `json:"name"`
	} `json:"groups"`
}

func HueConfigFromJSON(jsonstr string) (*HueConfig, error) {
	config := &HueConfig{}
	err := json.Unmarshal([]byte(jsonstr), config)
	return config, err
}

func (m *HueConfig) ToJSON() string {
	data, err := json.Marshal(m)
	if err == nil {
		return string(data)
	}
	return ""
}
