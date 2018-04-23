package config

import (
	"encoding/json"
)

const (
	KitchenMotionSensor          = "Kitchen Motion"
	LivingroomMotionSensor       = "Livingroom Motion"
	BedroomMotionSensor          = "Bedroom Motion"
	BedroomPowerPlug             = "Bedroom Plug"
	BedroomCeilingLightSwitch    = "Bedroom Ceiling Light-Switch"
	BedroomChandelierLightSwitch = "Bedroom Chandelier Light-Switch"
	BedroomSwitch                = "Bedroom Switch"
	SophiaRoomSwitch             = "Sophia Switch"
	FrontdoorMagnetSensor        = "Front Door Magnet"
)

type XiaomiConfig struct {
	Name    string      `json:"name"`
	IP      string      `json:"ip"`
	MAC     string      `json:"mac"`
	Key     CryptString `json:"key"`
	Motions []struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		BType string `json:"battery_type"`
	} `json:"motion"`
	Plugs []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"plug"`
	Switches []struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		BType string `json:"battery_type"`
	} `json:"switch"`
	Magnets []struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		BType string `json:"battery_type"`
	} `json:"magnet"`
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
