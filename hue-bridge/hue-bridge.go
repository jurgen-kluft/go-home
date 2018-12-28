package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
	huebridge "github.com/pborges/huejack"
)

/*
This service exists to have Alexa also be able to control other devices that are not supported directly.
Also to be able to have virtual 'variables/switches' that can be controlled with Alexa.
For example:

- Turning On/Off 'Livingroom TV'.
- Turning On/Off 'Bedroom TV'.
- Turning On/Off 'Story Mode'; which is a mode that puts certain light groups in a state for reading in the evening.
- Turning On/Off 'Holiday Mode'; which disables waking up kids and parents in the morning according to the calendar.
- Turning On/Off 'Flux Mode'; which disables automatic adjustments of lights
- Turning On/Off 'Bedroom ceiling light'
- Turning On/Off 'Bedroom chandelier'
- Turning On/Off 'Bedroom power switch'

Wild ideas:
- Turn On/Off 'Party Mode' / 'Halloween Mode' / 'Christmas Mode'
- Turn On/Off 'Music Mode' (use the MxChip Azure Devkit, can register sound ?)

Ok, so the configuration is mostly about defining 'variables' which are mostly routed to
service 'automation' which in turn will execute the logic.



*/

// context holds all necessary information
type context struct {
	name   string
	config *config.HueBridgeConfig
	vars   map[string]bool
	update []string
	log    *logpkg.Logger
}

// New creates a new instance of hue instance
func new() *context {
	c := &context{}
	c.name = "huebridge"
	c.vars = map[string]bool{}
	c.update = []string{}
	return c
}

func (c *context) initialize() (err error) {
	if c.config != nil {
		// huebridge.SetLogger(os.Stdout)
		// For every 'device' register a handler:
		for _, dev := range c.config.EmulatedDevices {
			huebridge.Handle(dev.Name, func(req huebridge.Request, res *huebridge.Response) {
				fmt.Println("HueBridge request from:", req.RemoteAddr, req.RequestedOnState)

				res.OnState = req.RequestedOnState
				c.vars[dev.Name] = req.RequestedOnState

				// These will be send on the "state/sensor/vars/" channel
				c.update = append(c.update, dev.Name)

				// Set ErrorState to true to have the echo respond with "unable to reach device"
				// res.ErrorState = true
				return
			})
		}
		// it is very important to use a full IP here or the UPNP does not work correctly.
		// one day ill fix this
		err = huebridge.ListenAndServe(c.config.IPPort)
	} else {
		err = fmt.Errorf("Cannot initialize Hue-Bridge without a configuration")
	}
	return err
}

func (c *context) getChannelsToRegister() []string {
	channels := []string{}
	for _, ch := range c.config.RegisterChannels {
		channels = append(channels, ch)
	}
	for _, dev := range c.config.EmulatedDevices {
		channels = append(channels, dev.Name)
	}
	return channels
}

func (c *context) getChannelsToSubscribe() []string {
	channels := []string{}
	channels = append(channels, "config/"+c.name+"/")
	for _, ch := range c.config.SubscribeChannels {
		channels = append(channels, ch)
	}
	return channels
}

func (c *context) process(client *pubsub.Context) {
	vars := config.NewSensorState("state.sensor.vars")
	for _, ev := range c.update {
		state := c.vars[ev]
		vars.AddBoolAttr(ev, state)
	}

	jsonstr, err := vars.ToJSON()
	if err == nil {
		fmt.Println(jsonstr)
		client.Publish("state/sensor/vars/", jsonstr)
	}
}

func main() {
	c := new()

	c.log = logpkg.New(c.name)
	c.log.AddEntry("emitter")
	c.log.AddEntry(c.name)

	for {
		client := pubsub.New(config.EmitterIOCfg)
		register := c.getChannelsToRegister()
		subscribe := c.getChannelsToSubscribe()
		err := client.Connect(c.name, register, subscribe)
		if err == nil {
			c.log.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/"+c.name+"/" {
						config, err := config.HueBridgeConfigFromJSON(string(msg.Payload()))
						if err == nil {
							c.log.LogInfo(c.name, "received configuration")
							c.config = config
							err = c.initialize()
							if err != nil {
								c.log.LogError(c.name, err.Error())
							}
						} else {
							c.log.LogError(c.name, err.Error())
						}
					} else if topic == "client/disconnected/" {
						c.log.LogInfo("emitter", "disconnected")
						connected = false
					}

				case <-time.After(time.Second * 1):
					// Drain events and send them to 'state/sensor/var'
					c.process(client)
				}
			}
		}

		if err != nil {
			c.log.LogError(c.name, err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}
