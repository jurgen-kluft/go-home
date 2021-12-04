package deconz

import (
	"errors"
	"fmt"
)

// CachedDeviceStore is a cached typestore which provides LookupType for event passing
// it will be our default store
type CachedDeviceStore struct {
	DeviceGetter
	cache *Devices
}

// DeviceGetter defines how we like to ask for devices
type DeviceGetter interface {
	Devices() (*Devices, error)
}

// SupportsResource returns true if this Store supports the resource type
func (c *CachedDeviceStore) SupportsResource(restype string) bool {
	// TODO: determine the resource types we are interested in
	return restype != "Unknown"
}

// LookupType lookups deCONZ event types though a cache
// TODO: if we where unable to lookup an ID we should try to refetch the cache
// - there could have been an device added we dont know about
func (c *CachedDeviceStore) LookupType(id string) (string, error) {
	var err error
	if c.cache == nil {
		err = c.populateCache()
		if err != nil {
			return "", fmt.Errorf("unable to populate devices: %s", err)
		}
	}

	if s, found := (*c.cache)[id]; found {
		return s.Type, nil
	}

	return "", errors.New("no such device")
}

// LookupDevice returns a device for an device id
func (c *CachedDeviceStore) LookupDevice(id string) (*Device, error) {
	var err error
	if c.cache == nil {
		err = c.populateCache()
		if err != nil {
			return nil, fmt.Errorf("unable to populate devices: %s", err)
		}
	}

	if s, found := (*c.cache)[id]; found {
		return &s, nil
	}

	return nil, errors.New("no such device")
}

func (c *CachedDeviceStore) populateCache() error {
	var err error
	c.cache, err = c.Devices()
	if err != nil {
		return err
	}

	//log.Printf("DeviceStore updated, found %d devices", len((*c.cache)))

	return nil
}
