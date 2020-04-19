package pubsub

import (
	"fmt"
	nats "github.com/nats-io/nats.go"
)

type Context struct {
	Config map[string]string
	InMsgs chan *nats.Msg
	Subs   map[string]*nats.Subscription
	Client *nats.Conn
}

func New(config map[string]string) *Context {
	ctx := &Context{}
	ctx.Config = config
	ctx.Subs = map[string]*nats.Subscription{}
	return ctx
}

func (ctx *Context) Connect(username string, register, subscribe []string) error {
	var err error

	ctx.Client, err = nats.Connect(ctx.Config["host"],
		nats.Name(username),
		nats.Token(ctx.Config["secret"]),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			msg := &nats.Msg{Subject: "client.disconnected"}
			ctx.InMsgs <- msg
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			msg := &nats.Msg{Subject: "client.reconnected", Data: []byte(nc.ConnectedUrl())}
			ctx.InMsgs <- msg
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			msg := &nats.Msg{Subject: "client.closed"}
			ctx.InMsgs <- msg
		}),
	)

	ctx.InMsgs = make(chan *nats.Msg, 128)

	return err
}

func (ctx *Context) Close() {
	ctx.Client.Close()
	ctx.Client = nil
	ctx.InMsgs = nil
}

func (ctx *Context) Register(channel string) error {
	return nil
}

func (ctx *Context) Subscribe(channel string) error {
	_, exists := ctx.Subs[channel]
	if !exists {
		sub, err := ctx.Client.ChanSubscribe(channel, ctx.InMsgs)
		ctx.Subs[channel] = sub
		return err
	}
	return fmt.Errorf("PubSub.Register failed for channel %s", channel)
}

func (ctx *Context) Publish(channel string, message string) error {
	_, exists := ctx.Subs[channel]
	if exists {
		ctx.Client.Publish(channel, []byte(message))
		return nil
	}
	return fmt.Errorf("PubSub.Publish failed for channel %s", channel)
}

func (ctx *Context) PublishTTL(channel string, message string, ttl int) error {
	_, exists := ctx.Subs[channel]
	if exists {
		ctx.Client.Publish(channel, []byte(message))
		return nil
	}
	return fmt.Errorf("PubSub.PublishTTL failed for channel %s", channel)
}
