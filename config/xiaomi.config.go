package config

import (
	"encoding/json"
)

type XiaomiConfig struct {
	Name    string      `json:"name"`
	IP      string      `json:"ip"`
	MAC     string      `json:"mac"`
	Key     CryptString `json:"key"`
	Motions []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"motions"`
	Plugs []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"plugs"`
	Switches []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"switches"`
	Magnets []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"magnets"`
}

func XiaomiConfigFromJSON(jsonstr string) (*XiaomiConfig, error) {
	config := &XiaomiConfig{}
	err := json.Unmarshal([]byte(jsonstr), config)
	return config, err
}

func (m *XiaomiConfig) ToJSON() (data []byte, err error) {
	data, err = json.Marshal(m)
	return
}
