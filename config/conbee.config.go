package config

import "encoding/json"

// ConbeeConfigFromJSON parser the incoming JSON string and returns an Config instance for Aqi
func ConbeeConfigFromJSON(jsonstr string) (*ConbeeConfig, error) {
	r := &ConbeeConfig{}
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

// ToJSON converts a ConbeeConfig to a JSON string
func (r *ConbeeConfig) ToJSON() (string, error) {
	data, err := json.Marshal(r)
	if err == nil {
		return string(data), nil
	}
	return "", err
}

// ConbeeConfig contains information to connect to a Conbee-II device
type ConbeeConfig struct {
	Addr   string `json:"url"`
	APIKey string `json:"apikey"`
}
