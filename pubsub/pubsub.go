package pubsub

import (
	"fmt"

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
	ctx.Secret = config.SecretKey
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

	// Create a new emitter client and connect to the broker
	ctx.Client = emitter.NewClient(options)
	sToken := ctx.Client.Connect()

	if sToken.Wait() && sToken.Error() == nil {
		return nil
	}

	return sToken.Error()
}

func (ctx *Context) Subscribe(channel string) error {
	ctx.Client.Subscribe(ctx.Secret, channel)
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
