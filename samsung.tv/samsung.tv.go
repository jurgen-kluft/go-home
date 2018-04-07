package main

// Features:
// - Samsung TV: Turn On/Off

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
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

	for {
		client := pubsub.New("tcp://10.0.0.22:8080")
		register := []string{"config/tv/samsung/", "state/tv/samsung/"}
		subscribe := []string{"config/tv/samsung/", "state/tv/samsung/"}
		err := client.Connect("tv.samsung", register, subscribe)
		if err == nil {
			fmt.Println("Connected to emitter")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/tv/samsung/" {
						//huelighting.config, err = config.HueConfigFromJSON(string(msg.Payload()))
					} else if topic == "state/tv/samsung/" {
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
						connected = false
					}

				case <-time.After(time.Second * 10):

				}
			}
		}

		if err != nil {
			fmt.Println("Error: " + err.Error())
			time.Sleep(5 * time.Second)
		}
	}
}
