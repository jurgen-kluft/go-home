package pubsub

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type AtomInt int32

func (b *AtomInt) Set(value int32) {
	atomic.StoreInt32((*int32)(b), value)
}

func (b *AtomInt) Get() int32 {
	return atomic.LoadInt32((*int32)(b))
}

type Msg struct {
	Subject string
	Data    []byte
}

// Context contains the necessary information to run the MQTT client
type Context struct {
	Config        map[string]string
	InMsgs        chan Msg
	SubToIndex    map[string]int
	SubChannels   []string
	Client        mqtt.Client
	Connected     *AtomInt
	TickFrequency time.Duration
	Tick          Msg
}

func New(config map[string]string, tickFrequency time.Duration) *Context {
	ctx := &Context{}
	ctx.Config = config
	ctx.SubToIndex = make(map[string]int)
	ctx.SubChannels = make([]string, 0, 10)
	ctx.Connected = new(AtomInt)
	ctx.Connected.Set(-1)
	ctx.TickFrequency = tickFrequency
	ctx.Tick = Msg{Subject: "tick", Data: []byte("do your thing")}

	var broker = config["mqtt.broker.host"]
	var port = config["mqtt.broker.port"]
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(config["mqtt.broker.clientId"])
	opts.SetUsername(config["mqtt.broker.username"])
	opts.SetPassword(config["mqtt.broker.password"])
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {})

	opts.OnConnect = func(client mqtt.Client) {
		msg := Msg{Subject: "client/connected", Data: []byte(fmt.Sprintf("tcp://%s:%s", broker, port))}
		if ctx.Connected.Get() == -1 {
			// subscribe to all registered topics
			for _, channel := range ctx.SubChannels {
				subschannel := strings.Replace(channel, ".", "/", -1)
				subschannel = strings.TrimSuffix(subschannel, "/")
				subschannel = strings.TrimPrefix(subschannel, "/")
				if err := ctx.Subscribe(subschannel); err != nil {
					fmt.Printf("Error subscribing to %s: %v\n", subschannel, err)
				}
			}
		}
		ctx.Connected.Set(1)
		ctx.InMsgs <- msg
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		ctx.Connected.Set(0)
		msg := Msg{Subject: "client/disconnected"}
		ctx.InMsgs <- msg
	}

	ctx.Client = mqtt.NewClient(opts)

	return ctx
}

func (ctx *Context) Topic(msg Msg) string {
	return msg.Subject
}

func (ctx *Context) Payload(msg Msg) []byte {
	return msg.Data
}

func (ctx *Context) Connect(username string, register, subscribe []string) error {
	var err error

	ctx.InMsgs = make(chan Msg, 128)
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

		go func() {
			for ctx.Connected.Get() == 1 {
				time.Sleep(time.Duration(1) * time.Second)
				ctx.InMsgs <- ctx.Tick
			}
		}()
	}

	return err
}

func (ctx *Context) Close() {
	ctx.Connected.Set(0)
	ctx.Client.Disconnect(100)
	ctx.Client = nil
	time.Sleep(1)
}

func (ctx *Context) Register(channel string) error {
	return nil
}

func (ctx *Context) Subscribe(channel string) (err error) {
	_, exists := ctx.SubToIndex[channel]
	if !exists {
		sub := strings.Replace(channel, ".", "/", -1)
		sub = strings.TrimSuffix(sub, "/")
		sub = strings.TrimPrefix(sub, "/")
		ctx.Client.Subscribe(sub, 0, nil)
		index := len(ctx.SubChannels)
		ctx.SubToIndex[channel] = index
		ctx.SubChannels = append(ctx.SubChannels, sub)
		return err
	}
	return fmt.Errorf("PubSub.Subscribe failed for channel %s", channel)
}

func (ctx *Context) PublishStr(channel string, message string) error {
	index, exists := ctx.SubToIndex[channel]
	if exists {
		ctx.Client.Publish(ctx.SubChannels[index], 0, false, []byte(message))
		return nil
	}
	return fmt.Errorf("PubSub.Publish failed for channel %s", channel)
}

func (ctx *Context) Publish(channel string, message []byte) error {
	index, exists := ctx.SubToIndex[channel]
	if exists {
		ctx.Client.Publish(ctx.SubChannels[index], 0, false, message)
		return nil
	}
	return fmt.Errorf("PubSub.Publish failed for channel %s", channel)
}
