package config

import "encoding/json"

// AqiConfigFromJSON parser the incoming JSON string and returns an Config instance for Aqi
func AqiConfigFromJSON(jsonstr string) (*AqiConfig, error) {
	r := &AqiConfig{}
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

func (r *AqiConfig) FromJSON(jsonstr string) (Config, error) {
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

// ToJSON converts a AqiConfig to a JSON string
func (r *AqiConfig) ToJSON() (string, error) {
	data, err := json.Marshal(r)
	if err == nil {
		return string(data), nil
	}
	return "", err
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
