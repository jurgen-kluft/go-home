package deconz

import (
	"github.com/jurgen-kluft/go-home/conbee.sensors/deconz/event"
)

// DeviceEvent is a sensor and a event embedded
type DeviceEvent struct {
	*Device
	*event.Event
}
