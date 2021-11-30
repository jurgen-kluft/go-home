package deconz

import "encoding/json"

// WebsocketUpdate contains the data deserialized from the async channel
type WebsocketUpdate struct {
	Meta WebsocketUpdateMetadata

	// These are conditionally filled in by parsing the State json.RawMessage field
	GroupState  *GroupState
	LightState  *LightState
	SensorState *SensorState

	// These are conditionally filled in by parsing the relevant json.RawMessage field
	Group  *Group
	Light  *Light
	Sensor *Sensor
}

// WebsocketUpdateMetadata contains the common metadata fields about the update.
type WebsocketUpdateMetadata struct {
	Type       string `json:"t"`
	Event      string `json:"e"`
	Resource   string `json:"r"`
	ResourceID string `json:"id"`
	UniqueID   string `json:"uniqueid"`

	// The following are set on `changed` events
	Config json.RawMessage `json:"config"`
	Name   string          `json:"name"`
	State  json.RawMessage `json:"state"`

	// The following fields are only set on `scene-called` events
	GroupID string `json:"gid"`
	SceneID string `json:"scid"`

	// The following fields are set on the `added` event for the relevant resource type
	Group  json.RawMessage `json:"group"`
	Light  json.RawMessage `json:"light"`
	Sensor json.RawMessage `json:"sensor"`
}

// UnmarshalJSON allows us to conditionally deserialize the websocket update
// so that only the relevant fields are available.
func (wsu *WebsocketUpdate) UnmarshalJSON(b []byte) error {
	meta := WebsocketUpdateMetadata{}
	err := json.Unmarshal(b, &meta)
	if err != nil {
		return err
	}

	wsu.Meta = meta

	if meta.Resource == "sensors" {
		if meta.Event == "changed" {
			state := &SensorState{}
			err = json.Unmarshal(meta.State, state)
			if err != nil {
				return err
			}

			wsu.SensorState = state
		} else if meta.Event == "added" {
			sensor := &Sensor{}
			err = json.Unmarshal(meta.Sensor, sensor)
			if err != nil {
				return err
			}

			wsu.Sensor = sensor
		}
	} else if meta.Resource == "lights" {
		if meta.Event == "changed" {
			state := &LightState{}
			err = json.Unmarshal(meta.State, state)
			if err != nil {
				return err
			}

			wsu.LightState = state
		} else if meta.Event == "added" {
			light := &Light{}
			err = json.Unmarshal(meta.Light, light)
			if err != nil {
				return err
			}

			wsu.Light = light
		}
	} else if meta.Resource == "groups" {
		if meta.Event == "changed" {
			state := &GroupState{}
			err = json.Unmarshal(meta.State, state)
			if err != nil {
				return err
			}

			wsu.GroupState = state
		} else if meta.Event == "added" {
			group := &Group{}
			err = json.Unmarshal(meta.Group, group)
			if err != nil {
				return err
			}

			wsu.Group = group
		}
	}

	return nil
}
