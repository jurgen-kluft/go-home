package aqi

import "encoding/json"

func unmarshalConfig(data []byte) (*Config, error) {
	r := &Config{}
	err := json.Unmarshal(data, r)
	return r, err
}

func (r *Config) marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Config struct {
	Token  string     `json:"token"`
	City   string     `json:"city"`
	URL    string     `json:"url"`
	Levels []AqiLevel `json:"levels"`
}

type AqiLevel struct {
	LessThan     float64 `json:"lessthan"`
	Tag          string  `json:"tag"`
	Implications string  `json:"implications"`
	Caution      string  `json:"caution"`
}
