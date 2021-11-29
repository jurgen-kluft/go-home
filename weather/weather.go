package main

import (
	"fmt"
	"time"

	"github.com/adlio/darksky"
	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
)

func converFToC(fahrenheit float64) float64 {
	return ((fahrenheit - 32.0) * 5.0 / 9.0)
}

type instance struct {
	config    *config.WeatherConfig
	location  *time.Location
	darksky   *darksky.Client
	darkargs  map[string]string
	latitude  float64
	longitude float64
	update    time.Time
}

func new() *instance {
	c := &instance{}
	c.darkargs = map[string]string{}
	c.darksky = nil
	c.darkargs = map[string]string{}
	c.darkargs["units"] = "si"
	c.update = time.Now()
	return c
}

func (c *instance) initialize(config *config.WeatherConfig) {
	c.config = config
	c.darksky = darksky.NewClient(c.config.Key.String)
	c.latitude = c.config.Location.Latitude
	c.longitude = c.config.Location.Longitude
}

func (c *instance) getRainDescription(rain float64) string {
	for _, r := range c.config.Rain {
		if r.Range.Min <= rain && rain <= r.Range.Max {
			return r.Name
		}
	}
	return ""
}

func (c *instance) getCloudsDescription(clouds float64) string {
	for _, r := range c.config.Clouds {
		if r.Range.Min <= clouds && clouds <= r.Range.Max {
			return r.Description
		}
	}
	return ""
}

func (c *instance) getTemperatureDescription(temperature float64) string {
	for _, t := range c.config.Temperature {
		if t.Range.Min <= temperature && temperature < t.Range.Max {
			return t.Description
		}
	}
	return ""
}

func (c *instance) getWindDescription(wind float64) string {
	for _, w := range c.config.Wind {
		if w.Range.Min <= wind && wind < w.Range.Max {
			return w.Description
		}
	}
	return ""
}

func (c *instance) addHourly(from time.Time, until time.Time, hourly *darksky.DataBlock) {

	for _, dp := range hourly.Data {
		hfrom := time.Unix(dp.Time.Unix(), 0)
		huntil := hoursLater(hfrom, 1.0)

		if timeRangeInGlobalRange(from, until, hfrom, huntil) {
			forecast := config.Forecast{}
			forecast.From = hfrom
			forecast.Until = huntil

			forecast.Rain = dp.PrecipProbability
			forecast.Clouds = dp.CloudCover
			forecast.Wind = dp.WindSpeed
			forecast.Temperature = dp.ApparentTemperature

			forecast.RainDescr = c.getRainDescription(dp.PrecipProbability)
			forecast.CloudDescr = c.getCloudsDescription(dp.CloudCover)
			forecast.WindDescr = c.getWindDescription(dp.WindSpeed)
			forecast.TempDescr = c.getTemperatureDescription(dp.ApparentTemperature)

		}
	}
}

func timeRangeInGlobalRange(globalFrom time.Time, globalUntil time.Time, from time.Time, until time.Time) bool {
	gf := globalFrom.Unix()
	gu := globalUntil.Unix()
	f := from.Unix()
	u := until.Unix()
	return f >= gf && f < gu && u > gf && u <= gu
}

func chanceOfRain(from time.Time, until time.Time, hourly *darksky.DataBlock) (chanceOfRain string) {

	precipProbability := 0.0
	for _, dp := range hourly.Data {
		hfrom := time.Unix(dp.Time.Unix(), 0)
		huntil := hoursLater(hfrom, 1.0)
		if timeRangeInGlobalRange(from, until, hfrom, huntil) {
			if dp.PrecipProbability > precipProbability {
				precipProbability = dp.PrecipProbability
			}
		}
	}

	// Finished the sentence:
	// "The chance of rain is " +
	if precipProbability < 0.1 {
		chanceOfRain = "none, as likely as seeing a dinosaur alive."
	} else if precipProbability < 0.3 {
		chanceOfRain = "unlikely, but probably."
	} else if precipProbability >= 0.3 && precipProbability < 0.5 {
		chanceOfRain = "possible, you can risk it."
	} else if precipProbability >= 0.5 && precipProbability < 0.7 {
		chanceOfRain = "likely, you may want to bring an umbrella."
	} else if precipProbability >= 0.7 && precipProbability < 0.9 {
		chanceOfRain = "definitely, so have an umbrella ready."
	} else {
		chanceOfRain = "for sure, so open your umbrella and hold it up."
	}

	return
}

const (
	daySeconds = 60.0 * 60.0 * 24.0
)

func hoursLater(date time.Time, h float64) time.Time {
	return time.Unix(date.Unix()+int64(h*float64(daySeconds)/24.0), 0)
}

func atHour(date time.Time, h int, m int) time.Time {
	now := time.Date(date.Year(), date.Month(), date.Day(), h, m, 0, 0, date.Location())
	return now
}

func (c *instance) process(name string) ([]byte, error) {
	now := time.Now()

	// Weather update every 5 minutes
	if now.Unix() >= c.update.Unix() {
		c.update = time.Unix(now.Unix()+5*60, 0)

		state := config.NewSensorState(name, "weather")

		lat := c.latitude
		lng := c.longitude
		forecast, err := c.darksky.GetForecast(fmt.Sprint(lat), fmt.Sprint(lng), c.darkargs)
		if err == nil {

			from := now
			until := hoursLater(from, 3.0)

			state.AddTimeWndAttr("forecast", from, until)
			state.AddFloatAttr("rain", forecast.Currently.PrecipProbability)
			state.AddFloatAttr("clouds", forecast.Currently.CloudCover)
			state.AddFloatAttr("wind", forecast.Currently.WindSpeed)
			state.AddFloatAttr("temperature", forecast.Currently.ApparentTemperature)

			//			c.addHourly(atHour(now, 6, 0), atHour(now, 20, 0), forecast.Hourly, state)
		}
		jsonbytes, err := state.ToJSON()
		return jsonbytes, err
	}
	return nil, nil
}

func main() {
	c := new()
	register := []string{"config/weather/", "config/request/", "state/sensor/weather/"}
	subscribe := []string{"config/weather/"}

	m := microservice.New("weather")
	m.RegisterAndSubscribe(register, subscribe)

	m.RegisterHandler("config/weather/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		weatherConfig, err := config.WeatherConfigFromJSON(msg)
		if err == nil {
			c.initialize(weatherConfig)
		} else {
			m.Logger.LogError(m.Name, err.Error())
		}
		return true
	})

	m.RegisterHandler("*", func(m *microservice.Service, topic string, msg []byte) bool {
		fmt.Printf("message received, topic:'%s', msg:'%s'\n", topic, string(msg))
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if tickCount%5 == 0 {
			if c.darksky != nil {
				jsonbytes, err := c.process(m.Name)
				if err == nil {
					if jsonbytes != nil {
						m.Pubsub.Publish("state/sensor/weather/", jsonbytes)
					}
				} else {
					m.Logger.LogError(m.Name, err.Error())
				}
			} else {
				m.Pubsub.PublishStr("config/request/", m.Name)
			}
		}
		tickCount++
		return true
	})

	m.Loop()
}
