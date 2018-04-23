package config

import (
	"encoding/json"
)

const (
	BedroomLights      = "Bedroom"
	LivingroomLights   = "Living Room"
	KitchenLights      = "Kitchen"
	SophiaRoomLights   = "Sophia"
	JenniferRoomLights = "Jennifer"
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

func (m *HueConfig) ToJSON() (data []byte, err error) {
	data, err = json.Marshal(m)
	return
}
