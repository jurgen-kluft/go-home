package config

import (
	"encoding/json"
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

func XiaomiConfigFromJSON(data []byte) (*XiaomiConfig, error) {
	config := &XiaomiConfig{}
	err := json.Unmarshal(data, config)
	return config, err
}

// FromJSON converts a json string to a XiaomiConfig instance
func (r *XiaomiConfig) FromJSON(data []byte) error {
	c := XiaomiConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

func (m *XiaomiConfig) ToJSON() (data []byte, err error) {
	data, err = json.Marshal(m)
	return
}
