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
			client.Subscribe(config.EmitterSensorLightChannelKey, "sensor/tv.sony/+")

			for {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "tv.sony/config" {
						//huelighting.config, err = config.HueConfigFromJSON(string(msg.Payload()))
					} else if topic == "tv.sony/state" {
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							if state.Value == "off" {
								sony.poweroff(state.Name)
							} else if state.Value == "on" {
								sony.poweron(state.Name)
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
