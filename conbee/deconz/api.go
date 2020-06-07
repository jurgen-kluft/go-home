package deconz

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jurgen-kluft/go-home/conbee/deconz/event"
)

// API represents the deCONZ rest api
type API struct {
	Config      Config
	deviceCache *CachedDeviceStore
}

// Devices returns a map of devices
func (a *API) Devices() (*Devices, error) {
	var resp *http.Response
	var err error
	var url string
	devices := Devices{}

	url = fmt.Sprintf("%s/%s/sensors", a.Config.Addr, a.Config.APIKey)
	resp, err = http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Devices() unable to get %s: %s", url, err)
	}
	defer resp.Body.Close()
	var sensors map[int]Device
	var dec *json.Decoder
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&sensors)
	if err != nil {
		return nil, fmt.Errorf("unable to decode deCONZ response: %s", err)
	}
	for _, s := range sensors {
		// TODO: Check if the DeviceID is formatted consistently the same way
		fmt.Printf("Sensor %s with unique-ID '%s'\n", s.Name, s.DeviceID)
		devices[s.DeviceID] = s
	}

	url = fmt.Sprintf("%s/%s/lights", a.Config.Addr, a.Config.APIKey)
	resp, err = http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Devices() unable to get %s: %s", url, err)
	}
	defer resp.Body.Close()
	var lights map[int]Device
	dec = json.NewDecoder(resp.Body)
	err = dec.Decode(&lights)
	if err != nil {
		return nil, fmt.Errorf("unable to decode deCONZ response: %s", err)
	}
	for _, l := range lights {
		// TODO: Check if the DeviceID is formatted consistently the same way
		fmt.Printf("Light %s with unique-ID '%s'\n", l.Name, l.DeviceID)
		devices[l.DeviceID] = l
	}

	return &devices, nil

}

// EventReader returns a event.Reader with a default cached type store
func (a *API) EventReader() (*event.Reader, error) {

	if a.deviceCache == nil {
		a.deviceCache = &CachedDeviceStore{DeviceGetter: a}
	}

	if a.Config.wsAddr == "" {
		err := a.Config.discoverWebsocket()
		if err != nil {
			return nil, err
		}
	}

	return &event.Reader{TypeStore: a.deviceCache, WebsocketAddr: a.Config.wsAddr}, nil
}

// DeviceEventReader takes an event reader and returns an sensor event reader
func (a *API) DeviceEventReader(r *event.Reader) *DeviceEventReader {

	if a.deviceCache == nil {
		a.deviceCache = &CachedDeviceStore{DeviceGetter: a}
	}

	return &DeviceEventReader{lookup: a.deviceCache, reader: r}
}
