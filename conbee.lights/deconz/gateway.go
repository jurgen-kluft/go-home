package deconz

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// CreateAPIKey attempts to generate an API key to use for subequent operations.
// The username is optional - if not specified a random string will be generated.
func (c *Client) CreateAPIKey(ctx context.Context, reqType *CreateAPIKeyRequest) (string, error) {
	req, err := json.Marshal(reqType)
	if err != nil {
		return "", err
	}

	r, err := http.NewRequest(http.MethodPost, "http://"+c.hostname+":"+strconv.Itoa(c.port)+"/api", bytes.NewBuffer(req))
	if err != nil {
		return "", err
	}

	r = r.WithContext(ctx)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	deconzResp := &Response{}
	err = json.NewDecoder(resp.Body).Decode(deconzResp)
	if err != nil {
		return "", err
	}

	if len(*deconzResp) < 1 {
		return "", ErrMalformedResponse
	}
	if id, ok := (*deconzResp)[0].Success["username"]; ok {
		if strID, ok := id.(string); ok {
			return strID, nil
		}
		return "", errors.New("new user id not string")
	}

	return "", (*deconzResp)[0].Error
}

// DeleteAPIKey is used to delete the specified API key.
func (c *Client) DeleteAPIKey(ctx context.Context, keyToDelete string) error {
	return c.delete(ctx, "config/whitelist/"+keyToDelete)
}

// GetGatewayState collects the current state from the gateway and returns it
func (c *Client) GetGatewayState(ctx context.Context) (*GatewayState, error) {
	gwState := &GatewayState{}

	err := c.get(ctx, "config", gwState)
	if err != nil {
		return nil, err
	}

	return gwState, nil
}

// GatewayState contains the current state of the gateway
type GatewayState struct {
	APIVersion          string              `json:"apiversion"`
	SoftwareVersion     string              `json:"swversion"`
	SoftwareUpdateState SoftwareUpdateState `json:"swupdate"`

	MACAddress    string `json:"mac"`
	ZigbeeChannel int    `json:"zigbeechannel"`
	ZigbeePANID   int    `json:"panid"`
	GatewayID     string `json:"uuid"`

	WebsocketNotifyAll bool `json:"websocketnotifyall"`
	WebsocketPort      int  `json:"websocketport"`
	LinkButtonPressed  bool `json:"linkbutton"`

	Name       string `json:"name"`
	LocalTime  string `json:"localtime"`
	UTCTime    string `json:"utc"`
	TimeFormat string `json:"timeformat"`
	Timezone   string `json:"timezone"`

	UsingDHCP bool   `json:"dhcp"`
	GatewayIP string `json:"gateway"`
	IP        string `json:"ipaddress"`
	Netmask   string `json:"netmask"`
}

// SoftwareUpdateState contains the important data about the current software update profile of the gateway
type SoftwareUpdateState struct {
	Notify      bool   `json:"notify"`
	Text        string `json:"text"`
	UpdateState int    `json:"updatestate"`
	URL         string `json:"url"`
}

// GetGatewayResponse contains the returned data from the full gateway API call
type GetGatewayResponse struct {
	GatewayState GatewayState `json:"config"`
	//Groups map[string]Group `json:"groups`
	Lights GetLightsResponse `json:"lights"`
	// Rules map[string]Rule `json:"rules"`
	// Schedules map[string]Schedule `json:"schedules"`
	Sensors map[int]Sensor `json:"sensors"`
}

// CreateAPIKeyRequest contains the fields which will be used to request an API key.
type CreateAPIKeyRequest struct {
	ApplicationName string `json:"devicetype"`
	Username        string `json:"username,omitempty"`
}

// SetConfigRequest contains the set of possible gateway configuration parameters.
type SetConfigRequest struct {
	Name        string `json:"name,omitempty"`
	RFConnected bool   `json:"rfconnected,omitempty"`
	// UpdateChannel can be set to one of stable, alpha, beta
	UpdateChannel string `json:"updatechannel,omitempty"`
	// PermitJoin when set to 0 indicates no Zigbee devices can join
	// 255 means the network is open
	// 1..254 represents the time in seconds the network will be open
	// These values decrement automatically
	PermitJoin int `json:"permitjoin,omitempty"`
	// GroupDelay contains the time between two group commands, in milliseconds
	GroupDelay        int  `json:"groupdelay,omitempty"`
	OTAUActive        bool `json:"otauactive,omitempty"`
	GWDiscoveryActive bool `json:"discovery,omitempty"`
	// Unlock being set to a value > 0 (and less than 600, the max) indicates the number of seconds the gateway is open for pairing
	Unlock int `json:"unlock,omitempty"`
	// ZigbeeChannel specifies one of 11, 15, 20 or 25 (the valid Zigbee channel numbers)
	ZigbeeChannel int    `json:"zigbeechannel,omitempty"`
	Timezone      string `json:"timezone,omitempty"`
	UTC           string `json:"utc,omitempty"`
	// TimeFormat is specified as either 12h or 24h
	TimeFormat string `json:"timeformat,omitempty"`
}
