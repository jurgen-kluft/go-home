package microservice

import (
	"strings"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	pubsub "github.com/jurgen-kluft/go-home/mqtt"
)

// Delegate is a handler that the user can register on a certain received topic
type Delegate func(m *Service, topic string, message []byte) bool

type Message struct {
	Topic   string
	Payload []byte
}

// Service is a convenience setup to implement a micro-service
type Service struct {
	Name            string
	Logger          *logpkg.Logger
	PubsubRegister  []string
	PubsubSubscribe []string
	Pubsub          *pubsub.Context
	Handlers        map[string]Delegate
	CatchHandler    Delegate
	ProcessMessages chan *Message
	TickFrequency   time.Duration
}

func New(name string, tickFrequency time.Duration) *Service {
	service := &Service{}

	service.Name = name
	service.Logger = logpkg.New(name)
	service.Logger.AddEntry("pubsub")
	service.Logger.AddEntry(name)

	service.PubsubRegister = make([]string, 0, 10)
	service.PubsubSubscribe = make([]string, 0, 10)
	service.Handlers = make(map[string]Delegate)

	service.ProcessMessages = make(chan *Message, 128)
	service.TickFrequency = tickFrequency
	return service
}

func (m *Service) Register(r string) error {
	if m.Pubsub == nil {
		// Not connected yet, just add it to the list
		m.PubsubRegister = append(m.PubsubRegister, r)
	} else {
		// We are connected, also call Register on pubsub
		m.PubsubRegister = append(m.PubsubRegister, r)
		m.Pubsub.Register(r)
	}
	return nil
}

func (m *Service) Subscribe(r string) error {
	if m.Pubsub == nil {
		// Not connected yet, just add it to the list
		m.PubsubSubscribe = append(m.PubsubSubscribe, r)
	} else {
		// We are connected, also call Subscribe on pubsub
		m.PubsubSubscribe = append(m.PubsubSubscribe, r)
		m.Pubsub.Subscribe(r)
	}
	return nil
}

func (m *Service) RegisterAndSubscribe(register []string, subscribe []string) {
	for _, r := range register {
		m.Register(r)
	}
	for _, r := range subscribe {
		m.Subscribe(r)
	}
}

func (m *Service) RegisterHandler(topic string, delegate Delegate) {
	m.Handlers[topic] = delegate
	natstopic := strings.Replace(topic, "/", ".", -1)
	natstopic = strings.TrimSuffix(natstopic, ".")
	m.Handlers[natstopic] = delegate
}

func matchTopic(etopic string, itopic string) bool {
	ei := 0
	ii := 0
	cs := true
	for ei < len(etopic) && ii < len(itopic) {
		echar := etopic[ei]
		ichar := itopic[ii]
		if (echar == '/' || echar == '.') && (ichar == '/' || ichar == '.') {
			cs = true
			ei++
			ii++
		} else if cs && etopic[ei] == '*' { // chapter start ?
			// consume chapter from itopic
			ei++
			for ii < len(itopic) && (itopic[ii] != '/' && itopic[ii] != '.') {
				ii++
			}
			if ii == len(itopic) {
				return false
			}
			cs = false
		} else if echar != ichar {
			return false
		} else {
			ei++
			ii++
			cs = false
		}
	}
	return ei == len(etopic) && ii == len(itopic)
}

func (m *Service) FindHandler(itopic string) (delegate Delegate, exists bool) {
	for etopic, edelegate := range m.Handlers {
		if matchTopic(etopic, itopic) {
			return edelegate, true
		}
	}
	return nil, false
}

func (m *Service) Loop() {
	quit := false
	for !quit {
		m.Pubsub = pubsub.New(config.PubSubCfg, m.TickFrequency)
		err := m.Pubsub.Connect(m.Name, m.PubsubRegister, m.PubsubSubscribe)
		if err == nil {
			m.Logger.LogInfo("pubsub", "connected")

			connected := true
			for connected {
				select {
				case msg := <-m.ProcessMessages:
					topic := msg.Topic
					delegate, exists := m.FindHandler(topic)
					if exists {
						if !delegate(m, topic, msg.Payload) {
							connected = false
							quit = true
						}
					}

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

						if topic == "client/disconnected" {
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

		if !quit {
			m.Logger.LogInfo("pubsub", "Waiting 5 seconds before re-connecting..")
			time.Sleep(5 * time.Second)
		} else {
			m.Logger.LogInfo("pubsub", "End.")
		}
	}
}
