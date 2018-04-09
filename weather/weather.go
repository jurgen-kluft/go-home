package weather

import (
	"fmt"
	"time"

	"github.com/adlio/darksky"
	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
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

func New() *Client {
	c := &Client{}
	c.darkargs = map[string]string{}
	c.update = time.Now()
	return c
}

func (c *Client) getRainDescription(rain float64) string {
	for _, r := range c.config.Rain {
		if r.Range.Min <= rain && rain <= r.Range.Max {
			return r.Name
		}
	}
	return ""
}

func (c *Client) getCloudsDescription(clouds float64) string {
	for _, r := range c.config.Clouds {
		if r.Range.Min <= clouds && clouds <= r.Range.Max {
			return r.Description
		}
	}
	return ""
}

func (c *Client) getTemperatureDescription(temperature float64) string {
	for _, t := range c.config.Temperature {
		if t.Range.Min <= temperature && temperature < t.Range.Max {
			return t.Description
		}
	}
	return ""
}

func (c *Client) getWindDescription(wind float64) string {
	for _, w := range c.config.Wind {
		if w.Range.Min <= wind && wind < w.Range.Max {
			return w.Description
		}
	}
	return ""
}

func (c *Client) addHourly(from time.Time, until time.Time, hourly *darksky.DataBlock) {

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

func (c *Client) Process(client *pubsub.Context) time.Duration {
	now := time.Now()

	state := config.NewSensorState("weather")

	// Weather update every 5 minutes
	if now.Unix() >= c.update.Unix() {
		c.update = time.Unix(now.Unix()+5*60, 0)

		lat := c.config.Location.Latitude
		lng := c.config.Location.Longitude
		forecast, err := c.darksky.GetForecast(fmt.Sprint(lat), fmt.Sprint(lng), c.darkargs)
		if err == nil {

			from := now
			until := hoursLater(from, 3.0)

			state.AddTimeWndAttr("rain", from, until)
			state.AddFloatAttr("rain", forecast.Currently.PrecipProbability)
			state.AddFloatAttr("clouds", forecast.Currently.CloudCover)
			state.AddFloatAttr("wind", forecast.Currently.WindSpeed)
			state.AddFloatAttr("temperature", forecast.Currently.ApparentTemperature)

			//			c.addHourly(atHour(now, 6, 0), atHour(now, 20, 0), forecast.Hourly, state)
		}
	}

	jsonstr, err := state.ToJSON()
	if err == nil {
		client.Publish("state/sensor/weather/", jsonstr)
	}

	wait := time.Duration(c.update.Unix()-time.Now().Unix()) * time.Second
	if wait < 0 {
		wait = 0
	}
	return wait
}

func main() {
	weather := New()

	logger := logpkg.New("shout")
	logger.AddEntry("emitter")
	logger.AddEntry("shout")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/weather/", "state/sensor/weather/"}
		subscribe := []string{"config/weather/"}
		err := client.Connect("weather", register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/weather/" {
						if weather.config == nil {
							logger.LogInfo("weather", "received configuration")
							weather.config, err = config.WeatherConfigFromJSON(string(msg.Payload()))
							if err != nil {
								weather.config = nil
							}
						}
					} else if topic == "client/disconnected/" {
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}

				case <-time.After(time.Second * 60):
					if weather.config != nil {
						weather.Process(client)
					}
				}
			}
		}

		if err != nil {
			logger.LogError("shout", err.Error())
		}

		time.Sleep(5 * time.Second)
	}

}
