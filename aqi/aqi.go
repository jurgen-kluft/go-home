package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/metrics"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
)

type instance struct {
	name    string
	config  *config.AqiConfig
	update  time.Time
	metrics *metrics.Metrics
}

func construct() (c *instance) {
	c = &instance{}
	c.name = "aqi"
	c.update = time.Now()
	c.metrics, _ = metrics.New()

	c.metrics.Register(c.name, map[string]string{c.name: "quality"}, map[string]interface{}{"pm2.5": 50.0})
	return c
}

func (c *instance) getResponse() (AQI float64, err error) {
	url := c.config.URL
	url = strings.Replace(url, "${CITY}", c.config.City, 1)
	url = strings.Replace(url, "${TOKEN}", c.config.Token.String, 1)
	if strings.HasPrefix(url, "http") {
		var resp *http.Response
		resp, err = http.Get(url)
		AQI = 80.0
		if err == nil {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				var caqi CaqiResponse
				caqi, err = unmarshalCaqiResponse(body)
				if err == nil {
					AQI = float64(caqi.Data.Aqi)
				}
			}
			resp.Body.Close()
		}
	} else if strings.HasPrefix(url, "print") {
		fmt.Printf("HTTP Get, '%s'\n", url)
	}
	return
}

func (c *instance) getAiqTagAndDescr(aqi float64) (level config.AqiLevel) {
	for _, l := range c.config.Levels {
		if aqi < l.LessThan {
			level = l
			return
		}
	}
	level = c.config.Levels[1]
	return
}

func (c *instance) shouldPoll(now time.Time, force bool) bool {
	if force || (now.Unix() >= c.update.Unix()) {
		return true
	}
	return false
}

func (c *instance) computeNextPoll(now time.Time, err error) {
	if err != nil {
		c.update = now.Add(time.Second * time.Duration(c.config.Interval))
	} else {
		c.update = now.Add(time.Duration(c.config.Interval) * time.Second)
	}
}

// Poll will get AQI information and returns a JSON string
func (c *instance) Poll() (aqiStateJSON []byte, err error) {
	aqiStateJSON = []byte{}
	aqi, err := c.getResponse()
	if err == nil {

		// Metrics
		c.metrics.Begin(c.name)
		c.metrics.Set(c.name, "pm2.5", aqi)
		c.metrics.Send(c.name)

		// MQTT: As a sensor
		sensor := config.NewSensorState("sensor.weather.aqi", "airquality")
		sensor.AddFloatAttr(c.name, aqi)
		level := c.getAiqTagAndDescr(aqi)
		sensor.AddStringAttr("name", level.Tag)
		sensor.AddStringAttr("caution", level.Caution)
		sensor.AddStringAttr("implications", level.Implications)
		aqiStateJSON, err = sensor.ToJSON()
	}
	return aqiStateJSON, err
}

func main() {
	register := []string{"config/aqi/", "config/request/", "state/sensor/aqi/"}
	subscribe := []string{"config/aqi/"}

	c := construct()
	m := microservice.New("aqi", time.Minute*10)
	m.RegisterAndSubscribe(register, subscribe)

	m.RegisterHandler("config/aqi/", func(m *microservice.Service, topic string, msg []byte) bool {
		configAqi, err := config.AqiConfigFromJSON(msg)
		if err == nil {
			m.Logger.LogInfo(m.Name, "received configuration")
			c.config = configAqi
		} else {
			m.Logger.LogError(m.Name, "received bad configuration, "+err.Error())
		}
		return true
	})

	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if c != nil && c.config != nil {
			m.Logger.LogInfo(m.Name, "polling Aqi")
			stateAsJson, err := c.Poll()
			if err == nil {
				m.Logger.LogInfo(m.Name, "publish Aqi")
				_ = m.Pubsub.Publish("state/sensor/aqi/", stateAsJson)
			} else {
				m.Logger.LogError(m.Name, err.Error())
			}
		} else if c != nil && c.config == nil {
			// Try and request our configuration
			_ = m.Pubsub.PublishStr("config/request/", "aqi")
		}
		return true
	})

	m.Loop()
}
