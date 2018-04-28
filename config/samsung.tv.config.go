package config

import "encoding/json"

const (
	BedroomSamsungTV = "Bedroom Samsung-TV"
)

// SamsungTVConfigFromJSON parser the incoming JSON string and returns an Config instance for Samsung.TV
func SamsungTVConfigFromJSON(jsonstr string) (*SamsungTVConfig, error) {
	r := &SamsungTVConfig{}
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

// ToJSON converts a AqiConfig to a JSON string
func (r *SamsungTVConfig) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type SamsungTVConfig struct {
	Devices []SamsungTV `json:"devices"`
}

type SamsungTV struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	ID   string `json:"id"`
}
