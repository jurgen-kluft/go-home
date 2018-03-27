package weather

import (
	"fmt"
	"time"

	"github.com/adlio/darksky"
	"github.com/jurgen-kluft/hass-go/dynamic"
	"github.com/jurgen-kluft/hass-go/state"
	"github.com/spf13/viper"
)

func converFToC(fahrenheit float64) float64 {
	return ((fahrenheit - 32.0) * 5.0 / 9.0)
}

type cloud struct {
	name        string
	description string
	min         float64
	max         float64
}

type rain struct {
	name          string
	unit          string
	intensity_min float64
	intensity_max float64
}

type wind struct {
	unit        string
	speed       float64
	description []string
}

type temperature struct {
	unit        string
	min         float64
	max         float64
	description string
}

type Client struct {
	viper        *viper.Viper
	location     *time.Location
	darksky      *darksky.Client
	latitude     float64
	longitude    float64
	darkargs     map[string]string
	clouds       []cloud
	rains        []rain
	winds        []wind
	temperatures []temperature
	update       time.Time
}

func New() (*Client, error) {
	c := &Client{}
	c.viper = viper.New()
	c.update = time.Now()

	// Viper command-line package
	c.viper.SetConfigName("weather") // name of config file (without extension)
	c.viper.AddConfigPath("config/") // optionally look for config in the working directory
	err := c.viper.ReadInConfig()    // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		return nil, err
	}

	c.location, _ = time.LoadLocation(c.viper.GetString("location.timezone"))
	c.darksky = darksky.NewClient(c.viper.GetString("darksky.key"))
	c.darkargs = map[string]string{}
	c.darkargs["units"] = "si"

	clouds := dynamic.Dynamic{Item: c.viper.Get("cloud")}
	for _, e := range clouds.ArrayIter() {
		o := cloud{}
		o.name = e.Get("name").AsString()
		o.description = e.Get("description").AsString()
		o.min = e.Get("min").AsFloat64()
		o.max = e.Get("max").AsFloat64()
		c.clouds = append(c.clouds, o)
	}

	rains := dynamic.Dynamic{Item: c.viper.Get("rain")}
	for _, e := range rains.ArrayIter() {
		o := rain{}
		o.name = e.Get("name").AsString()
		o.unit = e.Get("unit").AsString()
		o.intensity_min = e.Get("intensity_min").AsFloat64()
		o.intensity_max = e.Get("intensity_max").AsFloat64()
		c.rains = append(c.rains, o)
	}

	winds := dynamic.Dynamic{Item: c.viper.Get("wind")}
	for _, e := range winds.ArrayIter() {
		o := wind{}
		o.unit = e.Get("unit").AsString()
		o.speed = e.Get("speed").AsFloat64()
		o.description = []string{}
		descr := e.Get("description").ArrayIter()
		for _, e := range descr {
			o.description = append(o.description, e.AsString())
		}
		c.winds = append(c.winds, o)
	}

	temperatures := dynamic.Dynamic{Item: c.viper.Get("temperature")}
	for _, e := range temperatures.ArrayIter() {
		o := temperature{}
		o.unit = e.Get("unit").AsString()
		o.min = e.Get("min").AsFloat64()
		o.max = e.Get("max").AsFloat64()
		o.description = e.Get("description").AsString()
		c.temperatures = append(c.temperatures, o)
	}

	return c, nil
}

func (c *Client) getCloudsDescription(clouds float64) string {
	for _, cloud := range c.clouds {
		if clouds >= cloud.min && clouds <= cloud.max {
			return cloud.name
		}
	}
	return ""
}

func (c *Client) getRainDescription(rain float64) string {
	for _, r := range c.rains {
		if r.intensity_min <= rain && rain <= r.intensity_max {
			return r.name
		}
	}
	return ""
}

func (c *Client) getTemperatureDescription(temperature float64) string {
	for _, t := range c.temperatures {
		if t.min <= temperature && temperature < t.max {
			return t.description
		}
	}
	return ""
}

func (c *Client) getWindDescription(wind float64) string {
	for _, w := range c.winds {
		if wind < w.speed {
			if len(w.description) > 0 {
				return w.description[0]
			}
			break
		}
	}
	return ""
}

func (c *Client) updateHourly(from time.Time, until time.Time, states *state.Instance, hourly *darksky.DataBlock) {
	states.RemoveAnyStartingWith("weather.hour")

	for _, dp := range hourly.Data {
		hfrom := time.Unix(dp.Time.Unix(), 0)
		huntil := hoursLater(hfrom, 1.0)
		hour := hfrom.Hour()
		if timeRangeInGlobalRange(from, until, hfrom, huntil) {
			states.SetTimeState("weather."+fmt.Sprintf("hour[%d]:from", hour), hfrom)
			states.SetTimeState("weather."+fmt.Sprintf("hour[%d]:until", hour), huntil)

			states.SetFloatState("weather."+fmt.Sprintf("hour[%d]:rain", hour), dp.PrecipProbability)
			states.SetFloatState("weather."+fmt.Sprintf("hour[%d]:clouds", hour), dp.CloudCover)
			states.SetFloatState("weather."+fmt.Sprintf("hour[%d]:temperature", hour), dp.ApparentTemperature)
			states.SetFloatState("weather."+fmt.Sprintf("hour[%d]:wind", hour), dp.WindSpeed)

			states.SetStringState("weather."+fmt.Sprintf("hour[%d]:rain", hour), c.getRainDescription(dp.PrecipProbability))
			states.SetStringState("weather."+fmt.Sprintf("hour[%d]:clouds", hour), c.getCloudsDescription(dp.CloudCover))
			states.SetStringState("weather."+fmt.Sprintf("hour[%d]:temperature", hour), c.getTemperatureDescription(dp.ApparentTemperature))
			states.SetStringState("weather."+fmt.Sprintf("hour[%d]:wind", hour), c.getWindDescription(dp.WindSpeed))
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

func chanceOfRain(from time.Time, until time.Time, states *state.Instance, hourly *darksky.DataBlock) (chanceOfRain string) {

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

func (c *Client) Process(states *state.Instance) time.Duration {
	now := states.GetTimeState("time.now", time.Now())

	// Weather update every 5 minutes
	if now.Unix() >= c.update.Unix() {
		c.update = time.Unix(now.Unix()+5*60, 0)
		//fmt.Println("WEATHER: UPDATE")

		lat := states.GetFloatState("geo.latitude", c.latitude)
		lng := states.GetFloatState("geo.longitude", c.longitude)
		forecast, err := c.darksky.GetForecast(fmt.Sprint(lat), fmt.Sprint(lng), c.darkargs)
		if err == nil {

			from := now
			until := hoursLater(from, 3.0)

			states.SetTimeState("weather.currently:from", from)
			states.SetTimeState("weather.currently:until", until)
			states.SetStringState("weather.currently:rain", chanceOfRain(from, until, states, forecast.Hourly))
			states.SetFloatState("weather.currently:rain", forecast.Currently.PrecipProbability)
			states.SetFloatState("weather.currently:clouds", forecast.Currently.CloudCover)
			states.SetFloatState("weather.currently:temperature", forecast.Currently.ApparentTemperature)

			c.updateHourly(atHour(now, 6, 0), atHour(now, 20, 0), states, forecast.Hourly)
		}
	}

	wait := time.Duration(c.update.Unix()-time.Now().Unix()) * time.Second
	if wait < 0 {
		wait = 0
	}
	return wait
}
