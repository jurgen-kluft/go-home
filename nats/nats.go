package pubsub

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	server "github.com/nats-io/nats.go"
)

type AtomBool int32

func (b *AtomBool) Set(value bool) {
	var i int32
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

// Context contains the necessary information to run the Nats client
type Context struct {
	Config        map[string]string
	InMsgs        chan *server.Msg
	SubToIndex    map[string]int
	Subscriptions []*server.Subscription
	SubChannels   []string
	Client        *server.Conn
	Connected     *AtomBool
	TickFrequency time.Duration
	Tick          *server.Msg
}

func New(config map[string]string) *Context {
	ctx := &Context{}
	ctx.Config = config
	ctx.SubToIndex = make(map[string]int)
	ctx.Subscriptions = make([]*server.Subscription, 0, 10)
	ctx.SubChannels = make([]string, 0, 10)
	ctx.Connected = new(AtomBool)
	ctx.TickFrequency = time.Duration(1) * time.Second
	ctx.Tick = &server.Msg{Subject: "tick/", Data: nil}
	return ctx
}

func (ctx *Context) Topic(msg *server.Msg) string {
	return msg.Subject
}

func (ctx *Context) Payload(msg *server.Msg) []byte {
	return msg.Data
}

func (ctx *Context) Connect(username string, register, subscribe []string) error {
	var err error

	ctx.Client, err = server.Connect(ctx.Config["host"],
		server.Name(username),
		server.Token(ctx.Config["secret"]),
		server.DisconnectErrHandler(func(nc *server.Conn, err error) {
			ctx.Connected.Set(false)
			msg := &server.Msg{Subject: "client/disconnected/"}
			ctx.InMsgs <- msg
		}),
		server.ReconnectHandler(func(nc *server.Conn) {
			msg := &server.Msg{Subject: "client/reconnected/", Data: []byte(nc.ConnectedUrl())}
			ctx.InMsgs <- msg
		}),
		server.ClosedHandler(func(nc *server.Conn) {
			msg := &server.Msg{Subject: "client/closed/"}
			ctx.InMsgs <- msg
		}),
	)

	ctx.InMsgs = make(chan *server.Msg, 128)
	for _, s := range subscribe {
		err := ctx.Subscribe(s)
		if err != nil {
			return err
		}
	}

	if err == nil {
		for _, r := range register {
			ctx.Register(r)
		}

		ctx.Connected.Set(true)
		go func() {
			for ctx.Connected.IsTrue() {
				time.Sleep(ctx.TickFrequency)
				ctx.InMsgs <- ctx.Tick
			}
		}()
	}

	return err
}

func (ctx *Context) Close() {
	ctx.Connected.Set(false)
	ctx.Client.Close()
	ctx.Client = nil
	time.Sleep(1)
}

func (ctx *Context) Register(channel string) error {
	_, exists := ctx.SubToIndex[channel]
	if !exists {
		subschannel := strings.Replace(channel, "/", ".", -1)
		subschannel = strings.TrimSuffix(subschannel, ".")
		index := len(ctx.SubChannels)
		ctx.SubToIndex[channel] = index
		ctx.Subscriptions = append(ctx.Subscriptions, nil)
		ctx.SubChannels = append(ctx.SubChannels, subschannel)
	}
	return nil
}

func (ctx *Context) Subscribe(channel string) (err error) {
	_, exists := ctx.SubToIndex[channel]
	if !exists {
		subschannel := strings.Replace(channel, "/", ".", -1)
		subschannel = strings.TrimSuffix(subschannel, ".")
		subscription, err := ctx.Client.ChanSubscribe(subschannel, ctx.InMsgs)
		index := len(ctx.SubChannels)
		ctx.SubToIndex[channel] = index
		ctx.Subscriptions = append(ctx.Subscriptions, subscription)
		ctx.SubChannels = append(ctx.SubChannels, subschannel)
		return err
	}
	return fmt.Errorf("PubSub.Subscribe failed for channel %s", channel)
}

func (ctx *Context) PublishStr(channel string, message string) error {
	index, exists := ctx.SubToIndex[channel]
	if exists {
		ctx.Client.Publish(ctx.SubChannels[index], []byte(message))
		return nil
	}
	return fmt.Errorf("PubSub.Publish failed for channel %s", channel)
}

func (ctx *Context) Publish(channel string, message []byte) error {
	index, exists := ctx.SubToIndex[channel]
	if exists {
		ctx.Client.Publish(ctx.SubChannels[index], message)
		return nil
	}
	return fmt.Errorf("PubSub.Publish failed for channel %s", channel)
}

func (ctx *Context) PublishTTLStr(channel string, message string, ttl int) error {
	index, exists := ctx.SubToIndex[channel]
	if exists {
		ctx.Client.Publish(ctx.SubChannels[index], []byte(message))
		return nil
	}
	return fmt.Errorf("PubSub.PublishTTL failed for channel %s", channel)
}

func (ctx *Context) PublishTTL(channel string, message []byte, ttl int) error {
	index, exists := ctx.SubToIndex[channel]
	if exists {
		ctx.Client.Publish(ctx.SubChannels[index], message)
		return nil
	}
	return fmt.Errorf("PubSub.PublishTTL failed for channel %s", channel)
}
