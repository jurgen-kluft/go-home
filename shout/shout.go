package shout

import (
	"fmt"
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
		shout.client = slack.New(config.Key)
	}
	return shout, err
}

// postMessage posts a message to a channel
func (s *Instance) postMessage(jsonmsg string) {
	m, err := config.ShoutMsgFromJSON(jsonmsg)
	if err == nil {
		params := slack.PostMessageParameters{}
		params.Username = m.Username
		attachment := slack.Attachment{
			Pretext: m.Pretext,
			Text:    m.Prebody,
		}
		params.Attachments = []slack.Attachment{attachment}
		_, timestamp, err := s.client.PostMessage(m.Channel, m.Message, params)
		if err != nil {
			fmt.Printf("Error '%s' at %s\n", err, timestamp)
		}
		//fmt.Printf("Message successfully sent to channel %s at %s\n", channel, timestamp)
	}
}

func main() {

	var shout *Instance

	logger := logpkg.New("shout")
	logger.AddEntry("emitter")
	logger.AddEntry("shout")

	for {
		connected := true
		for connected {
			client := pubsub.New(config.EmitterSecrets["host"])
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
						} else if topic == "client/disconnected/" {
							logger.LogInfo("emitter", "disconnected")
							connected = false
						} else if topic == "shout/message/" {
							// Is this a message to send over slack ?
							if shout != nil {
								logger.LogInfo("shout", "message")
								shout.postMessage(string(msg.Payload()))
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
