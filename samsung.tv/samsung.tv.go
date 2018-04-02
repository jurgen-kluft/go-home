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
	mac    string
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
	tv.remote, err = samote.Dial(tv.host, tv.name, tv.mac)
	if err == nil {
		x.tvs[id] = tv
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

	for {
		client := pubsub.New()
		err := client.Connect("tv.samsung")
		if err == nil {

			fmt.Println("Connected to emitter")
			client.Subscribe("config/tv/samsung")
			client.Subscribe("state/tv/samsung")

			for {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "tv.samsung/config" {
						//huelighting.config, err = config.HueConfigFromJSON(string(msg.Payload()))
					} else if topic == "tv.samsung/state" {
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							if state.Value == "off" {
								samsung.poweroff(state.Name)
							} else if state.Value == "on" {
								samsung.poweron(state.Name)
							}
						}
					}
					break
				case <-time.After(time.Second * 10):
					// do something if messages are taking too long
					// or if we haven't received enough state info.

					break
				}
			}
		} else {
			fmt.Println(err.Error())
		}

		// Wait for 10 seconds before retrying
		fmt.Println("Connecting to emitter (retry every 10 seconds)")
		time.Sleep(10 * time.Second)
	}
}
