package weather

import (
	"fmt"
	"time"

	"github.com/adlio/darksky"
	"github.com/jurgen-kluft/go-home/config"
)

func converFToC(fahrenheit float64) float64 {
	return ((fahrenheit - 32.0) * 5.0 / 9.0)
}

type Client struct {
	config    *config.WeatherConfig
	location  *time.Location
	darksky   *darksky.Client
	latitude  float64
	longitude float64
	darkargs  map[string]string
	update    time.Time
}

func New() (*Client, error) {
	c := &Client{}
	c.update = time.Now()

	return c, nil
}

func (c *Client) getRainDescription(rain float64) string {
	for _, r := range c.config.Rain {
		if r.Intensity.Min <= rain && rain <= r.Intensity.Max {
			return r.Name
		}
	}
	return ""
}

func (c *Client) getCloudsDescription(clouds float64) string {
	for _, r := range c.config.Clouds {
		if r.Cover.Min <= clouds && clouds <= r.Cover.Max {
			return r.Description
		}
	}
	return ""
}

func (c *Client) getTemperatureDescription(temperature float64) string {
	for _, t := range c.config.Temperature {
		if t.Min <= temperature && temperature < t.Max {
			return t.Description
		}
	}
	return ""
}

func (c *Client) getWindDescription(wind float64) string {
	for _, w := range c.config.Wind {
		if wind < w.Speed {
			if len(w.Description) > 0 {
				return w.Description[0]
			}
			break
		}
	}
	return ""
}

type Forecast struct {
	From        time.Time `json:"from"`
	Until       time.Time `json:"until"`
	Rain        float64   `json:"rain"`
	RainDescr   string    `json:"rainDescr"`
	Wind        float64   `json:"wind"`
	WindDescr   string    `json:"windDescr"`
	Clouds      float64   `json:"clouds"`
	CloudDescr  string    `json:"cloudsDescr"`
	Temperature float64   `json:"temperature"`
	TempDescr   string    `json:"temperatureDescr"`
}

type State struct {
	Current Forecast   `json:"current"`
	Hourly  []Forecast `json:"hourly"`
}

func (c *Client) addHourly(from time.Time, until time.Time, hourly *darksky.DataBlock, state *State) {

	for _, dp := range hourly.Data {
		hfrom := time.Unix(dp.Time.Unix(), 0)
		huntil := hoursLater(hfrom, 1.0)

		if timeRangeInGlobalRange(from, until, hfrom, huntil) {
			forecast := Forecast{}
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
		chanceOfRain = "as likely as seeing a dinosaur alive."
	} else if precipProbability < 0.3 {
		chanceOfRain = "likely but probably not."
	} else if precipProbability >= 0.3 && precipProbability < 0.5 {
		chanceOfRain = "possible but you can risk it."
	} else if precipProbability >= 0.5 && precipProbability < 0.7 {
		chanceOfRain = "likely and you may want to bring an umbrella."
	} else if precipProbability >= 0.7 && precipProbability < 0.9 {
		chanceOfRain = "definitely so have an umbrella ready."
	} else {
		chanceOfRain = "for sure so open your umbrella and hold it up."
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

func (c *Client) Process() time.Duration {
	now := time.Now()
	state := &State{Hourly: []Forecast{}}

	// Weather update every 5 minutes
	if now.Unix() >= c.update.Unix() {
		c.update = time.Unix(now.Unix()+5*60, 0)

		lat := c.config.Location.Latitude
		lng := c.config.Location.Longitude
		forecast, err := c.darksky.GetForecast(fmt.Sprint(lat), fmt.Sprint(lng), c.darkargs)
		if err == nil {

			from := now
			until := hoursLater(from, 3.0)

			current := Forecast{}
			current.From = from
			current.Until = until

			current.Rain = forecast.Currently.PrecipProbability
			current.Clouds = forecast.Currently.CloudCover
			current.Wind = forecast.Currently.WindSpeed
			current.Temperature = forecast.Currently.ApparentTemperature

			state.Current = current

			c.addHourly(atHour(now, 6, 0), atHour(now, 20, 0), forecast.Hourly, state)
		}
	}

	wait := time.Duration(c.update.Unix()-time.Now().Unix()) * time.Second
	if wait < 0 {
		wait = 0
	}
	return wait
}
