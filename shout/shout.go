package main

import (
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/nlopes/slack"
)

// Instance is our instant-messenger instance (currently Slack)
type instance struct {
	name   string
	config *config.ShoutConfig
	client *slack.Client
}

func new() *instance {
	s := &instance{}
	s.name = "shout"
	return s
}

// New creates a new instance of Slack
func (s *instance) initialize(jsondata []byte) error {
	s.name = "shout"
	config, err := config.ShoutConfigFromJSON(jsondata)
	if err == nil {
		s.config = config
		s.client = slack.New(config.Key.String)
	}
	return err
}

// postMessage posts a message to a channel
func (s *instance) postMessage(jsondata []byte) (err error) {
	m, err := config.ShoutMsgFromJSON(jsonmsg)
	if err == nil {
		_, _, err = s.client.PostMessage(m.Channel, slack.MsgOptionText("Some text", false), slack.MsgOptionUsername("g0-h0m3"), slack.MsgOptionAsUser(true))
	}
	return err
}

func main() {

	c := new()

	logger := logpkg.New(c.name)
	logger.AddEntry("emitter")
	logger.AddEntry(c.name)

	for {
		connected := true
		for connected {
			client := pubsub.New(config.PubSubCfg)
			register := []string{"config/shout/", "shout/message/"}
			subscribe := []string{"config/shout/", "shout/message/", "config/request/"}
			err := client.Connect(c.name, register, subscribe)
			if err == nil {
				logger.LogInfo("emitter", "connected")

				for connected {
					select {
					case msg := <-client.InMsgs:
						topic := msg.Topic()
						if topic == "config/shout/" {
							logger.LogInfo(c.name, "received configuration")
							err = c.initialize(msg.Payload())
							if err != nil {
								logger.LogError(c.name, err.Error())
							}
						} else if topic == "client/disconnected/" {
							logger.LogInfo("emitter", "disconnected")
							connected = false
						} else if topic == "shout/message/" {
							// Is this a message to send over slack ?
							if c.client != nil {
								logger.LogInfo(c.name, "message")
								err = c.postMessage(msg.Payload())
								if err != nil {
									logger.LogError(c.name, err.Error())
								}
							}
						}
					case <-time.After(time.Minute * 1): // Try and request our configuration
						if c.config == nil {
							client.Publish("config/request/", "shout")
						}
					}
				}
			} else {
				connected = false
			}

			if err != nil {
				logger.LogError(c.name, err.Error())
			}
		}

		time.Sleep(5 * time.Second)
	}
}
