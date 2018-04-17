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
	config  *config.AqiConfig
	update  time.Time
	metrics *metrics.Metrics
}

func construct() (c *instance) {
	c = &instance{}
	c.update = time.Now()
	c.metrics, _ = metrics.New()

	c.metrics.Register("aqi", map[string]string{"aqi": "quality"}, map[string]interface{}{"pm2.5": 50.0})
	return c
}

func (c *instance) getResponse() (AQI float64, err error) {
	url := c.config.URL
	url = strings.Replace(url, "${CITY}", c.config.City, 1)
	url = strings.Replace(url, "${TOKEN}", c.config.Token.String, 1)
	if strings.HasPrefix(url, "http") {
		var resp *http.Response
		resp, err = http.Get(url)
		if err != nil {
			AQI = 80.0
		} else {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			var caqi CaqiResponse
			caqi, err = unmarshalCaqiResponse(body)
			AQI = float64(caqi.Data.Aqi)
			if err != nil {
				fmt.Print(string(body))
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
		c.metrics.Begin("aqi")
		c.metrics.Set("aqi", "pm2.5", aqi)
		c.metrics.Send("aqi")

		// MQTT: As a sensor
		sensor := config.NewSensorState("sensor.weather.aqi")
		sensor.AddFloatAttr("aqi", aqi)
		level := c.getAiqTagAndDescr(aqi)
		sensor.AddStringAttr("name", level.Tag)
		sensor.AddStringAttr("caution", level.Caution)
		sensor.AddStringAttr("implications", level.Implications)
		aqiStateJSON, err = sensor.ToJSON()
	}
	return aqiStateJSON, err
}

func main() {

	aqi := construct()

	logger := logpkg.New("aqi")
	logger.AddEntry("emitter")
	logger.AddEntry("aqi")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/aqi/", "state/sensor/aqi/"}
		subscribe := []string{"config/aqi/"}
		err := client.Connect("aqi", register, subscribe)
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
							logger.LogInfo("aqi", "received configuration")
							aqi.config = config
							pollCount = 0
						} else {
							logger.LogError("aqi", "received bad configuration, "+err.Error())
						}
					} else if topic == "client/disconnected/" {
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}

				case <-time.After(time.Second * 10):
					if aqi != nil && aqi.config != nil {
						if aqi.shouldPoll(time.Now(), pollCount == 0) {
							logger.LogInfo("aqi", "polling Aqi")
							jsonstate, err := aqi.Poll()
							if err == nil {
								logger.LogInfo("aqi", "publish Aqi")
								fmt.Println(jsonstate)
								client.PublishTTL("state/sensor/aqi/", jsonstate, 5*60)
							} else {
								logger.LogError("aqi", err.Error())
							}
							pollCount++
							aqi.computeNextPoll(time.Now(), err)
						}
					}
				}
			}
		}
		if err != nil {
			logger.LogError("aqi", err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}
