package main

import (
	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/micro-service"
	"github.com/nlopes/slack"
)

// Instance is our instant-messenger instance (currently Slack)
type instance struct {
	name    string
	config  *config.ShoutConfig
	client  *slack.Client
	service *microservice.Service
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
	m, err := config.ShoutMsgFromJSON(jsondata)
	if err == nil {
		_, _, err = s.client.PostMessage(m.Channel, slack.MsgOptionText("Some text", false), slack.MsgOptionUsername("g0-h0m3"), slack.MsgOptionAsUser(true))
	}
	return err
}

func main() {
	register := []string{"config/shout/", "shout/message/"}
	subscribe := []string{"config/shout/", "shout/message/", "config/request/"}

	m := microservice.New("shout")
	m.RegisterAndSubscribe(register, subscribe)

	c := new()
	c.service = m

	m.RegisterHandler("config/shout/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		err := c.initialize(msg)
		if err != nil {
			m.Logger.LogError(m.Name, err.Error())
		}
		return true
	})

	m.RegisterHandler("shout/message/", func(m *microservice.Service, topic string, msg []byte) bool {
		// Is this a message to send over slack ?
		if c.client != nil {
			m.Logger.LogInfo(m.Name, "message")
			err := c.postMessage(msg)
			if err != nil {
				m.Logger.LogError(c.name, err.Error())
			}
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if (tickCount % 30) == 0 {
			if c.config == nil {
				m.Pubsub.Publish("config/request/", m.Name)
			}
		}
		tickCount++
		return true
	})

	m.Loop()
}
