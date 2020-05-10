package config

import "encoding/json"

// AqiConfigFromJSON parser the incoming JSON string and returns an Config instance for Aqi
func AqiConfigFromJSON(data []byte) (*AqiConfig, error) {
	r := &AqiConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

// FromJSON converts json data into AqiConfig
func (r *AqiConfig) FromJSON(data []byte) error {
	c := AqiConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

// ToJSON converts a AqiConfig to a JSON string
func (r *AqiConfig) ToJSON() ([]byte, error) {
	data, err := json.Marshal(r)
	if err == nil {
		return data, nil
	}
	return nil, err
}

type AqiConfig struct {
	Token    CryptString `json:"token"`
	City     string      `json:"city"`
	URL      string      `json:"url"`
	Interval int         `json:"interval"`
	Levels   []AqiLevel  `json:"levels"`
}

type AqiLevel struct {
	LessThan     float64 `json:"lessthan"`
	Tag          string  `json:"tag"`
	Implications string  `json:"implications"`
	Caution      string  `json:"caution"`
}
