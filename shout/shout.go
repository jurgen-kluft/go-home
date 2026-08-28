package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
)

// Shout microservice instance
// This microservice shouts messages to specific Apple Home Pod devices using
// Play audio on a HomePod by switching the `audio output` on the Mac Mini, we can do this from the terminal.
// This means we can write a script that can convert text to an mp4 and then play it on one or more HomePods.
// This is useful for automations that need to announce something, like when the door bell is pressed or when
// the front door is opened or when the wash machine or dryer is done.
// We can do this by writing to an announcement.queue file, and we have a process that tails the file and plays
// the announcements one at a time on the designated HomePod(s).
//
// Example announcement on HomePod1 and HomePod2:
// ```
// say "Hello, this is a test announcement" --output announcement.mp4
// select audio output as HomePod1
// afplay announcement.mp4
// select audio output as HomePod2
// afplay announcement.mp4
// ```

type instance struct {
	name    string
	config  *config.ShoutConfig
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

	}
	return err
}

// postMessage posts a message to a channel
func (s *instance) postMessage(jsondata []byte) {
}

func main() {

	thisConfigFilepath := "config/shout.config.json"
	servicesConfigFilepath := "config/shout.config.json"

	m, err := microservice.New("shout", thisConfigFilepath, servicesConfigFilepath, time.Second*15)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := new()
	c.service = m

	m.RegisterHandler("config/request", func(m *microservice.Service, msg *microservice.Message) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		err := c.initialize(msg.Payload)
		if err != nil {
			m.Logger.LogError(m.Name, err.Error())
		} else {
			c.postMessage([]byte("service connected"))
		}
		return true
	})

	m.RegisterHandler("shout/message", func(m *microservice.Service, msg *microservice.Message) bool {

		return true
	})

	tickCount := 0
	m.RegisterHandler("tick", func(m *microservice.Service, msg *microservice.Message) bool {
		if (tickCount % 30) == 0 {
			if c.config == nil {
				msg, err := m.NewTextMessage(m.Name)
				if err == nil {
					m.SendTo("config/request", msg)
				}
			}
		}
		tickCount++
		return true
	})

	m.Start()
	m.Loop()
}
