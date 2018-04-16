package config

import (
	"encoding/json"
)

type ShoutConfig struct {
	Key CryptString `json:"key"`
}

func ShoutConfigFromJSON(jsonstr string) (*ShoutConfig, error) {
	config := &ShoutConfig{}
	err := json.Unmarshal([]byte(jsonstr), config)
	return config, err
}

func (m *ShoutConfig) ToJSON() (data []byte, err error) {
	data, err = json.Marshal(m)
	return
}

type ShoutMsg struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Pretext  string `json:"pretext"`
	Prebody  string `json:"prebody"`
}

func ShoutMsgFromJSON(jsonstr string) (ShoutMsg, error) {
	var msg ShoutMsg
	err := json.Unmarshal([]byte(jsonstr), &msg)
	return msg, err
}

func (m *ShoutMsg) ToJSON() string {
	data, err := json.Marshal(m)
	if err == nil {
		return string(data)
	}
	return ""
}
