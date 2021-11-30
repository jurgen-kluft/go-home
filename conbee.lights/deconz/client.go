package deconz

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var (
	// ErrMalformedResponse is returned if the deconz response JSON isn't formatted as expected
	ErrMalformedResponse = errors.New("malformed deconz response")
)

// EmptyRequest is a placeholder struct used for any request which has no parameters.
type EmptyRequest struct {
}

// Response is a generic response returned by the API
type Response []ResponseEntry

// ResponseEntry is one of the multiple response entries returned by the API
type ResponseEntry struct {
	Success map[string]interface{} `json:"success"`
	Error   ResponseError          `json:"error"`
}

// ResponseError contains a general error which was detected.
type ResponseError struct {
	Type        int    `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

// Error allows the response error to be returned as an Error compatible type.
func (re ResponseError) Error() string {
	return fmt.Sprintf("%s: %d (%s)", re.Address, re.Type, re.Description)
}

// Client represents a handle to the deconz API
type Client struct {
	httpClient *http.Client
	apiKey     string
	hostname   string
	port       int
}

// NewClient creates a new deconz API client
func NewClient(httpClient *http.Client, hostname string, port int, apiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		hostname:   hostname,
		port:       port,
		apiKey:     apiKey,
	}
}

func (c *Client) getURLBase() string {
	return "http://" + c.hostname + ":" + strconv.Itoa(c.port) + "/api/" + c.apiKey + "/"
}

func (c *Client) get(ctx context.Context, path string, respType interface{}) error {
	r, err := http.NewRequest(http.MethodGet, c.getURLBase()+path, nil)
	if err != nil {
		return err
	}

	r = r.WithContext(ctx)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return json.NewDecoder(resp.Body).Decode(respType)
	}

	deconzResp := Response{}
	err = json.NewDecoder(resp.Body).Decode(&deconzResp)
	if err != nil {
		return err
	}

	if len(deconzResp) < 1 {
		return ErrMalformedResponse
	}

	return deconzResp[0].Error
}

func (c *Client) post(ctx context.Context, path string, reqType interface{}) (*Response, error) {
	req, err := json.Marshal(reqType)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, c.getURLBase()+path, bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	r = r.WithContext(ctx)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	deconzResp := Response{}
	err = json.NewDecoder(resp.Body).Decode(&deconzResp)
	if err != nil {
		return nil, err
	}

	if len(deconzResp) < 1 {
		return nil, ErrMalformedResponse
	}
	for _, deconsRespEntry := range deconzResp {
		if len(deconsRespEntry.Success) < 1 {
			return nil, deconsRespEntry.Error
		}
	}

	return &deconzResp, nil
}

func (c *Client) put(ctx context.Context, path string, reqType interface{}) error {
	req, err := json.Marshal(reqType)
	if err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPut, c.getURLBase()+path, bytes.NewBuffer(req))
	if err != nil {
		return err
	}

	r = r.WithContext(ctx)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	deconzResp := Response{}
	err = json.NewDecoder(resp.Body).Decode(&deconzResp)
	if err != nil {
		return err
	}

	if len(deconzResp) < 1 {
		return ErrMalformedResponse
	}
	for _, deconsRespEntry := range deconzResp {
		if len(deconsRespEntry.Success) < 1 {
			return deconsRespEntry.Error
		}
	}

	return nil
}

func (c *Client) delete(ctx context.Context, path string) error {
	r, err := http.NewRequest(http.MethodDelete, c.getURLBase()+path, nil)
	if err != nil {
		return err
	}

	r = r.WithContext(ctx)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	deconzResp := Response{}
	err = json.NewDecoder(resp.Body).Decode(&deconzResp)
	if err != nil {
		return err
	}

	if len(deconzResp) < 1 {
		return ErrMalformedResponse
	}
	for _, deconsRespEntry := range deconzResp {
		if len(deconsRespEntry.Success) < 1 {
			return deconsRespEntry.Error
		}
	}

	return nil
}
