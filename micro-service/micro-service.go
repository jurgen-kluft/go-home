package microservice

import (
	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/nats"
	"strings"
	"time"
)

// Delegate is a handler that the user can register on a certain received topic
type Delegate func(m *Service, topic string, message []byte) bool

// Service is a convenience setup to implement a micro-service
type Service struct {
	Name              string
	UpdateIntervalSec int
	UpdateTimeStamp   time.Time
	Logger            *logpkg.Logger
	PubsubRegister    []string
	PubsubSubscribe   []string
	Pubsub            *pubsub.Context
	Handlers          map[string]Delegate
	CatchHandler      Delegate
}

func New(name string) *Service {
	service := &Service{}

	service.Name = name
	service.Logger = logpkg.New(name)
	service.Logger.AddEntry("pubsub")
	service.Logger.AddEntry(name)

	service.PubsubRegister = make([]string, 0, 10)
	service.PubsubSubscribe = make([]string, 0, 10)
	service.Handlers = make(map[string]Delegate)

	return service
}

func (m *Service) RegisterAndSubscribe(register []string, subscribe []string) {
	for _, r := range register {
		m.PubsubRegister = append(m.PubsubRegister, r)
	}
	for _, r := range subscribe {
		m.PubsubSubscribe = append(m.PubsubSubscribe, r)
	}
}

func (m *Service) RegisterHandler(topic string, delegate Delegate) {
	m.Handlers[topic] = delegate
	natstopic := strings.Replace(topic, "/", ".", -1)
	natstopic = strings.TrimSuffix(natstopic, ".")
	m.Handlers[natstopic] = delegate
}

func (m *Service) Loop() {
	quit := false
	for !quit {
		m.Pubsub = pubsub.New(config.PubSubCfg)
		err := m.Pubsub.Connect(m.Name, m.PubsubRegister, m.PubsubSubscribe)
		if err == nil {
			m.Logger.LogInfo("pubsub", "connected")
			m.Pubsub.Publish("config/request/", m.Name)

			connected := true
			for connected {
				select {
				case msg := <-m.Pubsub.InMsgs:
					topic := m.Pubsub.Topic(msg)
					delegate, exists := m.Handlers[topic]
					if exists {
						if !delegate(m, topic, m.Pubsub.Payload(msg)) {
							connected = false
							quit = true
						}
					} else {
						delegate, exists := m.Handlers["*"]
						if exists {
							if !delegate(m, topic, m.Pubsub.Payload(msg)) {
								connected = false
								quit = true
							}
						}

						if topic == "client/disconnected/" || topic == "client.disconnected" {
							m.Logger.LogInfo("pubsub", "disconnected")
							connected = false
						}
					}
				}
			}
			m.Pubsub.Close()
		}

		if err != nil {
			m.Logger.LogError(m.Name, err.Error())
		}

		m.Logger.LogInfo("pubsub", "Waiting 5 seconds before re-connecting..")
		time.Sleep(5 * time.Second)
	}
}
