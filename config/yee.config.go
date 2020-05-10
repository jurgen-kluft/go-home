package config

import (
	"encoding/json"
)

const (
	FrontdoorHallLight = "Frontdoor hall light"
)

type YeeConfig struct {
	Name   string `json:"name"`
	Lights []struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
		Port string `json:"port"`
	} `json:"lights"`
}

func YeeConfigFromJSON(data []byte) (*YeeConfig, error) {
	config := &YeeConfig{}
	err := json.Unmarshal(data, config)
	return config, err
}

// FromJSON converts a json string to a YeeConfig instance
func (r *YeeConfig) FromJSON(data []byte) error {
	c := YeeConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

func (m *YeeConfig) ToJSON() (data []byte, err error) {
	data, err = json.Marshal(m)
	return
}
