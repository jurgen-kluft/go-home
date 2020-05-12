package pubsub

import (
	"fmt"
	"sync/atomic"
	"time"

	emitter "github.com/emitter-io/go"
)

type AtomBool int32

func (b *AtomBool) Set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32((*int32)(b), int32(i))
}

func (b *AtomBool) IsTrue() bool {
	if atomic.LoadInt32((*int32)(b)) != 0 {
		return true
	}
	return false
}

type Context struct {
	EmitterCfg  map[string]string
	ChannelKeys map[string]string
	InMsgs      chan emitter.Message
	Client      emitter.Emitter
	Connected   *AtomBool
	TickMsg     emitter.Message
	KeyRequest  chan bool
}

type DisconnectMessage struct {
}

func (d *DisconnectMessage) Topic() string {
	return "client/disconnected/"
}

func (d *DisconnectMessage) Payload() []byte {
	return []byte{}
}

func New(emittercfg map[string]string) *Context {
	ctx := &Context{}
	ctx.EmitterCfg = emittercfg
	ctx.ChannelKeys = map[string]string{}
	ctx.InMsgs = make(chan emitter.Message)
	ctx.TickMsg = &TickMessage{}
	ctx.KeyRequest = make(chan bool)
	return ctx
}

type TickMessage struct {
}

func (p *TickMessage) Topic() string {
	return "tick"
}
func (p *TickMessage) Payload() []byte {
	return nil
}

type PubsubMessage struct {
	topic   string
	payload []byte
}

func (p *PubsubMessage) Topic() string {
	return p.topic
}
func (p *PubsubMessage) Payload() []byte {
	return p.payload
}

func (ctx *Context) Topic(msg *emitter.Message) string {
	return msg.Topic()
}

func (ctx *Context) Payload(msg *emitter.Message) []byte {
	return msg.Payload()
}
func (ctx *Context) Connect(username string, register, subscribe []string) error {
	// Create the options with default values
	options := emitter.NewClientOptions()
	options.SetUsername(username)

	// Set the message handler
	options.SetOnMessageHandler(func(client emitter.Emitter, msg emitter.Message) {
		ctx.InMsgs <- msg
	})

	// Set the presence notification handler
	options.SetOnPresenceHandler(func(_ emitter.Emitter, p emitter.PresenceEvent) {
		fmt.Printf("Occupancy: %v\n", p.Occupancy)
	})

	// Set the connection lost handler
	options.SetOnConnectionLostHandler(func(_ emitter.Emitter, e error) {
		msg := &DisconnectMessage{}
		ctx.InMsgs <- msg
	})

	options.SetOnKeyGenHandler(func(_ emitter.Emitter, r emitter.KeyGenResponse) {
		fmt.Printf("KeyGenResponse from emitter: '%s' = '%s' (status: %d)\n", r.Channel, r.Key, r.Status)
		ctx.ChannelKeys[r.Channel] = r.Key
		if r.Channel == "" || r.Key == "" {
			ctx.KeyRequest <- false
		} else {
			ctx.KeyRequest <- true
		}

	})

	options.AddBroker(ctx.EmitterCfg["host"])

	// Create a new emitter client and connect to the broker
	ctx.Client = emitter.NewClient(options)
	sToken := ctx.Client.Connect()
	if sToken.Error() != nil {
		return sToken.Error()
	}

	if !sToken.WaitTimeout(time.Second * 5) {
		return fmt.Errorf("Timeout when connecting to emitter.io server")
	}

	if ctx.Client.IsConnected() == false {
		return fmt.Errorf("Unknown error when connecting to emitter.io server")
	}

	for _, reg := range register {
		err := ctx.Register(reg)
		if err != nil {
			return err
		}
	}
	for _, sub := range subscribe {
		err := ctx.Subscribe(sub)
		if err != nil {
			return err
		}
	}

	go func() {
		for ctx.Connected.IsTrue() {
			time.Sleep(1)
			ctx.InMsgs <- ctx.TickMsg
		}
	}()

	return nil
}

func (ctx *Context) Register(channel string) error {

	keygenRequest := emitter.NewKeyGenRequest()
	keygenRequest.Key = ctx.EmitterCfg["secret"]
	keygenRequest.Channel = channel
	keygenRequest.Type = "rwslp"
	keygenToken := ctx.Client.GenerateKey(keygenRequest)
	if !keygenToken.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("Emitter.GenerateKey did not succeed for channel %s due to a timeout", channel)
	}

	var err error
	select {
	case request := <-ctx.KeyRequest:
		if !request {
			err = fmt.Errorf("Emitter.Register did not succeed for channel %s due to a fail", channel)
		}
	case <-time.After(5 * time.Second):
		err = fmt.Errorf("Emitter.Register did not succeed for channel %s due to a timeout", channel)
	}
	return err
}

func (ctx *Context) Subscribe(channel string) error {
	key, exists := ctx.ChannelKeys[channel]
	if exists {
		ctx.Client.Subscribe(key, channel)
		return nil
	}
	return fmt.Errorf("Emitter.Subscribe did not succeed for channel %s since no channel key was configured", channel)
}

func (ctx *Context) Publish(channel string, message string) error {
	key, exists := ctx.ChannelKeys[channel]
	if exists {
		ctx.Client.Publish(key, channel, message)
		return nil
	}
	return fmt.Errorf("Emitter.Publish did not succeed for channel %s since no channel key was configured", channel)
}

func (ctx *Context) PublishTTL(channel string, message string, ttl int) error {
	key, exists := ctx.ChannelKeys[channel]
	if exists {
		ctx.Client.PublishWithTTL(key, channel, message, ttl)
		return nil
	}
	return fmt.Errorf("Emitter.PublishTTL did not succeed for channel %s since no channel key was configured", channel)
}
