package config

import "encoding/json"

// AqiConfigFromJSON parser the incoming JSON string and returns an Config instance for Aqi
func AqiConfigFromJSON(jsonstr string) (*AqiConfig, error) {
	r := &AqiConfig{}
	err := json.Unmarshal([]byte(jsonstr), r)
	return r, err
}

// ToJSON converts a AqiConfig to a JSON string
func (r *AqiConfig) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type AqiConfig struct {
	Token    string     `json:"token"`
	City     string     `json:"city"`
	URL      string     `json:"url"`
	Interval int        `json:"interval"`
	Levels   []AqiLevel `json:"levels"`
}

type AqiLevel struct {
	LessThan     float64 `json:"lessthan"`
	Tag          string  `json:"tag"`
	Implications string  `json:"implications"`
	Caution      string  `json:"caution"`
}
