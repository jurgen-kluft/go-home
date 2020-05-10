package config

import "encoding/json"

const (
	LivingroomBraviaTV = "Livingroom Sony Bravia-TV"
)

// BraviaTVConfigFromJSON parser the incoming JSON string and returns an Config instance for Bravia.TV
func BraviaTVConfigFromJSON(data []byte) (*BraviaTVConfig, error) {
	r := &BraviaTVConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

// FromJSON converts a json string to a BraviaTVConfig instance
func (r *BraviaTVConfig) FromJSON(data []byte) error {
	c := BraviaTVConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
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
