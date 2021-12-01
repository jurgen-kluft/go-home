package main

import (
	"fmt"

	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
	"github.com/slack-go/slack"
)

// Instance is our instant-messenger instance (currently Slack)
type instance struct {
	name   string
	config *config.ShoutConfig
	client *slack.Client
	//socket_client *socketmode.Client
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
		s.client = slack.New(s.config.UserToken.String, slack.OptionDebug(true), slack.OptionAppLevelToken(s.config.AppToken.String))

		/// The below does work, we can receive messages :-)

		///		// go-slack comes with a SocketMode package that we need to use that accepts a Slack client and outputs a Socket mode client instead
		///		s.socket_client = socketmode.New(
		///			s.client,
		///			socketmode.OptionDebug(true),
		///			// Option to set a custom logger
		///			socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
		///		)
		///
		///		// Create a context that can be used to cancel goroutine
		///		ctx, cancel := context.WithCancel(context.Background())
		///		// Make this cancel called properly in a real program , graceful shutdown etc
		///		defer cancel()
		///
		///		go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		///			// Create a for loop that selects either the context cancellation or the events incomming
		///			for {
		///				select {
		///				// inscase context cancel is called exit the goroutine
		///				case <-ctx.Done():
		///					log.Println("Shutting down socketmode listener")
		///					return
		///				case event := <-socketClient.Events:
		///					// We have a new Events, let's type switch the event
		///					// Add more use cases here if you want to listen to other events.
		///					switch event.Type {
		///					// handle EventAPI events
		///					case slackevents.Message:
		///
		///						// The application has been mentioned since this Event is a Mention event
		///						log.Println(event)
		///
		///					case socketmode.EventTypeEventsAPI:
		///						// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
		///						eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
		///						if !ok {
		///							log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
		///							continue
		///						}
		///						// We need to send an Acknowledge to the slack server
		///						socketClient.Ack(*event.Request)
		///						// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
		///						log.Println(eventsAPIEvent)
		///					}
		///
		///				}
		///			}
		///		}(ctx, s.client, s.socket_client)
		///
		///		s.socket_client.Run()
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
	register := []string{"config/request/"}
	subscribe := []string{"config/shout/", "shout/message/"}

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
