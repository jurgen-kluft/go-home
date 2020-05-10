package config

import (
	"encoding/json"
)

func PresenceConfigFromJSON(data []byte) (*PresenceConfig, error) {
	r := &PresenceConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

func (r *PresenceConfig) FromJSON(data []byte) error {
	c := PresenceConfig{}
	err := json.Unmarshal(data, &c)
	if err == nil {
		*r = c
	}
	return err
}

func (r *PresenceConfig) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type PresenceConfig struct {
	Name              string `json:"name"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	User              string `json:"user"`
	Password          string `json:"password"`
	UpdateHistory     int    `json:"update_history"`
	UpdateIntervalSec int    `json:"update_interval_sec"`
	Devices           []struct {
		Name string `json:"name"`
		Mac  string `json:"mac"`
	} `json:"devices"`
}
