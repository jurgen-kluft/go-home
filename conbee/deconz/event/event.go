package event

import (
	"encoding/json"
	"fmt"
)

// TypeLookuper is the interface that we require to lookup types from id's
type TypeLookuper interface {
	SupportsResource(string) bool
	LookupType(string) (string, error)
}

// Event represents a deconz sensor event
type Event struct {
	Type     string          `json:"t"`
	Event    string          `json:"e"`
	Resource string          `json:"r"`
	UniqueID string          `json:"uniqueid"`
	ID       int             `json:"id,string"`
	RawState json.RawMessage `json:"state"`
	State    interface{}
}

// Decoder is able to decode deCONZ events
type Decoder struct {
	TypeStore TypeLookuper
}

// Parse parses events from bytes
func (d *Decoder) Parse(b []byte) (*Event, error) {
	var e Event
	err := json.Unmarshal(b, &e)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal json: %s", err)
	}

	// If there is no state, dont try to parse it
	// TODO: figure out what to do with these
	//       some of them seems to be battery updates
	if !d.TypeStore.SupportsResource(e.Resource) || len(e.RawState) == 0 {
		if !d.TypeStore.SupportsResource(e.Resource) {
			fmt.Printf("Unsupported device type '%s': %s\n", e.Resource, string(e.RawState))
		}
		e.State = &EmptyState{}
		return &e, nil
	}

	err = e.ParseState(d.TypeStore)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal state: %s", err)
	}

	return &e, nil
}

// ParseState tries to unmarshal the appropriate state based
// on looking up the id though the TypeStore
func (e *Event) ParseState(tl TypeLookuper) error {

	t, err := tl.LookupType(e.UniqueID)
	if err != nil {
		return fmt.Errorf("unable to lookup event id %s: %s", e.UniqueID, err)
	}

	switch t {
	case "ZHAFire":
		var s ZHAFire
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHATemperature":
		var s ZHATemperature
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAPressure":
		var s ZHAPressure
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAHumidity":
		var s ZHAHumidity
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAWater":
		var s ZHAWater
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHASwitch":
		var s ZHASwitch
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAPresence":
		var s ZHAPresence
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAOpenClose":
		var s ZHAOpenClose
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "CLIPPresence":
		var s ZHAPresence
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "CLIPGenericStatus":
		err = nil
		e.State = string(e.RawState)
		break
	case "Daylight":
		var s Daylight
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
	case "Extended color light":
		var s LightState
		fmt.Printf("light-state: %s\n", string(e.RawState))
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
	case "Dimmable light":
		var s LightState
		fmt.Printf("light-state: %s\n", string(e.RawState))
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
	default:
		e.Resource = "Unknown"
		fmt.Printf("event state: %s is not a known type", t)
		//err = fmt.Errorf("unable to unmarshal event state: %s is not a known type", t)
	}

	// err should continue to be null if everythings ok
	return err
}

// State is for embedding into event states
type State struct {
	Lastupdated string
}

// LightState represent the state of a extended color light type
type LightState struct {
	State
	ColorMode string
	Bri       int
	CT        int
	On        bool
	Reachable bool
}

// IsLightState returns true if s is of type LightState
func IsLightState(s interface{}) bool {
	switch s.(type) {
	case *LightState:
		return true
	}
	return false
}

// ZHAHumidity represents a presure change
type ZHAHumidity struct {
	State
	Humidity int
}

// Fields returns timeseries data for influxdb
func (z *ZHAHumidity) Fields() map[string]interface{} {
	return map[string]interface{}{
		"humidity": float64(z.Humidity) / 100,
	}
}

// ZHAPressure represents a presure change
type ZHAPressure struct {
	State
	Pressure int
}

// Fields returns timeseries data for influxdb
func (z *ZHAPressure) Fields() map[string]interface{} {
	return map[string]interface{}{
		"pressure": z.Pressure,
	}
}

// ZHATemperature represents a temperature change
type ZHATemperature struct {
	State
	Temperature int
}

// Fields returns timeseries data for influxdb
func (z *ZHATemperature) Fields() map[string]interface{} {
	return map[string]interface{}{
		"temperature": float64(z.Temperature) / 100,
	}
}

// ZHAWater respresents a change from a flood sensor
type ZHAWater struct {
	State
	Lowbattery bool
	Tampered   bool
	Water      bool
}

// Fields returns timeseries data for influxdb
func (z *ZHAWater) Fields() map[string]interface{} {
	return map[string]interface{}{
		"lowbattery": z.Lowbattery,
		"tampered":   z.Tampered,
		"water":      z.Water,
	}
}

// ZHAFire represents a change from a smoke detector
type ZHAFire struct {
	State
	Fire       bool
	Lowbattery bool
	Tampered   bool
}

// Fields returns timeseries data for influxdb
func (z *ZHAFire) Fields() map[string]interface{} {
	return map[string]interface{}{
		"lowbattery": z.Lowbattery,
		"tampered":   z.Tampered,
		"fire":       z.Fire,
	}
}

// ZHASwitch represents a change from a button or switch
type ZHASwitch struct {
	State
	Buttonevent int
}

// Fields returns timeseries data for influxdb
func (z *ZHASwitch) Fields() map[string]interface{} {
	return map[string]interface{}{
		"buttonevent": z.Buttonevent,
	}
}

// ZHAPresence represents aaaa
type ZHAPresence struct {
	State
	Presence bool
}

// Fields returns timeseries data for influxdb
func (z *ZHAPresence) Fields() map[string]interface{} {
	return map[string]interface{}{
		"presence": z.Presence,
	}
}

// ZHAOpenClose represents a door/window sensor that can have 2 states, open or close
type ZHAOpenClose struct {
	State
	Open bool
}

// Fields returns timeseries data for influxdb
func (z *ZHAOpenClose) Fields() map[string]interface{} {
	return map[string]interface{}{
		"open": z.Open,
	}
}

// Daylight represents a change in daylight
type Daylight struct {
	State
	Daylight bool
	Status   int
}

// Fields returns timeseries data for influxdb
func (z *Daylight) Fields() map[string]interface{} {
	return map[string]interface{}{
		"daylight": z.Daylight,
		"status":   z.Status,
	}
}

// EmptyState is an empty struct used to indicate no state was parsed
type EmptyState struct {
	State
}

// IsEmptyState returns true if s is of type EmptyState
func IsEmptyState(s interface{}) bool {
	switch s.(type) {
	case *EmptyState:
		return true
	}
	return false
}
