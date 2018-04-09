package main

// Features:
// - Samsung TV: Turn On/Off

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/saljam/samote"
)

type tv struct {
	host   string
	name   string
	id     string
	remote samote.Remote
}

type instance struct {
	tvs map[string]*tv
}

// New ...
func New() *instance {
	x := &instance{}
	x.tvs = map[string]*tv{}
	return x
}

func (x *instance) Add(host string, name string, id string) error {

	tv := &tv{}
	tv.host = host
	tv.name = name
	tv.id = id

	var err error
	tv.remote, err = samote.Dial(tv.host, tv.name, tv.id)
	if err == nil {
		x.tvs[name] = tv
	}
	return err
}

func (x *instance) poweron(name string) {
	tv, exists := x.tvs[name]
	if exists {
		tv.remote.SendKey(samote.KEY_POWERON)
	}
}
func (x *instance) poweroff(name string) {
	tv, exists := x.tvs[name]
	if exists {
		tv.remote.SendKey(samote.KEY_POWEROFF)
	}
}

func main() {
	samsung := New()
	err := samsung.Add("10.0.0.76:55000", "Bedroom Samsung-TV", "Remote")
	if err != nil {
		fmt.Println(err)
	}

	logger := logpkg.New("samsung.tv")
	logger.AddEntry("emitter")
	logger.AddEntry("samsung.tv")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/tv/samsung/", "state/tv/samsung/"}
		subscribe := []string{"config/tv/samsung/", "state/tv/samsung/"}
		err := client.Connect("tv.samsung", register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/tv/samsung/" {
						//huelighting.config, err = config.HueConfigFromJSON(string(msg.Payload()))
					} else if topic == "state/tv/samsung/" {
						logger.LogInfo("samsung.tv", "received configuration")
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							power := state.GetValueAttr("power", "idle")
							if power == "off" {
								samsung.poweroff(state.Name)
							} else if power == "on" {
								samsung.poweron(state.Name)
							}
						}
					} else if topic == "client/disconnected/" {
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}

				case <-time.After(time.Second * 10):

				}
			}
		}

		if err != nil {
			logger.LogError("samsung.tv", err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}
