package main

// https://github.com/czerwe/gobravia

// Features:
// - Sony Bravia TVs: Turn On/Off

import (
	"fmt"
	"time"

	"github.com/czerwe/gobravia"
	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/pubsub"
)

type tv struct {
	name string
	host string
	mac  string
	key  string
	tv   *gobravia.BraviaTV
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

func (x *instance) AddTV(host string, mac string, name string) {
	tv := &tv{}
	tv.name = name
	tv.host = "10.0.0.77"
	tv.mac = "C4:3A:BE:95:0C:1E"

	tv.tv = gobravia.GetBravia(tv.host, "0000", tv.mac)
	tv.tv.GetCommands()

	x.tvs[name] = tv
}

func (x *instance) poweron(name string) {
	tv, exists := x.tvs[name]
	if exists {
		tv.tv.Poweron(tv.host)
	}
}
func (x *instance) poweroff(name string) {
	tv, exists := x.tvs[name]
	if exists {
		tv.tv.SendAlias("poweroff")
	}
}

func main() {
	sony := New()

	for {
		client := pubsub.New()
		err := client.Connect("tv.sony")
		if err == nil {

			fmt.Println("Connected to emitter")

			client.Subscribe("config/tv/sony")
			client.Subscribe("state/tv/sony")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/tv/sony" {
						//huelighting.config, err = config.HueConfigFromJSON(string(msg.Payload()))
					} else if topic == "state/tv/sony" {
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							power := state.GetValue("power", "idle")
							if power == "off" {
								sony.poweroff(state.Name)
							} else if power == "on" {
								sony.poweron(state.Name)
							}
						}
					} else if topic == "client/disconnected" {
						connected = false
					}

				case <-time.After(time.Second * 10):

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
