package config

import (
	"encoding/json"
	"io/ioutil"
)

// ConbeeConfigFromJSON parser the incoming JSON string and returns an Config instance for Aqi
func ConbeeSensorsConfigFromJSON(data []byte) (*ConbeeSensorsConfig, error) {
	r := &ConbeeSensorsConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

// FromJSON converts a json string to a ConbeeSensorsConfig instance
func (r *ConbeeSensorsConfig) FromJSON(data []byte) error {
	c := ConbeeSensorsConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

// ToJSON converts a ConbeeSensorsConfig to a JSON string
func (r *ConbeeSensorsConfig) ToJSON() ([]byte, error) {
	data, err := json.Marshal(r)
	if err == nil {
		return data, nil
	}
	return nil, err
}

// LoadConbeeSensorsConfig loads a ConbeeSensorsConfig from a file
func LoadConbeeSensorsConfig(filename string) (*ConbeeSensorsConfig, error) {
	filedata, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := &ConbeeSensorsConfig{}
	err = json.Unmarshal(filedata, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ConbeeSensorsConfig contains information to connect to a Conbee-II device
type ConbeeSensorsConfig struct {
	Host        string          `json:"Host"`
	Port        int             `json:"Port"`
	APIKey      string          `json:"APIKey"`
	SwitchesOut string          `json:"switches.out"`
	SensorsOut  string          `json:"sensors.out"`
	Switches    []ConbeeSwitch  `json:"switch"`
	Motion      []ConbeeMotion  `json:"motion"`
	Contact     []ConbeeContact `json:"contact"`
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
