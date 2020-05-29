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

State {Read-Only} [
Bedroom Motion Sensor
Kitchen Motion Sensor
Livingroom Motion Sensor 1
Livingroom Motion Sensor 2
Frontdoor Magnet Sensor
]

State {Read/Write} [
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

When turning ON a light from automation logic we modify the state but we do keep
reading the state which will take higher priority.


*/


func main() {
	config := defaultConfiguration()

	deconzConfig := deconz.Config{Addr: config.Addr, APIKey: config.APIKey}
	sensorChan, err := sensorEventChan(deconzConfig)
	if err != nil {
		panic(err)
	}

	log.Printf("Connected to deCONZ at %s", deconzConfig.Addr)

	//TODO: figure out how to create a timer that is stopped
	timeout := time.NewTimer(1 * time.Second)
	timeout.Stop()

	for {

		select {
		case sensorEvent := <-sensorChan:
			_, fields, err := sensorEvent.Timeseries()

			if err != nil {
				//log.Printf("skip event: '%s'", err)
				continue
			}

			for k, v := range fields {
				if strings.HasPrefix(k, "presence") {
					log.Printf("motion:  %s -> %s = %v (uuid: %d)", sensorEvent.Name, k, v, sensorEvent.ID)
				} else if strings.HasPrefix(k, "open") {
					log.Printf("magnet:  %s -> %s = %v (uuid: %d)", sensorEvent.Name, k, v, sensorEvent.ID)
				} else if strings.HasPrefix(k, "button") {
					log.Printf("switch:  %s -> %s = %v (uuid: %d)", sensorEvent.Name, k, v, sensorEvent.ID)
				}
			}

			timeout.Reset(1 * time.Second)

		case <-timeout.C:
			// when timer fires: save batch points, initialize a new batch
			// err := influxdb.Write(batch)
			// if err != nil {
			// 	panic(err)
			// }

			// log.Printf("Saved %d records to influxdb", len(batch.Points()))
			// // influx batch point
			// batch, err = client.NewBatchPoints(client.BatchPointsConfig{
			// 	Database:  config.InfluxdbDatabase,
			// 	Precision: "s",
			// })
		}
	}
}

func sensorEventChan(c deconz.Config) (chan *deconz.SensorEvent, error) {
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
	sensorEventReader := d.SensorEventReader(reader)
	channel := make(chan *deconz.SensorEvent)
	// start it, it starts its own thread
	sensorEventReader.Start(channel)
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
