package config

import (
	"encoding/json"
)

type YeeConfig struct {
	Name   string `json:"name"`
	Lights []struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
		Port string `json:"port"`
	} `json:"lights"`
}

func YeeConfigFromJSON(jsonstr string) (*YeeConfig, error) {
	config := &YeeConfig{}
	err := json.Unmarshal([]byte(jsonstr), config)
	return config, err
}

func (m *YeeConfig) ToJSON() string {
	data, err := json.Marshal(m)
	if err == nil {
		return string(data)
	}
	return ""
}
