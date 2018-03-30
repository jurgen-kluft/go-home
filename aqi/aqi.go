package aqi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/emitter-io/go"
	"github.com/jurgen-kluft/go-home/config"
)

type instance struct {
	config *config.AqiConfig
	update time.Time
	period time.Duration
}

func (c *instance) readConfig(jsonstr string) (*config.AqiConfig, error) {
	obj, err := config.AqiConfigFromJSON(jsonstr)
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
		sensor := config.SensorState{Domain: "sensor", Product: "weather", Name: "aqi", Type: "float", Value: fmt.Sprintf("%f", aqi), Time: time.Now()}
		jsonbytes, err := json.Marshal(sensor)
		if err == nil {
			aqiStateJSON = string(jsonbytes)
		}
	}
	return aqiStateJSON, err
}

func strContains(tag string, tags []string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

type context struct {
	inmsgs chan emitter.Message
	inpres chan emitter.PresenceEvent
}

type DisconnectMessage struct {
}

func (d *DisconnectMessage) Topic() string {
	return "client/disconnected"
}
func (d *DisconnectMessage) Payload() []byte {
	return []byte{}
}

func main() {

	aqi := construct()

	for {
		connected := true
		for connected {
			// Create the options with default values
			emitterOptions := emitter.NewClientOptions()
			emitterOptions.SetUsername("aqi")

			ctx := context{}
			ctx.inmsgs = make(chan emitter.Message)

			// Set the message handler
			emitterOptions.SetOnMessageHandler(func(client emitter.Emitter, msg emitter.Message) {
				ctx.inmsgs <- msg
			})

			// Set the presence notification handler
			emitterOptions.SetOnPresenceHandler(func(_ emitter.Emitter, p emitter.PresenceEvent) {
				fmt.Printf("Occupancy: %v\n", p.Occupancy)
			})

			// Set the connection lost handler
			emitterOptions.SetOnConnectionLostHandler(func(_ emitter.Emitter, e error) {
				msg := &DisconnectMessage{}
				ctx.inmsgs <- msg
			})

			// Create a new emitter client and connect to the broker
			emitterClient := emitter.NewClient(emitterOptions)
			sToken := emitterClient.Connect()
			if sToken.Wait() && sToken.Error() == nil {

				// Subscribe to the presence demo channel
				emitterClient.Subscribe(config.SecretKey, "aqi/+")

				for connected {
					select {
					case msg := <-ctx.inmsgs:
						topic := msg.Topic()
						if topic == "aqi/config" {
							jsonmsg := string(msg.Payload())
							config, err := aqi.readConfig(jsonmsg)
							if err == nil {
								aqi.config = config
							}
						} else if topic == "client/disconnected" {
							connected = false
						}
						break
					case <-time.After(time.Second * 10):
						if aqi != nil && aqi.config != nil {
							if aqi.shouldPoll(time.Now(), false) {
								jsonstate, err := aqi.Poll()
								if err == nil {
									emitterClient.PublishWithTTL(config.SecretKey, "sensor/weather/aqi", jsonstate, 5*60)
								}
								aqi.computeNextPoll(time.Now(), err)
							}
						}
						break

					}
				}
			} else {
				panic("Error on Client.Connect(): " + sToken.Error().Error())
			}
		}

		// Wait for 10 seconds before retrying
		time.Sleep(10 * time.Second)
	}
}
