package config

import (
	"encoding/json"
)

func ServicesConfigFromJSON(data []byte) (*ServicesConfig, error) {
	r := &ServicesConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

func (r *ServicesConfig) FromJSON(data []byte) error {
	c := ServicesConfig{}
	err := json.Unmarshal(data, &c)
	if err == nil {
		*r = c
	}
	return err
}

func (r *ServicesConfig) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type ServicesConfig struct {
	Services []struct {
		ID        string   `json:"id"`
		DependsOn []string `json:"depends_on"`
	} `json:"services"`
}

func (r *ServicesConfig) GetDependencies(serviceID string) []string {
	for _, s := range r.Services {
		if s.ID == serviceID {
			return s.DependsOn
		}
	}
	return nil
}
