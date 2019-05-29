package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/metrics"
	"github.com/jurgen-kluft/go-home/pubsub"
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
func (c *instance) Poll() (aqiStateJSON string, err error) {
	aqiStateJSON = ""
	aqi, err := c.getResponse()
	if err == nil {

		// Metrics
		c.metrics.Begin(c.name)
		c.metrics.Set(c.name, "pm2.5", aqi)
		c.metrics.Send(c.name)

		// MQTT: As a sensor
		sensor := config.NewSensorState("sensor.weather.aqi")
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

	c := construct()

	logger := logpkg.New(c.name)
	logger.AddEntry("emitter")
	logger.AddEntry(c.name)

	for {
		client := pubsub.New(config.EmitterIOCfg)
		register := []string{"config/aqi/", "state/sensor/aqi/"}
		subscribe := []string{"config/aqi/", "config/request/"}
		err := client.Connect(c.name, register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			pollCount := int64(0)
			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/aqi/" {
						jsonmsg := string(msg.Payload())
						config, err := config.AqiConfigFromJSON(jsonmsg)
						if err == nil {
							logger.LogInfo(c.name, "received configuration")
							c.config = config
							pollCount = 0
						} else {
							logger.LogError(c.name, "received bad configuration, "+err.Error())
						}
					} else if topic == "client/disconnected/" {
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}

				case <-time.After(time.Second * 10):
					if c != nil && c.config != nil {
						if c.shouldPoll(time.Now(), pollCount == 0) {
							logger.LogInfo(c.name, "polling Aqi")
							jsonstate, err := c.Poll()
							if err == nil {
								logger.LogInfo(c.name, "publish Aqi")
								fmt.Println(jsonstate)
								client.PublishTTL("state/sensor/aqi/", jsonstate, 5*60)
							} else {
								logger.LogError(c.name, err.Error())
							}
							pollCount++
							c.computeNextPoll(time.Now(), err)
						}
					}

				case <-time.After(time.Minute * 1):
					if c != nil && c.config == nil {
						// Try and request our configuration
						client.Publish("config/request/", "aqi")
					}
				}
			}
		}
		if err != nil {
			logger.LogError(c.name, err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}
