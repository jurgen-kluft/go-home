package main

import (
	"log"
	"time"

	"github.com/jurgen-kluft/go-home/conbee/deconz"
	"github.com/jurgen-kluft/go-home/conbee/deconz/event"
	"github.com/jurgen-kluft/go-home/config"
)

/*
STATE

State {Read} [
]

State {Write} [
]

When turning ON a light from automation logic we inform Conbee. We will keep
reading the state which will be the only factual state.

*/

type lightState struct {
	Name      string
	IDs       []string
	LastSeen  time.Time
	CT        float32
	BRI       float32
	Reachable bool
	OnOff     bool
}

type motionSensorState struct {
	Name     string
	ID       string
	LastSeen time.Time
	Motion   bool
}

type contactSensorState struct {
	Name     string
	ID       string
	LastSeen time.Time
	Contact  bool
}

type switchState struct {
	Name     string
	ID       string
	LastSeen time.Time
	Button   bool
}

type fullState struct {
	switches       map[string]switchState
	motionSensors  map[string]motionSensorState
	contactSensors map[string]contactSensorState
	lights         map[string]*lightState
}

func fullStateFromConfig(c config.ConbeeConfig) fullState {
	full := fullState{}
	full.switches = make(map[string]switchState)
	full.motionSensors = make(map[string]motionSensorState)
	full.contactSensors = make(map[string]contactSensorState)
	full.lights = make(map[string]*lightState)

	for _, e := range c.Switches {
		state := switchState{Name: e.Name, ID: e.ID, LastSeen: time.Now(), Button: false}
		full.switches[state.ID] = state
	}
	for _, e := range c.Sensors.Motion {
		state := motionSensorState{Name: e.Name, ID: e.ID, LastSeen: time.Now(), Motion: false}
		full.motionSensors[state.ID] = state
	}
	for _, e := range c.Sensors.Contact {
		state := contactSensorState{Name: e.Name, ID: e.ID, LastSeen: time.Now(), Contact: false}
		full.contactSensors[state.ID] = state
	}
	for _, e := range c.Lights {
		state := &lightState{Name: e.Name, IDs: e.IDS, LastSeen: time.Now(), Reachable: false, OnOff: false}
		for _, id := range state.IDs {
			full.lights[id] = state
		}
	}

	return full
}

func main() {
	config := defaultConfiguration()

	eventChan, err := eventChan(config.Addr, config.APIKey)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to deCONZ at %s", config.Addr)

	fullState := fullStateFromConfig(config)

	//TODO: figure out how to create a timer that is stopped
	timeout := time.NewTimer(1 * time.Second)
	timeout.Stop()

	for {

		select {
		case ev := <-eventChan:
			//fields, err := ev.Fields()
			//if err != nil {
			//	log.Printf("skip event: '%s'", err)
			//	continue
			//}

			cstate, exist := fullState.contactSensors[ev.UniqueID]
			if exist {
				dstate := ev.State.(event.ZHAOpenClose)
				log.Printf("contact:  %s -> %v = %v", cstate.Name, cstate.Contact, fields["open"].(bool))
				cstate.Contact = fields["open"].(bool)
				fullState.contactSensors[ev.UniqueID] = cstate
			} else {
				mstate, exist := fullState.motionSensors[ev.UniqueID]
				if exist {
					log.Printf("motion:  %s -> %v = %v", mstate.Name, mstate.Motion, fields["presence"].(bool))
					mstate.Motion = fields["presence"].(bool)
					fullState.motionSensors[ev.UniqueID] = mstate
				} else {
					sstate, exist := fullState.switches[ev.UniqueID]
					if exist {
						log.Printf("switch:  %s -> %v = %v", sstate.Name, sstate.Button, fields["button"].(bool))
						sstate.Button = fields["button"].(bool)
						fullState.switches[ev.UniqueID] = sstate
					} else {
						lstate, exist := fullState.lights[ev.UniqueID]
						if exist {
							log.Printf("light:  %s -> %v = %v", lstate.Name, lstate.OnOff, fields["on"].(bool))
							lstate.OnOff = fields["on"].(bool)
							fullState.lights[ev.UniqueID] = lstate
						} else {
							log.Printf("unknown:  %s", ev.UniqueID)
						}
					}
				}
			}

			timeout.Reset(1 * time.Second)

		case <-timeout.C:
			// Currently does nothing
			// Request the state of all lights?
		}
	}
}

func eventChan(addr string, APIkey string) (chan *deconz.DeviceEvent, error) {
	// get an event reader from the API
	d := deconz.API{Config: deconz.Config{Addr: addr, APIKey: APIkey}}
	reader, err := d.EventReader()
	if err != nil {
		return nil, err
	}

	// Dial the reader
	err = reader.Dial()
	if err != nil {
		return nil, err
	}

	// create a new reader, embedding the event reader
	eventReader := d.DeviceEventReader(reader)
	channel := make(chan *deconz.DeviceEvent)
	// start it, it starts its own thread
	eventReader.Start(channel)
	// return the channel
	return channel, nil
}

func defaultConfiguration() config.ConbeeConfig {
	c := config.ConbeeConfig{}
	err := config.LoadFile("../config/conbee.config.json", &c)
	if err == nil {
		log.Printf("Addr: %s, APIKey: %s", c.Addr, c.APIKey)
	} else {
		panic(err)
	}
	return c
}
