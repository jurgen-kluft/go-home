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

func WemoConfigFromJSON(data []byte) (*WemoConfig, error) {
	config := &WemoConfig{}
	err := json.Unmarshal(data, config)
	return config, err
}

// FromJSON converts a json string to a WemoConfig instance
func (r *WemoConfig) FromJSON(data []byte) error {
	c := WemoConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

func (m *WemoConfig) ToJSON() (data []byte, err error) {
	data, err = json.Marshal(m)
	return
}
