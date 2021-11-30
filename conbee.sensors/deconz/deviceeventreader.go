package deconz

import (
	"errors"
	"log"
	"time"

	"github.com/jurgen-kluft/go-home/conbee.sensors/deconz/event"
)

// DeviceLookup represents an interface for device lookup
type DeviceLookup interface {
	LookupDevice(string) (*Device, error)
	SupportsResource(string) bool
}

// EventReader interface
type EventReader interface {
	ReadEvent() (*event.Event, error)
	Dial() error
	Close() error
}

// DeviceEventReader reads events from an event.reader and returns DeviceEvents
type DeviceEventReader struct {
	lookup  DeviceLookup
	reader  EventReader
	running bool
}

// Start starts a thread reading events into the given channel
// returns immediately
func (r *DeviceEventReader) Start(out chan *DeviceEvent) error {

	if r.lookup == nil {
		return errors.New("cannot run without a SensorLookup from which to lookup sensors")
	}
	if r.reader == nil {
		return errors.New("cannot run without a EventReader from which to read events")
	}

	if r.running {
		return errors.New("reader is already running")
	}

	r.running = true

	go func() {
	REDIAL:
		for r.running {
			// establish connection
			for r.running {
				err := r.reader.Dial()
				if err != nil {
					log.Printf("Error connecting Deconz websocket: %s\nAttempting reconnect in 5s...", err)
					time.Sleep(5 * time.Second) // TODO configurable delay
				} else {
					log.Printf("Deconz websocket connected")
					break
				}
			}
			// read events until connection fails
			for r.running {
				e, err := r.reader.ReadEvent()
				if err != nil {
					if eerr, ok := err.(event.EventError); ok && eerr.Recoverable() {
						log.Printf("Dropping event due to error: %s", err)
						continue
					}
					continue REDIAL
				}
				// we only care about sensor events
				if !r.lookup.SupportsResource(e.Resource) {
					log.Printf("unsupported resource %s with id: %d, type: %s, event: %s, rawstate: %s", e.Resource, e.ID, e.Type, e.Event, e.RawState)
					continue
				}

				// TODO: Check if the UniqueID is formatted consistently the same way
				device, err := r.lookup.LookupDevice(e.UniqueID)

				if err != nil {
					log.Printf("Dropping event. Could not lookup device for id %s: %s", e.UniqueID, err)
					continue
				}
				// send event on channel
				out <- &DeviceEvent{Event: e, Device: device}
			}
		}
		// if not running, close connection and return from goroutine
		r.reader.Close()
		log.Printf("Deconz websocket closed")
	}()
	return nil
}

// StopReadEvents closes the reader, closing the connection to deconz and terminating the goroutine
func (r *DeviceEventReader) StopReadEvents() {
	r.running = false
}
