package config

import "encoding/json"

const (
	BedroomSamsungTV = "Bedroom Samsung-TV"
)

// SamsungTVConfigFromJSON parser the incoming JSON string and returns an Config instance for Samsung.TV
func SamsungTVConfigFromJSON(data []byte) (*SamsungTVConfig, error) {
	r := &SamsungTVConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

// FromJSON converts a json string to a SamsungTVConfig instance
func (r *SamsungTVConfig) FromJSON(data []byte) error {
	c := SamsungTVConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
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
