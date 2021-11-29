package deconz

import (
	"context"
	"encoding/json"
)

// GetSensors retrieves all the sensors available on the gatway
func (c *Client) GetSensors(ctx context.Context) (GetSensorsResponse, error) {
	sensorsResp := GetSensorsResponse{}

	err := c.get(ctx, "sensors", &sensorsResp)
	if err != nil {
		return nil, err
	}

	return sensorsResp, nil
}

// GetSensor retrieves the specified sensor
func (c *Client) GetSensor(ctx context.Context, id string) (*Sensor, error) {
	sensor := &Sensor{}

	err := c.get(ctx, "sensors/"+id, sensor)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}

// SetSensorState specifies the new state of a sensor
func (c *Client) SetSensorState(ctx context.Context, id string, newState *SetSensorStateRequest) error {
	return c.put(ctx, "sensors/"+id+"/state", newState)
}

// SetSensor specifies the new options for a sensor
func (c *Client) SetSensor(ctx context.Context, id string, newConfig *SetSensorRequest) error {
	return c.put(ctx, "sensors/"+id, newConfig)
}

// SetSensorConfig specifies the new config of a sensor
func (c *Client) SetSensorConfig(ctx context.Context, id string, newConfig *SetSensorConfigRequest) error {
	return c.put(ctx, "sensors/"+id+"/config", newConfig)
}

// DeleteSensor removes the specified sensor from the gateway
func (c *Client) DeleteSensor(ctx context.Context, id string) error {
	return c.delete(ctx, "sensors/"+id)
}

// Sensor represents a generic sensor in a Zigbee network
type Sensor struct {
	SensorMetadata

	AlarmState          *ZHAAlarm
	CarbonMonoxideState *ZHACarbonMonoxide
	ConsumptionState    *ZHAConsumption
	FireState           *ZHAFire
	HumidityState       *ZHAHumidity
	LightLevelState     *ZHALightLevel
	OpenCloseState      *ZHAOpenClose
	PowerState          *ZHAPower
	PresenceState       *ZHAPresence
	SwitchState         *ZHASwitch
	PressureState       *ZHAPressure
	TemperatureState    *ZHATemperature
	ThermostatState     *ZHAThermostat
	VibrationState      *ZHAVibration
	WaterState          *ZHAWater
	ButtonState         *ZGPSwitch
}

// SensorMetadata contains a bunch of fields about all sensors
type SensorMetadata struct {
	ID               int          `json:"ep"`
	Config           SensorConfig `json:"config"`
	ETag             string       `json:"etag"`
	ManufacturerName string       `json:"manufacturername"`
	ModelID          string       `json:"modelid"`
	Mode             int          `json:"mode"`
	Name             string       `json:"name"`
	SoftwareVersion  string       `json:"swversion"`
	Type             string       `json:"type"`
	UniqueID         string       `json:"uniqueid"`

	StateRaw json.RawMessage `json:"state"`
}

// UnmarshalJSON is a custom unmarshaler for the State object
func (s *Sensor) UnmarshalJSON(b []byte) error {
	meta := &SensorMetadata{}
	err := json.Unmarshal(b, meta)
	if err != nil {
		return err
	}

	s.SensorMetadata = *meta

	switch meta.Type {
	case "ZHAAlarm":
		s.AlarmState = &ZHAAlarm{}
		err = json.Unmarshal(meta.StateRaw, s.AlarmState)
	case "ZHACarbonMonoxide":
		s.CarbonMonoxideState = &ZHACarbonMonoxide{}
		err = json.Unmarshal(meta.StateRaw, s.CarbonMonoxideState)
	case "ZHAConsumption":
		s.ConsumptionState = &ZHAConsumption{}
		err = json.Unmarshal(meta.StateRaw, s.ConsumptionState)
	case "ZHAFire":
		s.FireState = &ZHAFire{}
		err = json.Unmarshal(meta.StateRaw, s.FireState)
	case "ZHAHumidity":
		s.HumidityState = &ZHAHumidity{}
		err = json.Unmarshal(meta.StateRaw, s.HumidityState)
	case "ZHALightLevel":
		s.LightLevelState = &ZHALightLevel{}
		err = json.Unmarshal(meta.StateRaw, s.LightLevelState)
	case "ZHAOpenClose":
		s.OpenCloseState = &ZHAOpenClose{}
		err = json.Unmarshal(meta.StateRaw, s.OpenCloseState)
	case "ZHAPower":
		s.PowerState = &ZHAPower{}
		err = json.Unmarshal(meta.StateRaw, s.PowerState)
	case "ZHAPresence":
		s.PresenceState = &ZHAPresence{}
		err = json.Unmarshal(meta.StateRaw, s.PresenceState)
	case "ZHASwitch":
		s.SwitchState = &ZHASwitch{}
		err = json.Unmarshal(meta.StateRaw, s.SwitchState)
	case "ZHAPressure":
		s.PressureState = &ZHAPressure{}
		err = json.Unmarshal(meta.StateRaw, s.PressureState)
	case "ZHATemperature":
		s.TemperatureState = &ZHATemperature{}
		err = json.Unmarshal(meta.StateRaw, s.TemperatureState)
	case "ZHAThermostat":
		s.ThermostatState = &ZHAThermostat{}
		err = json.Unmarshal(meta.StateRaw, s.ThermostatState)
	case "ZHAVibration":
		s.VibrationState = &ZHAVibration{}
		err = json.Unmarshal(meta.StateRaw, s.VibrationState)
	case "ZHAWater":
		s.WaterState = &ZHAWater{}
		err = json.Unmarshal(meta.StateRaw, s.WaterState)
	case "ZGPSwitch":
		s.ButtonState = &ZGPSwitch{}
		err = json.Unmarshal(meta.StateRaw, s.ButtonState)
	}

	return err
}

// SensorConfig contains the settable properties of a sensor
type SensorConfig struct {
	On           bool `json:"on"`
	Reachable    bool `json:"reachable"`
	BatteryLevel int  `json:"battery"`
}

// SensorState contains the reported, immutable properties of a sensor.
// This is a generic type which contains state for all possible Zigbee sensors.
// Specific sensor types are subclassed and exposed with only their relevant fields.
type SensorState struct {
	LastUpdated string `json:"lastupdated"`
	LowBattery  bool   `json:"lowbattery"`
	Tampered    bool   `json:"tampered"`

	Alarm          bool   `json:"alarm"`
	CarbonMonoxide bool   `json:"carbonmonoxide"`
	Consumption    int    `json:"consumption"`
	Power          int    `json:"power"`
	Fire           bool   `json:"fire"`
	Humidity       int    `json:"humidity"`
	Lux            int    `json:"lux"`
	LightLevel     int    `json:"lightlevel"`
	Dark           bool   `json:"dark"`
	Daylight       bool   `json:"daylight"`
	Open           bool   `json:"open"`
	Current        int    `json:"current"`
	Voltage        int    `json:"voltage"`
	Presence       bool   `json:"presence"`
	ButtonEvent    int    `json:"buttonevent"`
	Gesture        int    `json:"gesture"`
	EventDuration  int    `json:"eventduration"`
	X              int    `json:"x"`
	Y              int    `json:"y"`
	Angle          int    `json:"angle"`
	Pressure       int    `json:"pressure"`
	Temperature    int    `json:"temperature"`
	Valve          int    `json:"valve"`
	WindowOpen     string `json:"windowopen"`
}

// ZHAAlarm represents a Zigbee Home Automation Alarm
type ZHAAlarm struct {
	Alarm       bool   `json:"alarm"`
	LastUpdated string `json:"lastupdated"`
	LowBattery  bool   `json:"lowbattery"`
	Tampered    bool   `json:"tampered"`
}

// ZHACarbonMonoxide represents a Zigbee Home Automation Carbon Monoxide detector
type ZHACarbonMonoxide struct {
	CarbonMonoxide bool   `json:"carbonmonoxide"`
	LastUpdated    string `json:"lastupdated"`
	LowBattery     bool   `json:"lowbattery"`
	Tampered       bool   `json:"tampered"`
}

// ZHAConsumption represents a Zigbee Home Automation consumption monitor
type ZHAConsumption struct {
	Consumption int    `json:"consumption"`
	LastUpdated string `json:"lastupdated"`
	Power       int    `json:"power"`
}

// ZHAFire represents a Zigbee Home Automation fire detector
type ZHAFire struct {
	Fire        bool   `json:"fire"`
	LastUpdated string `json:"lastupdated"`
	LowBattery  bool   `json:"lowbattery"`
	Tampered    bool   `json:"tampered"`
}

// ZHAHumidity represents a Zigbee Home Automation humidity monitor
type ZHAHumidity struct {
	Humidity    int    `json:"humidity"`
	LastUpdated string `json:"lastupdated"`
}

// ZHALightLevel represents a Zigbee Home Automation light level sensor
type ZHALightLevel struct {
	Lux         int    `json:"lux"`
	LastUpdated string `json:"lastupdated"`
	LightLevel  int    `json:"lightlevel"`
	Dark        bool   `json:"dark"`
	Daylight    bool   `json:"daylight"`
}

// ZHAOpenClose represents an open/close sensor
type ZHAOpenClose struct {
	Open        bool   `json:"open"`
	LastUpdated string `json:"lastupdated"`
}

// ZHAPower represents a Zigbee power monitor
type ZHAPower struct {
	Current     int    `json:"current"`
	LastUpdated string `json:"lastupdated"`
	Power       int    `json:"power"`
	Voltage     int    `json:"voltage"`
}

// ZHAPresence represents a Zigbee presence monitor
type ZHAPresence struct {
	Presence    bool   `json:"presence"`
	LastUpdated string `json:"lastupdated"`
}

// ZHASwitch represents a Zigbee switch
type ZHASwitch struct {
	ButtonEvent   int    `json:"buttonevent"`
	LastUpdated   string `json:"lastupdated"`
	Gesture       int    `json:"gesture"`
	EventDuration int    `json:"eventduration"`
	X             int    `json:"x"`
	Y             int    `json:"y"`
	Angle         int    `json:"angle"`
}

// ZHAPressure represents a Zigbee pressure monitor
type ZHAPressure struct {
	Pressure    int    `json:"pressure"`
	LastUpdated string `json:"lastupdated"`
}

// ZHATemperature represents a Zigbee temperature sensor
type ZHATemperature struct {
	Temperature int    `json:"temperature"`
	LastUpdated string `json:"lastupdated"`
}

// ZHAThermostat represents a Zigbee thermostat
type ZHAThermostat struct {
	On          bool   `json:"on"`
	LastUpdated string `json:"lastupdated"`
	Temperature int    `json:"temperature"`
	Valve       int    `json:"valve"`
	WindowOpen  string `json:"windowopen"`
}

// ZHAVibration represents a Zigbee vibration sensor
type ZHAVibration struct {
	Vibration         bool   `json:"vibration"`
	LastUpdated       string `json:"lastupdated"`
	OrientationX      int    `json:"orientation_x"`
	OrientationY      int    `json:"orientation_y"`
	OrientationZ      int    `json:"orientation_z"`
	TiltAngle         int    `json:"tiltangle"`
	VibrationStrength int    `json:"vibrationstrength"`
}

// ZHAWater represents a Zigbee water sensor
type ZHAWater struct {
	Water       bool   `json:"water"`
	LastUpdated string `json:"lastupdated"`
	LowBattery  bool   `json:"lowbattery"`
	Tampered    bool   `json:"tampered"`
}

// ZGPSwitch represents a Zigbee general button event
type ZGPSwitch struct {
	ButtonEvent int    `json:"buttonevent"`
	LastUpdated string `json:"lastupdated"`
}

// GetSensorsResponse contains the set of sensors in the gateway
type GetSensorsResponse map[string]Sensor

// SetSensorRequest allows for specific parts of a sensor to be changed
type SetSensorRequest struct {
	Name string `json:"name,omitempty"`
	// Mode is only available for dresden elektronik Lighting Switches
	// 1 represents Scene mode
	// 2 represents two-groups mode
	// 3 represents colour temperature mode
	Mode int `json:"mode,omitempty"`
}

// SetSensorConfigRequest contains the fields of a sensor which can be changed.
type SetSensorConfigRequest SensorConfig

// SetSensorStateRequest contains the relevant properties which can be set.
type SetSensorStateRequest struct {
	// ButtonEvent is settable for CLIPSwitch type
	ButtonEvent int `json:"buttonevent,omitempty"`

	// Open is settable for CLIPOpenClose
	Open bool `json:"open,omitempty"`

	// Presence is settable for CLIPPresence
	Presence bool `json:"presence,omitempty"`

	// Temperature is settable for CLIPTemperature
	Temperature int `json:"temperature,omitempty"`

	// Flag is settable for CLIPGenericFlag
	Flag bool `json:"flag,omitempty"`

	// Status is settable for CLIPGenericStatus
	Status int `json:"status,omitempty"`

	// Humidity is settable for CLIPHumidity
	Humidity int `json:"humidity,omitempty"`
}
