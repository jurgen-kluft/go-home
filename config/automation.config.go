package config

import "encoding/json"

// AutomationConfigFromJSON parser the incoming JSON string and returns an Config instance for Aqi
func AutomationConfigFromJSON(data []byte) (*AutomationConfig, error) {
	r := &AutomationConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

// FromJSON converts a json string to a AutomationConfig instance
func (r *AutomationConfig) FromJSON(data []byte) error {
	c := AutomationConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

// ToJSON converts a AutomationConfig to a JSON string
func (r *AutomationConfig) ToJSON() ([]byte, error) {
	data, err := json.Marshal(r)
	if err == nil {
		return data, nil
	}
	return nil, err
}

// AutomationConfig holds the configuration for automation
type AutomationConfig struct {
	SubChannels        []string                 `json:"subscribing_channels"`
	ChannelsToRegister []string                 `json:"register_channels"`
	DeviceControlCache map[string]DeviceControl `json:"device_control_json_cache"`
}

// DeviceControl holds the configuration to control a device
type DeviceControl struct {
	Channel string `json:"channel"`
	On      string `json:"on"`
	Off     string `json:"off"`
	Toggle  string `json:"toggle"`
}
