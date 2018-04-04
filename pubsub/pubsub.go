package pubsub

import (
	"fmt"
	"time"

	emitter "github.com/emitter-io/go"
	"github.com/jurgen-kluft/go-home/config"
)

type Context struct {
	Secret      string
	ChannelKeys map[string]string
	InMsgs      chan emitter.Message
	Client      emitter.Emitter
}

type DisconnectMessage struct {
}

func (d *DisconnectMessage) Topic() string {
	return "client/disconnected"
}
func (d *DisconnectMessage) Payload() []byte {
	return []byte{}
}

func New() *Context {
	ctx := &Context{}
	ctx.ChannelKeys = map[string]string{}
	ctx.Secret = config.EmitterSecrets["secret"]
	ctx.InMsgs = make(chan emitter.Message)
	return ctx
}

func (ctx *Context) Connect(username string) error {
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
	})

	options.AddBroker("tcp://127.0.0.1:8080")

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

	return nil
}

func (ctx *Context) Register(channel string) error {
	_, exists := ctx.ChannelKeys[channel]
	if !exists {
		keygenRequest := emitter.NewKeyGenRequest()
		keygenRequest.Key = "hVtJ6AfHC_UoZvy1DYHJf8gUnHsnKiKG"
		keygenRequest.Channel = channel + "/"
		keygenToken := ctx.Client.GenerateKey(keygenRequest)
		if !keygenToken.Wait() {
			return fmt.Errorf("Emitter.Register did not succeed for channel %s due to a timeout", channel)
		}
		_, exists = ctx.ChannelKeys[channel]
		if !exists {
			return fmt.Errorf("Emitter.Register did not succeed for channel %s, did not receive keygen response", channel)
		}
	}
	return nil
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
	key, exists := config.EmitterSecrets[channel]
	if exists {
		ctx.Client.Publish(key, channel, message)
		return nil
	}
	return fmt.Errorf("Emitter.Publish did not succeed for channel %s since no channel key was configured", channel)
}

func (ctx *Context) PublishTTL(channel string, message string, ttl int) error {
	key, exists := config.EmitterSecrets[channel]
	if exists {
		ctx.Client.PublishWithTTL(key, channel, message, ttl)
		return nil
	}
	return fmt.Errorf("Emitter.PublishTTL did not succeed for channel %s since no channel key was configured", channel)
}
