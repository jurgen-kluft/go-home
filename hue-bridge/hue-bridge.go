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

// Huecontext holds all necessary information
type Huecontext struct {
	key    string
	config *config.HueConfig

	log *logpkg.Logger
}

// New creates a new instance of hue instance
func New() *Huecontext {
	hue := &Huecontext{}
	hue.light = map[string]huelights.Light{}
	return hue
}

func (hue *Huecontext) initialize() (err error) {

	huebridge.SetLogger(os.Stdout)

	// Every 'device' (light/switch) is identical to 'handler'

	// For every 'device' register a handler:
	huebridge.Handle("test", func(req huebridge.Request, res *huebridge.Response) {
		fmt.Println("im handling test from", req.RemoteAddr, req.RequestedOnState)
		res.OnState = req.RequestedOnState

		// Set ErrorState to true to have the echo respond with "unable to reach device"
		// res.ErrorState = true
		return
	})

	// it is very important to use a full IP here or the UPNP does not work correctly.
	// one day ill fix this
	panic(huebridge.ListenAndServe("10.0.0.11:5000"))

	return nil
}

func main() {
	hue := New()

	hue.log = logpkg.New("hue-bridge")
	hue.log.AddEntry("emitter")
	hue.log.AddEntry("hue-bridge")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/huebridge/"}
		subscribe := []string{"config/hue/", "state/sensor/hue/", "sensor/light/hue/"}
		err := client.Connect("hue-bridge", register, subscribe)
		if err == nil {
			hue.log.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/huebridge/" {
						config, err := config.HueConfigFromJSON(string(msg.Payload()))
						if err == nil {
							hue.log.LogInfo("hue-bridge", "received configuration")
							hue.config = config
							err = hue.initialize()
							if err != nil {
								hue.log.LogError("hue-bridge", err.Error())
							}
						} else {
							hue.log.LogError("hue-bridge", err.Error())
						}
					} else if topic == "client/disconnected/" {
						hue.log.LogInfo("emitter", "disconnected")
						connected = false
					} else if topic == "state/sensor/hue/" {
						huesensor, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							hue.log.LogInfo("hue", "received flux")
							hue.CT = huesensor.GetFloatAttr("CT", 325.0)
							hue.BRI = huesensor.GetFloatAttr("BRI", 128.0)
						} else {
							hue.log.LogError("hue-bridge", err.Error())
						}
					} else if topic == "state/light/hue/" {
						if hue.config != nil {
							huesensor, err := config.SensorStateFromJSON(string(msg.Payload()))
							if err == nil {
								hue.log.LogInfo("hue-bridge", "received state")
								lightname := huesensor.GetValueAttr("name", "")
								if lightname != "" {
									light, exists := hue.light[lightname]
									if exists {
										huesensor.ExecValueAttr("power", func(power string) {
											if power == "on" {
												hue.lightState.On = new(bool)
												*hue.lightState.On = true
												*hue.lightState.CT = uint16(hue.CT)
												*hue.lightState.Bri = uint8(hue.BRI)
												hue.lights.SetLightState(light.ID, hue.lightState)
												hue.lightState.On = nil
											} else if power == "off" {
												*hue.lightState.On = false
												hue.lights.SetLightState(light.ID, hue.lightState)
												hue.lightState.On = nil
											}
										})
									}
								}
							} else {
								hue.log.LogError("hue-bridge", err.Error())
							}
						} else {
							hue.log.LogError("hue-bridge", fmt.Sprintf("error, receiving message on channel %s but we haven't received a configuration", topic))
						}
					}

				case <-time.After(time.Second * 60):
					hue.flux()
				}
			}
		}

		if err != nil {
			hue.log.LogError("hue-bridge", err.Error())
		}

		time.Sleep(5 * time.Second)
	}
}
