// To parse and unparse this JSON data, add this code to your project and do:
//
//    config, err := UnmarshalConfig(bytes)
//    bytes, err = config.Marshal()

package suncalc

import "encoding/json"

func configFromJSON(jsonstr string) (*Config, error) {
	data := []byte(jsonstr)
	r := &Config{}
	err := json.Unmarshal(data, r)
	return r, err
}

func (r *Config) configToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type Config struct {
	Config   ConfigConfig `json:"config"`
	Anglecfg []Anglecfg   `json:"anglecfg"`
	Moments  []Moment     `json:"moments"`
}

type Anglecfg struct {
	Angle float64 `json:"angle"`
	Rise  string  `json:"rise"`
	Set   string  `json:"set"`
}

type ConfigConfig struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Moment struct {
	Title string `json:"title"`
	Descr string `json:"descr"`
	Begin string `json:"begin"`
	End   string `json:"end  "`
}
