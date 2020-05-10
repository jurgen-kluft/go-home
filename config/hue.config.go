package config

import (
	"encoding/json"
)

// Group defines
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

// HueConfigFromJSON converts a json string to a HueConfig instance
func HueConfigFromJSON(data []byte) (*HueConfig, error) {
	config := &HueConfig{}
	err := json.Unmarshal(data, config)
	return config, err
}

// FromJSON converts a json string to a HueConfig instance
func (r *HueConfig) FromJSON(data []byte) error {
	c := HueConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

// ToJSON converts a HueConfig to a JSON string
func (r *HueConfig) ToJSON() (data []byte, err error) {
	data, err = json.Marshal(r)
	return
}
