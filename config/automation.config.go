package config

import "encoding/json"

const (
	BedroomLightStand         = "Bedroom Light Stand"
	BedroomLightMain          = "Bedroom Light Main"
	LivingroomLightStand      = "Living Room Stand"
	LivingroomLightMain       = "Living Room Main"
	LivingroomLightChandelier = "Living Room Chandelier"
	KitchenLights             = "Kitchen"
	SophiaRoomLightStand      = "Sophia Stand"
	SophiaRoomLightMain       = "Sophia Main"
	JenniferRoomLightMain     = "Jennifer"
)

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
	Subscribe   []string                 `json:"subscribe"`
	Register    []string                 `json:"register"`
	DeviceCache map[string]DeviceControl `json:"devices"`
}

// DeviceControl holds the configuration to control a device
type DeviceControl struct {
	Channel string `json:"channel"`
	On      string `json:"on"`
	Off     string `json:"off"`
	Toggle  string `json:"toggle"`
}
