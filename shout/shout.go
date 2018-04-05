package shout

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/nlopes/slack"
)

// Instance is our instant-messenger instance (currently Slack)
type Instance struct {
	config *config.ShoutConfig
	slack  *slack.Client
}

// New creates a new instance of Slack
func New(jsonstr string) (*Instance, error) {
	shout := &Instance{}
	config, err := config.ShoutConfigFromJSON(jsonstr)
	if err == nil {
		shout.config = config
		shout.slack = slack.New(config.Key)
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
		_, timestamp, err := s.slack.PostMessage(m.Channel, m.Message, params)
		if err != nil {
			fmt.Printf("Error '%s' at %s\n", err, timestamp)
		}
		//fmt.Printf("Message successfully sent to channel %s at %s\n", channel, timestamp)
	}
}

func main() {

	var shout *Instance

	for {
		connected := true
		for connected {
			client := pubsub.New()
			err := client.Connect("shout")
			if err == nil {

				// Subscribe to the presence demo channel
				client.Register("config/shout")
				client.Register("shout/message")

				client.Subscribe("config/shout")
				client.Subscribe("shout/message")

				for connected {
					select {
					case msg := <-client.InMsgs:
						topic := msg.Topic()
						if topic == "config/shout" {
							shout, err = New(string(msg.Payload()))
						} else if topic == "client/disconnected" {
							connected = false
						} else if topic == "shout/message" {
							// Is this a message to send over slack ?
							shout.postMessage(string(msg.Payload()))
						}
						break
					case <-time.After(time.Second * 10):

						break
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
