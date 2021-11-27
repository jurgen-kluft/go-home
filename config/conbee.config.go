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
func (r *ConbeeConfig) ToJSON() ([]byte, error) {
	data, err := json.Marshal(r)
	if err == nil {
		return data, nil
	}
	return nil, err
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
	Addr       string             `json:"Addr"`
	APIKey     string             `json:"APIKey"`
	LightsOut  string             `json:"lights.out"`
	SensorsOut string             `json:"sensors.out"`
	LightsIn   []string           `json:"lights.in"`
	Switches   []ConbeeSwitch     `json:"switch"`
	Motion     []ConbeeMotion     `json:"motion"`
	Contact    []ConbeeContact    `json:"contact"`
	Lights     []ConbeeLightGroup `json:"lights"`
}

// ConbeeLight is a structure that references a light object at Conbee
type ConbeeLightGroup struct {
	Name string   `json:"name"`
	IDS  []string `json:"ids"`
}

// ConbeeDevice is a structure that references a general device at Conbee
type ConbeeDevice struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BatteryType string `json:"battery_type"`
}

// ConbeeMotion is a structure that references a motion sensor device at Conbee
type ConbeeMotion struct {
	ConbeeDevice
	On  string `json:"on"`
	Off string `json:"off"`
}

// ConbeeContact is a structure that references a magnet sensor device at Conbee
type ConbeeContact struct {
	ConbeeDevice
	Open  string `json:"open"`
	Close string `json:"close"`
}

// ConbeeContact is a structure that references a switch/button device at Conbee
type ConbeeSwitch struct {
	ConbeeDevice
	SingleClick  string `json:"single_click"`
	DoubleClick  string `json:"double_click"`
	TrippleClick string `json:"tripple_click"`
}
