package config

import (
	"encoding/json"
)

type XiaomiConfig struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	Motions []struct {
		Name string `json:"name"`
	} `json:"motions"`
	Plugs []struct {
		Name string `json:"name"`
	} `json:"plugs"`
	Switches []struct {
		Name string `json:"name"`
	} `json:"switches"`
	Magnets []struct {
		Name string `json:"name"`
	} `json:"magnets"`
}

func XiaomiConfigFromJSON(jsonstr string) (*XiaomiConfig, error) {
	config := &XiaomiConfig{}
	err := json.Unmarshal([]byte(jsonstr), config)
	return config, err
}

func (m *XiaomiConfig) ToJSON() string {
	data, err := json.Marshal(m)
	if err == nil {
		return string(data)
	}
	return ""
}
