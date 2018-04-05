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
	tv.host = host
	tv.mac = mac

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
	sony.AddTV("10.0.0.77", "C4:3A:BE:95:0C:1E", "Livingroom TV")
	sony.poweroff("Livingroom TV")

	for {
		client := pubsub.New("tcp://10.0.0.22:8080")
		register := []string{"config/tv/sony/", "state/tv/sony/"}
		subscribe := []string{"config/tv/sony/", "state/tv/sony/"}
		err := client.Connect("bravia.tv", register, subscribe)
		if err == nil {
			fmt.Println("Connected to emitter")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/tv/sony/" {
						//huelighting.config, err = config.HueConfigFromJSON(string(msg.Payload()))
					} else if topic == "state/tv/sony/" {
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							power := state.GetValueAttr("power", "idle")
							if power == "off" {
								sony.poweroff(state.Name)
							} else if power == "on" {
								sony.poweron(state.Name)
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
			time.Sleep(1 * time.Second)
		}

		// Wait for 5 seconds before retrying
		fmt.Println("Connecting to emitter (retry every 5 seconds)")
		time.Sleep(5 * time.Second)
	}
}
