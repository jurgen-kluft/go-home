package config

import "encoding/json"

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

// ConbeeConfig contains information to connect to a Conbee-II device
type ConbeeConfig struct {
	Addr   string `json:"url"`
	APIKey string `json:"apikey"`
}
