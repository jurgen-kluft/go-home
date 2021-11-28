package deconz

import (
	"context"
)

// GetLights retrieves all the lights available on the gatway
func (c *Client) GetLights(ctx context.Context) (GetLightsResponse, error) {
	lightsResp := GetLightsResponse{}

	err := c.get(ctx, "lights", &lightsResp)
	if err != nil {
		return nil, err
	}

	return lightsResp, nil
}

// GetLight retrieves the specified light
func (c *Client) GetLight(ctx context.Context, id string) (*Light, error) {
	light := &Light{}

	err := c.get(ctx, "lights/"+id, light)
	if err != nil {
		return nil, err
	}

	return light, nil
}

// SetLightState specifies the new state of a light
func (c *Client) SetLightState(ctx context.Context, id string, newState *SetLightStateRequest) error {
	return c.put(ctx, "lights/"+id+"/state", newState)
}

// SetLightConfig specifies the new config of a light
func (c *Client) SetLightConfig(ctx context.Context, id string, newConfig *SetLightConfigRequest) error {
	return c.put(ctx, "lights/"+id, newConfig)
}

// DeleteLight removes the specified light from the gateway
func (c *Client) DeleteLight(ctx context.Context, id string) error {
	return c.delete(ctx, "lights/"+id)
}

// DeleteLightGroups removes the light from all its groups
func (c *Client) DeleteLightGroups(ctx context.Context, id string) error {
	return c.delete(ctx, "lights/"+id+"/groups")
}

// DeleteLightScenes removes the light from all its scenes
func (c *Client) DeleteLightScenes(ctx context.Context, id string) error {
	return c.delete(ctx, "lights/"+id+"/scenes")
}

// Light contains the fields of a light.
type Light struct {
	// ID contains the gateway-specified ID; could change.
	// Exists only for accessing by path; dedup using UniqueID instead
	ID              string
	CTMax           int        `json:"ctmax"`
	CTMin           int        `json:"ctmin"`
	LastAnnounced   string     `json:"lastannounced"`
	LastSeen        string     `json:"lastseen"`
	ETag            string     `json:"etag"`
	Manufacturer    string     `json:"manufacturer"`
	Name            string     `json:"name"`
	ModelID         string     `json:"modelid"`
	SoftwareVersion string     `json:"swversion"`
	Type            string     `json:"type"`
	State           LightState `json:"state"`
	UniqueID        string     `json:"uniqueid"`
}

// LightState contains the specific, controllable fields of a light.
type LightState struct {
	On         bool      `json:"on"`
	Brightness int       `json:"bri"`
	Hue        int       `json:"hue"`
	Saturation int       `json:"sat"`
	CT         int       `json:"ct"`
	XY         []float64 `json:"xy"`
	Alert      string    `json:"alert"`
	ColorMode  string    `json:"colormode"`
	Effect     string    `json:"effect"`
	Reachable  bool      `json:"reachable"`
}

// GetLightsResponse contains the result of all active lights.
type GetLightsResponse map[string]Light

// SetLightStateRequest lets a user update certain properties of the light.
// These are directly changing the active light and what it is showing.
type SetLightStateRequest struct {
	On         bool      `json:"on"`
	Brightness int       `json:"bri,omitempty"`
	Hue        int       `json:"hue,omitempty"`
	Saturation int       `json:"sat,omitempty"`
	CT         int       `json:"ct,omitempty"`
	XY         []float64 `json:"xy,omitempty"`
	Alert      string    `json:"alert,omitempty"`
	// Effect contains the light effect to apply. Either 'none' or 'colorloop'
	Effect string `json:"effect,omitempty"`
	// ColorLoopSpeed contains the speed of a colorloop.
	// 1 is very fast, 255 is very slow.
	// This is only read if the 'colorloop' effect is specifed
	ColorLoopSpeed int `json:"colorloopspeed,omitempty"`
	// TransitionTime is represented in 1/10th of a second between states
	TransitionTime int `json:"transitiontime,omitempty"`
}

// SetLightConfigRequest lets a user update certain properties of the light.
// This is metadata and not directly changing the light behaviour.
type SetLightConfigRequest struct {
	Name string `json:"name,omitempty"`
}
