package main

import (
	"fmt"
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
func (s *instance) postMessage(jsondata []byte) {
	if s.client != nil {
		channelID, timestamp, err := s.client.PostMessage(s.config.Channel, slack.MsgOptionText(string(jsondata), false), slack.MsgOptionUsername("g0-h0m3"), slack.MsgOptionAsUser(true))
		if err == nil {
			s.service.Logger.LogInfo(s.name, fmt.Sprintf("message '%s' send (%s, %s)", string(jsondata), channelID, timestamp))
		} else {
			s.service.Logger.LogError(s.name, fmt.Sprintf("message '%s' not send (%s, %s)", string(jsondata), s.config.Channel, timestamp))
			s.service.Logger.LogError(s.name, err.Error())
		}
	} else {
		s.service.Logger.LogError(s.name, "service not connected")
	}
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
		} else {
			c.postMessage([]byte("service connected"))
		}
		return true
	})

	m.RegisterHandler("shout/message/", func(m *microservice.Service, topic string, msg []byte) bool {
		// Is this a message to send over slack ?
		if c.client != nil {
			m.Logger.LogInfo(m.Name, "message")
			c.postMessage(msg)
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if (tickCount % 30) == 0 {
			if c.config == nil {
				m.Pubsub.Publish("config/request/", []byte(m.Name))
			}
		}
		tickCount++
		return true
	})

	m.Loop()
}
