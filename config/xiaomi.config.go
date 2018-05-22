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

const (
	WirelessSwitchSingleClick = "single click"
	WirelessSwitchDoubleClick = "double click"
	WirelessSwitchLongPress   = "long press"
	WirelessSwitchLongRelease = "long release"
)

type XiaomiConfig struct {
	Name   string      `json:"name"`
	IP     string      `json:"ip"`
	MAC    string      `json:"mac"`
	Key    CryptString `json:"key"`
	Motion []struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		BType string `json:"battery_type"`
	} `json:"motion"`
	Plug []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"plug"`
	Switch []struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		BType string `json:"battery_type"`
	} `json:"switch"`
	Magnet []struct {
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
