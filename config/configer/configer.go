package main

import (
	"io/ioutil"

	"github.com/jurgen-kluft/go-home/config"
)

// Configs holds all the config objects that we can have
type Configs struct {
	aqiConfig      *config.AqiConfig
	calendarConfig *config.CalendarConfig
	fluxConfig     *config.FluxConfig
	presenceConfig *config.PresenceConfig
	sensorState    *config.SensorState
	shoutConfig    *config.ShoutConfig
	suncalcConfig  *config.SuncalcConfig
	weatherConfig  *config.WeatherConfig
	wemoConfig     *config.WemoConfig
	xiaomiConfig   *config.XiaomiConfig
	yeeConfig      *config.YeeConfig
}

// Load will load the JSON based config file
func (c *Configs) Load(configtype string) (err error) {
	if configtype == "aqi" {
		data, err := ioutil.ReadFile("aqi.config.json")
		if err == nil {
			c.aqiConfig, err = config.AqiConfigFromJSON(string(data))
		}
	}

	return err
}

func main() {

	// The idea here is to be able to hot load updated configurations and send them
	// on the associated pubsub channel.

}
