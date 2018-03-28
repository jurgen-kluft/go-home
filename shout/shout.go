package shout

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nanopack/mist/clients"
	"github.com/nlopes/slack"
)

// Instance is our instant-messenger instance (currently Slack)
type Instance struct {
	config *Config
	slack  *slack.Client
}

type Config struct {
	Key string `json:"key"`
}

// New creates a new instance of Slack
func New(jsonstr string) (*Instance, error) {
	shout := &Instance{}
	config := &Config{}
	err := json.Unmarshal([]byte(jsonstr), config)
	if err == nil {
		shout.config = config
		shout.slack = slack.New(config.Key)
	}
	return shout, err
}

type Msg struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Pretext  string `json:"pretext"`
	Prebody  string `json:"prebody"`
}

// PostMessage posts a message to a channel
func (s *Instance) postMessage(jsonmsg string) {
	var m Msg
	err := json.Unmarshal([]byte(jsonmsg), &m)
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

func tagsContains(tag string, tags []string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func main() {

	var shout *Instance

	for {
		client, err := clients.New("127.0.0.1:1445", "authtoken.wicked")
		if err != nil {
			fmt.Println(err)
			continue
		}

		client.Ping()
		client.Subscribe([]string{"shout"})
		client.Publish([]string{"request", "config"}, "shout")

		for {
			select {
			case msg := <-client.Messages():
				if tagsContains("config", msg.Tags) {
					shout, err = New(msg.Data)
				} else {
					// Is this a message to send over slack ?
					shout.postMessage(msg.Data)
				}
				break
			case <-time.After(time.Second * 60):

				break
			}
		}

		// Disconnect from Mist
	}
}
