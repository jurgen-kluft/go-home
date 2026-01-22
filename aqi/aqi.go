package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
)

type instance struct {
	name   string
	config *config.AqiConfig
	update time.Time
}

func construct() (c *instance) {
	c = &instance{}
	c.name = "aqi"
	c.update = time.Now()
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
	peers := []string{"config/request"}

	c := construct()
	m, err := microservice.New("aqi", "sensor/aqi", time.Minute*10)
	if err != nil {
		panic(err)
	}
	m.ConnectTo(peers)

	m.RegisterHandler("config/request", func(m *microservice.Service, msg *microservice.Message) bool {
		configAqi, err := config.AqiConfigFromJSON(msg.Payload)
		if err == nil {
			m.Logger.LogInfo(m.Name, "received configuration")
			c.config = configAqi
		} else {
			m.Logger.LogError(m.Name, "received bad configuration, "+err.Error())
		}
		return true
	})

	m.RegisterHandler("tick", func(m *microservice.Service, msg *microservice.Message) bool {
		if c != nil && c.config != nil {
			m.Logger.LogInfo(m.Name, "polling Aqi")
			stateAsJson, err := c.Poll()
			if err == nil {
				m.Logger.LogInfo(m.Name, "publish Aqi")
				msg, err := m.NewTextMessage(string(stateAsJson))
				if err == nil {
					_ = m.Broadcast(msg)
				}
			} else {
				m.Logger.LogError(m.Name, err.Error())
			}
		} else if c != nil && c.config == nil {
			// Try and request our configuration
			msg, err := m.NewTextMessage(m.Name)
			if err == nil {
				_ = m.SendTo("config/request", msg)
			}
		}
		return true
	})

	m.Loop()
}
