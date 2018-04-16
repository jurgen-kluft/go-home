package config

import "encoding/json"

// BraviaTVConfigFromJSON parser the incoming JSON string and returns an Config instance for Bravia.TV
func BraviaTVConfigFromJSON(jsonstr string) (*BraviaTVConfig, error) {
	r := &BraviaTVConfig{}
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

// ToJSON converts a AqiConfig to a JSON string
func (r *BraviaTVConfig) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type BraviaTVConfig struct {
	Devices []BraviaTV `json:"devices"`
}

type BraviaTV struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	MAC  string `json:"mac"`
}
