package deconz

import (
	"context"
	"errors"
	"strconv"
)

// CreateGroup creates a new group on the gateway. The new ID is returned on success.
func (c *Client) CreateGroup(ctx context.Context, req *CreateGroupRequest) (int, error) {
	resp, err := c.post(ctx, "groups", req)
	if err != nil {
		return 0, err
	}

	if len(*resp) < 1 {
		return 0, errors.New("new group missing success entry")
	}
	if id, ok := (*resp)[0].Success["id"]; ok {
		if strID, ok := id.(string); ok {
			return strconv.Atoi(strID)
		}
		return 0, errors.New("new group id not string")
	}

	return 0, errors.New("new group missing id entry")
}

// GetGroups retrieves all the groups available on the gatway
func (c *Client) GetGroups(ctx context.Context) (*GetGroupsResponse, error) {
	groupsResp := &GetGroupsResponse{}

	err := c.get(ctx, "groups", groupsResp)
	if err != nil {
		return nil, err
	}

	return groupsResp, nil
}

// GetGroup retrieves the specified group
func (c *Client) GetGroup(ctx context.Context, id int) (*Group, error) {
	group := &Group{}

	err := c.get(ctx, "groups/"+strconv.Itoa(id), group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// SetGroupState specifies the new state of a group
func (c *Client) SetGroupState(ctx context.Context, id int, newState *SetGroupStateRequest) error {
	return c.put(ctx, "groups/"+strconv.Itoa(id)+"/action", newState)
}

// SetGroupConfig specifies the new config of a group
func (c *Client) SetGroupConfig(ctx context.Context, id int, newConfig *SetGroupConfigRequest) error {
	return c.put(ctx, "groups/"+strconv.Itoa(id), newConfig)
}

// DeleteGroup removes the specified group from the gateway
func (c *Client) DeleteGroup(ctx context.Context, id int) error {
	return c.delete(ctx, "groups/"+strconv.Itoa(id))
}

// Group represents a collection of lights and provides the foundation for scenes
type Group struct {
	LastAction Action   `json:"action"`
	DeviceIDs  []string `json:"devicemembership"`
	ETag       string   `json:"etag"`
	Hidden     bool     `json:"hidden"`
	ID         string   `json:"id"`
	// LightIDs contains a gateway-sorted list of all the light IDs in this group
	LightIDs []string `json:"lights"`
	// LightIDSequence contains a user-sorted list of a subset of all the light IDs in this group
	LightIDSequence []string `json:"lightsequence"`
	// MultiDeviceIDs contains the subsequent IDs of multi-device lights
	MultiDeviceIDs []string `json:"multideviceids"`
	Name           string   `json:"name"`
	Scenes         []struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		TransitionTime int    `json:"transitiontime"`
		LightCount     int    `json:"lightcount"`
	} `json:"scenes"`
	State GroupState `json:"state"`
}

// GroupState contains the fields relevant to the state of a group
type GroupState struct {
	AllOn bool `json:"all_on"`
	AnyOn bool `json:"any_on"`
}

// Action represents a state change which has occurred
type Action struct {
	On                *bool      `json:"on"`
	Brightness        int        `json:"bri"`
	Hue               *int       `json:"hue"`
	Saturation        *int       `json:"sat"`
	ColourTemperature int        `json:"ct"`
	XY                *[]float64 `json:"xy"`
	Effect            *string    `json:"effect"`
}

// CreateGroupRequest is used to create a new group with the specified name.
type CreateGroupRequest struct {
	Name string `json:"name"`
}

// GetGroupsResponse contains the fields returned by the 'list groups' API call
type GetGroupsResponse map[string]Group

// SetGroupConfigRequest sets the config options of the group
type SetGroupConfigRequest struct {
	Name            string   `json:"name,omitempty"`
	LightIDs        []string `json:"lights,omitempty"`
	Hidden          bool     `json:"hidden,omitempty"`
	LightIDSequence []string `json:"lightsequence,omitempty"`
	MultiDeviceIDs  []string `json:"multideviceids,omitempty"`
}

// SetGroupStateRequest sets the state of the specified group.
type SetGroupStateRequest struct {
	SetLightStateRequest
	// Toggle flips the state from on to off or vice versa.
	// This superscedes the values set directly.
	Toggle bool `json:"toggle,omitempty"`
}
