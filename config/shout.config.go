package config

import (
	"encoding/json"
)

type ShoutConfig struct {
	UserToken CryptString `json:"usertoken"`
	AppToken  CryptString `json:"apptoken"`
	Channel   string      `json:"channel"`
}

// ShoutConfigFromJSON converts a json string to a ShoutConfig instance
func ShoutConfigFromJSON(data []byte) (*ShoutConfig, error) {
	config := &ShoutConfig{}
	err := json.Unmarshal(data, config)
	return config, err
}

// FromJSON converts a json string to a ShoutConfig instance
func (m *ShoutConfig) FromJSON(data []byte) error {
	c := ShoutConfig{}
	err := json.Unmarshal(data, &c)
	*m = c
	return err
}

// ToJSON converts a ShoutConfig to a JSON string
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

func ShoutMsgFromJSON(jsondata []byte) (ShoutMsg, error) {
	var msg ShoutMsg
	err := json.Unmarshal(jsondata, &msg)
	return msg, err
}

func (m *ShoutMsg) FromJSON(data []byte) error {
	c := ShoutMsg{}
	err := json.Unmarshal(data, &c)
	if err == nil {
		*m = c
	}
	return err
}

func (m *ShoutMsg) ToJSON() []byte {
	data, err := json.Marshal(m)
	if err == nil {
		return data
	}
	return nil
}
