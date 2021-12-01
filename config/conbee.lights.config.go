package config

import (
	"encoding/json"
	"io/ioutil"
)

// ConbeeLightsConfigFromJSON parser the incoming JSON string and returns an Config instance for Aqi
func ConbeeLightsConfigFromJSON(data []byte) (*ConbeeLightsConfig, error) {
	r := &ConbeeLightsConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

// FromJSON converts a json string to a ConbeeLightsConfig instance
func (r *ConbeeLightsConfig) FromJSON(data []byte) error {
	c := ConbeeLightsConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

// ToJSON converts a ConbeeLightsConfig to a JSON string
func (r *ConbeeLightsConfig) ToJSON() ([]byte, error) {
	data, err := json.Marshal(r)
	if err == nil {
		return data, nil
	}
	return nil, err
}

// LoadConfig loads a ConbeeLightsConfig from a file
func LoadConfig(filename string) (*ConbeeLightsConfig, error) {
	filedata, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := &ConbeeLightsConfig{}
	err = json.Unmarshal(filedata, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ConbeeLightsConfig contains information to connect to a Conbee-II device
type ConbeeLightsConfig struct {
	Addr     string             `json:"Addr"`
	Port     int                `json:"Port"`
	APIKey   string             `json:"APIKey"`
	LightsIn []string           `json:"lights.in"`
	Lights   []ConbeeLightGroup `json:"lights"`
}

// ConbeeLight is a structure that references a light object at Conbee
type ConbeeLightGroup struct {
	Name  string   `json:"name"`
	Group int      `json:"group"`
	On    string   `json:"on"`
	Off   string   `json:"off"`
	CT    string   `json:"ct"`
	Alert string   `json:"alert"`
	IDS   []string `json:"ids"`
}
