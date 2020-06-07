package main

import (
	"log"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-home/conbee/deconz"
	"github.com/jurgen-kluft/go-home/config"
)

/*
STATE

State {Read} [
Bedroom Motion Sensor
Kitchen Motion Sensor
Livingroom Motion Sensor 1
Livingroom Motion Sensor 2
Frontdoor Magnet Sensor

Sophia Switch
Bedroom Switch
Remote Switch

Bedroom Light Stand
Bedroom Light Main
Kitchen Light Main
Jennifer Light Main
Sophia Light Main
Sophia Light Stand
Livingroom Light Main
Livingroom Light Stand
Livingroom Light Chandelier
]

State {Write} [
Bedroom Light Stand
Bedroom Light Main
Kitchen Light Main
Jennifer Light Main
Sophia Light Main
Sophia Light Stand
Livingroom Light Main
Livingroom Light Stand
Livingroom Light Chandelier
]

When turning ON a light from automation logic we inform Conbee. We will keep
reading the state which will be the only factual state.

*/

type lightState struct {
	CT        float32
	BRI       float32
	Reachable bool
	OnOff     bool
}

type motionSensorState struct {
	LastSeen time.Time
	Motion   bool
}

type contactSensorState struct {
	LastSeen time.Time
	Contact  bool
}

type fullstate struct {
	lights         map[string]lightState
	motionSensors  map[string]motionSensorState
	contactSensors map[string]contactSensorState
	rename         map[string]string
}

func main() {
	config := defaultConfiguration()

	deconzConfig := deconz.Config{Addr: config.Addr, APIKey: config.APIKey}
	eventChan, err := eventChan(deconzConfig)
	if err != nil {
		panic(err)
	}

	log.Printf("Connected to deCONZ at %s", deconzConfig.Addr)

	//TODO: figure out how to create a timer that is stopped
	timeout := time.NewTimer(1 * time.Second)
	timeout.Stop()

	for {

		select {
		case ev := <-eventChan:
			fields, err := ev.Fields()

			if err != nil {
				//log.Printf("skip event: '%s'", err)
				continue
			}

			for k, v := range fields {
				if strings.HasPrefix(k, "presence") {
					log.Printf("motion:  %s -> %s = %v (uuid: %s)", ev.Name, k, v, ev.UniqueID)
				} else if strings.HasPrefix(k, "open") {
					log.Printf("magnet:  %s -> %s = %v (uuid: %s)", ev.Name, k, v, ev.UniqueID)
				} else if strings.HasPrefix(k, "button") {
					log.Printf("switch:  %s -> %s = %v (uuid: %s)", ev.Name, k, v, ev.UniqueID)
				} else if strings.HasPrefix(k, "bri") {
					log.Printf("light:  %s -> %s = %v (uuid: %s)", ev.Name, k, v, ev.UniqueID)
				}
			}

			timeout.Reset(1 * time.Second)

		case <-timeout.C:
			// Currently does nothing
			// Request the state of all lights?
		}
	}
}

func eventChan(c deconz.Config) (chan *deconz.DeviceEvent, error) {
	// get an event reader from the API
	d := deconz.API{Config: c}
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

func defaultConfiguration() *config.ConbeeConfig {
	// this is the default configuration
	c := &config.ConbeeConfig{
		Addr:   "http://10.0.0.18/api",
		APIKey: "0A498B9909",
	}

	return c
}
