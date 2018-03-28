package aqi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/nanopack/mist/clients"
)

type instance struct {
	config *Config
	update time.Time
	period time.Duration
}

type SensorState struct {
	Domain  string    `json:"domain"`
	Product string    `json:"product"`
	Name    string    `json:"name"`
	Type    string    `json:"type"`
	Value   string    `json:"value"`
	Time    time.Time `json:"time"`
}

func (c *instance) readConfig(jsonstr string) (*Config, error) {
	jsonBytes := []byte(jsonstr)
	obj, err := unmarshalConfig(jsonBytes)
	return obj, err
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

func (c *instance) getAiqTagAndDescr(aiq float64) (level AqiLevel) {
	for _, l := range c.config.Levels {
		if aiq < l.LessThan {
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
		sensor := SensorState{Domain: "sensor", Product: "weather", Name: "aqi", Type: "float", Value: fmt.Sprintf("%f", aqi), Time: time.Now()}
		jsonbytes, err := json.Marshal(sensor)
		if err == nil {
			aqiStateJSON = string(jsonbytes)
		}
	}
	return aqiStateJSON, err
}

func tagsContains(tag string, tags []string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func main() {

	aqi := construct()

	for {
		client, err := clients.New("127.0.0.1:1445", "authtoken.wicked")
		if err != nil {
			fmt.Println(err)
			continue
		}

		client.Ping()
		client.Subscribe([]string{"aqi"})
		client.Publish([]string{"request", "config"}, "aqi")

		for {
			select {
			case msg := <-client.Messages():
				if tagsContains("config", msg.Tags) {
					aqi.config, err = aqi.readConfig(msg.Data)
				}
				break
			case <-time.After(time.Second * 10):
				if aqi != nil && aqi.config != nil {
					if aqi.shouldPoll(time.Now(), false) {
						jsonstate, err := aqi.Poll()
						if err == nil {
							client.Publish([]string{"sensor", "weather", "aqi"}, jsonstate)
						}
						aqi.computeNextPoll(time.Now(), err)
					}
				}
				break
			}
		}

		// Disconnect from Mist
	}
}
