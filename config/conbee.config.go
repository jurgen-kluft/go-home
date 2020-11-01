package config

import (
	"encoding/json"
	"io/ioutil"
)

// ConbeeConfigFromJSON parser the incoming JSON string and returns an Config instance for Aqi
func ConbeeConfigFromJSON(data []byte) (*ConbeeConfig, error) {
	r := &ConbeeConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

// FromJSON converts a json string to a ConbeeConfig instance
func (r *ConbeeConfig) FromJSON(data []byte) error {
	c := ConbeeConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

// ToJSON converts a ConbeeConfig to a JSON string
func (r *ConbeeConfig) ToJSON() (string, error) {
	data, err := json.Marshal(r)
	if err == nil {
		return string(data), nil
	}
	return "", err
}

// LoadConfig loads a ConbeeConfig from a file
func LoadConfig(filename string) (*ConbeeConfig, error) {
	filedata, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := &ConbeeConfig{}
	err = json.Unmarshal(filedata, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ConbeeConfig contains information to connect to a Conbee-II device
type ConbeeConfig struct {
	Addr        string         `json:"Addr"`
	APIKey      string         `json:"APIKey"`
	LightsOut   string         `json:"lights.out"`
	SwitchesOut string         `json:"switches.out"`
	SensorsOut  string         `json:"sensors.out"`
	LightsIn    []string       `json:"lights.in"`
	Switches    []ConbeeDevice `json:"switches"`
	Sensors     ConbeeSensors  `json:"sensors"`
	Lights      []ConbeeLight  `json:"lights"`
}

// ConbeeLight is a structure that references a light object at Conbee
type ConbeeLight struct {
	Name string   `json:"name"`
	IDS  []string `json:"ids"`
}

// ConbeeSensors is a structure that references sensor objects at Conbee
type ConbeeSensors struct {
	Motion  []ConbeeDevice `json:"motion"`
	Contact []ConbeeDevice `json:"contact"`
}

// ConbeeDevice is a structure that references a device at Conbee
type ConbeeDevice struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BatteryType string `json:"battery_type"`
}
