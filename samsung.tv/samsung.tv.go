package main

// Features:
// - Samsung TV: Turn On/Off

import (
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/saljam/samote"
)

type instance struct {
	config *config.SamsungTVConfig
	tvs    map[string]samote.Remote
}

// New ...
func New() *instance {
	x := &instance{}
	x.config = nil
	x.tvs = map[string]samote.Remote{}
	return x
}

func (x *instance) Add(host string, name string, id string) error {
	var err error
	var remote samote.Remote
	remote, err = samote.Dial(host, name, id)
	if err == nil {
		x.tvs[name] = remote
	}
	return err
}

func (x *instance) poweron(name string) {
	remote, exists := x.tvs[name]
	if exists {
		remote.SendKey(samote.KEY_POWERON)
	}
}
func (x *instance) poweroff(name string) {
	remote, exists := x.tvs[name]
	if exists {
		remote.SendKey(samote.KEY_POWEROFF)
	}
}

func main() {
	samsung := New()

	logger := logpkg.New("samsung.tv")
	logger.AddEntry("emitter")
	logger.AddEntry("samsung.tv")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/samsung.tv/", "state/samsung.tv/"}
		subscribe := []string{"config/samsung.tv/", "state/samsung.tv/"}
		err := client.Connect("tv.samsung", register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/samsung.tv/" {
						logger.LogInfo("samsung.tv", "received configuration")
						samsung.config, err = config.SamsungTVConfigFromJSON(string(msg.Payload()))
						for _, tv := range samsung.config.Devices {
							err = samsung.Add(tv.IP, tv.Name, tv.ID)
							if err == nil {
								logger.LogInfo("samsung.tv", "registered TV with name "+tv.Name)
							} else {
								logger.LogError("samsung.tv", err.Error())
							}
						}
					} else if topic == "state/samsung.tv/" {
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
