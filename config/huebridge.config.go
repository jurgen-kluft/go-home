package config

import (
	"encoding/json"
)

// HueBridgeConfig is a struct that holds information for our emulated Hue Bridge
type HueBridgeConfig struct {
	IPPort            string                    `json:"ip_port"`
	RegisterChannels  []string                  `json:"register_channels"`
	SubscribeChannels []string                  `json:"subscribe_channels"`
	EmulatedDevices   []HueBridgeEmulatedDevice `json:"emulated-devices"`
}

// HueBridgeEmulatedDevice is a device emulated on the Hue Bridge
type HueBridgeEmulatedDevice struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
	On      string `json:"on"`
	Off     string `json:"off"`
}

// HueBridgeConfigFromJSON will decode a JSON string to a HueBridgeConfig struct
func HueBridgeConfigFromJSON(jsonstr string) (*HueBridgeConfig, error) {
	config := &HueBridgeConfig{}
	err := json.Unmarshal([]byte(jsonstr), config)
	return config, err
}

// ToJSON will encode a HueBridgeConfig struct to a JSON string
func (m *HueBridgeConfig) ToJSON() (data []byte, err error) {
	data, err = json.Marshal(m)
	return
}
