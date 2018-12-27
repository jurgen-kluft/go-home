package main

import (
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/nlopes/slack"
)

// Instance is our instant-messenger instance (currently Slack)
type Instance struct {
	config *config.ShoutConfig
	client *slack.Client
}

// New creates a new instance of Slack
func New(jsonstr string) (*Instance, error) {
	shout := &Instance{}
	config, err := config.ShoutConfigFromJSON(jsonstr)
	if err == nil {
		shout.config = config
		shout.client = slack.New(config.Key.String)
	}
	return shout, err
}

// postMessage posts a message to a channel
func (s *Instance) postMessage(jsonmsg string) (err error) {
	m, err := config.ShoutMsgFromJSON(jsonmsg)
	if err == nil {
		_, _, err = s.client.PostMessage(m.Channel, slack.MsgOptionText("Some text", false), slack.MsgOptionUsername("g0-h0m3"), slack.MsgOptionAsUser(true))
	}
	return err
}

func main() {

	var shout *Instance

	logger := logpkg.New("shout")
	logger.AddEntry("emitter")
	logger.AddEntry("shout")

	for {
		connected := true
		for connected {
			client := pubsub.New(config.EmitterIOCfg)
			register := []string{"config/shout/", "shout/message/"}
			subscribe := []string{"config/shout/", "shout/message/"}
			err := client.Connect("shout", register, subscribe)
			if err == nil {
				logger.LogInfo("emitter", "connected")

				for connected {
					select {
					case msg := <-client.InMsgs:
						topic := msg.Topic()
						if topic == "config/shout/" {
							logger.LogInfo("shout", "received configuration")
							shout, err = New(string(msg.Payload()))
							if err != nil {
								logger.LogError("shout", err.Error())
							}
						} else if topic == "client/disconnected/" {
							logger.LogInfo("emitter", "disconnected")
							connected = false
						} else if topic == "shout/message/" {
							// Is this a message to send over slack ?
							if shout != nil {
								logger.LogInfo("shout", "message")
								err = shout.postMessage(string(msg.Payload()))
								if err != nil {
									logger.LogError("shout", err.Error())
								}
							}
						}
						break
					case <-time.After(time.Second * 10):

						break
					}
				}
			} else {
				connected = false
			}

			if err != nil {
				logger.LogError("shout", err.Error())
			}
		}

		time.Sleep(5 * time.Second)
	}
}
