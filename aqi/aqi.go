package aqi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/pubsub"
)

type instance struct {
	config *config.AqiConfig
	update time.Time
	period time.Duration
}

func construct() (c *instance) {
	c = &instance{}
	c.update = time.Now()
	c.period = time.Minute * 15
	return c
}

func (c *instance) getResponse() (AQI float64, err error) {
	url := c.config.URL
	url = strings.Replace(url, "${CITY}", c.config.City, 1)
	url = strings.Replace(url, "${TOKEN}", c.config.Token, 1)
	if strings.HasPrefix(url, "http") {
		var resp *http.Response
		fmt.Printf("HTTP Get, '%s'\n", url)
		resp, err = http.Get(url)
		if err != nil {
			AQI = 99.0
			resp.Body.Close()
		} else {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			var caqi CaqiResponse
			caqi, err = unmarshalCaqiResponse(body)
			AQI = float64(caqi.Data.Aqi)
			if err != nil {
				fmt.Print(string(body))
			}
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
		c.update = now.Add(time.Minute * 1)
	} else {
		c.update = now.Add(c.period)
	}
}

// Poll will get AQI information and returns a JSON string
func (c *instance) Poll() (aqiStateJSON string, err error) {
	aqiStateJSON = ""
	aqi, err := c.getResponse()
	if err == nil {
		aqiStateJSON, err = config.FloatSensorAsJSON("sensor.weather.aqi", "aqi", aqi)
	}
	return aqiStateJSON, err
}

func main() {

	aqi := construct()

	for {
		connected := true
		for connected {
			client := pubsub.New()
			err := client.Connect("aqi")
			if err == nil {

				client.Subscribe("config/aqi")

				for connected {
					select {
					case msg := <-client.InMsgs:
						topic := msg.Topic()
						if topic == "config/aqi" {
							jsonmsg := string(msg.Payload())
							config, err := config.AqiConfigFromJSON(jsonmsg)
							if err == nil {
								aqi.config = config
							}
						} else if topic == "client/disconnected" {
							connected = false
						}

					case <-time.After(time.Second * 10):
						if aqi != nil && aqi.config != nil {
							if aqi.shouldPoll(time.Now(), false) {
								jsonstate, err := aqi.Poll()
								if err == nil {
									client.PublishTTL("state/sensor/aqi", jsonstate, 5*60)
								}
								aqi.computeNextPoll(time.Now(), err)
							}
						}

					}
				}
			} else {
				panic("Error on Client.Connect(): " + err.Error())
			}
		}

		// Wait for 10 seconds before retrying
		time.Sleep(10 * time.Second)
	}
}
