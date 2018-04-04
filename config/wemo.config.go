package config

import (
	"encoding/json"
)

type WemoConfig struct {
	Name    string `json:"name"`
	Devices []struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
		Port string `json:"port"`
	} `json:"devices"`
}

func WemoConfigFromJSON(jsonstr string) (*WemoConfig, error) {
	config := &WemoConfig{}
	err := json.Unmarshal([]byte(jsonstr), config)
	return config, err
}

func (m *WemoConfig) ToJSON() string {
	data, err := json.Marshal(m)
	if err == nil {
		return string(data)
	}
	return ""
}
