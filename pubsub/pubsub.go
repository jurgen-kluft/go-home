package pubsub

import (
	"fmt"
	"time"

	emitter "github.com/emitter-io/go"
	"github.com/jurgen-kluft/go-home/config"
)

type Context struct {
	Secret string
	InMsgs chan emitter.Message
	Client emitter.Emitter
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
	ctx.Secret = config.EmitterSecretKey
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

	options.AddBroker("tcp://10.0.0.22:8080")

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

func (ctx *Context) Subscribe(key string, channel string) error {
	ctx.Client.Subscribe(key, channel)
	return nil
}

func (ctx *Context) Publish(channel string, message string) error {
	ctx.Client.Publish(ctx.Secret, channel, message)
	return nil
}

func (ctx *Context) PublishTTL(channel string, message string, ttl int) error {
	ctx.Client.PublishWithTTL(ctx.Secret, channel, message, ttl)
	return nil
}
