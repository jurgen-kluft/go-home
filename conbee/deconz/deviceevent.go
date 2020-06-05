package deconz

import (
	"fmt"

	"github.com/jurgen-kluft/go-home/conbee/deconz/event"
)

// DeviceEvent is a sensor and a event embedded
type DeviceEvent struct {
	*Device
	*event.Event
}

type fielder interface {
	Fields() map[string]interface{}
}

// Fields returns tags and fields
func (s *DeviceEvent) Fields() (map[string]interface{}, error) {
	f, ok := s.Event.State.(fielder)
	if !ok {
		return nil, fmt.Errorf("this event (%T:%s) has no time series data", s.State, s.Name)
	}
	return f.Fields(), nil
}
