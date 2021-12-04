package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jurgen-kluft/go-home/conbee.sensors/deconz"
	"github.com/jurgen-kluft/go-home/conbee.sensors/deconz/event"
	"github.com/jurgen-kluft/go-home/config"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
)

/*
This process will scan for events from Conbee, mainly sensors and will send those as
sensor states over NATS.
*/

type motionSensorState struct {
	Name     string
	ID       string
	LastSeen time.Time
	Motion   bool
	Conbee   config.ConbeeMotion
}

type contactSensorState struct {
	Name     string
	ID       string
	LastSeen time.Time
	Contact  bool
	Conbee   config.ConbeeContact
}

type switchState struct {
	Name     string
	ID       string
	LastSeen time.Time
	Button   string
	Conbee   config.ConbeeSwitch
}

type fullState struct {
	switches       map[string]switchState
	motionSensors  map[string]motionSensorState
	contactSensors map[string]contactSensorState
}

func fullStateFromConfig(c *config.ConbeeSensorsConfig) fullState {
	full := fullState{}
	full.switches = make(map[string]switchState)
	full.motionSensors = make(map[string]motionSensorState)
	full.contactSensors = make(map[string]contactSensorState)

	for _, e := range c.Switches {
		state := switchState{Name: e.Name, ID: e.ID, LastSeen: time.Now(), Button: "", Conbee: e}
		full.switches[state.ID] = state
	}
	for _, e := range c.Motion {
		state := motionSensorState{Name: e.Name, ID: e.ID, LastSeen: time.Now(), Motion: false, Conbee: e}
		full.motionSensors[state.ID] = state
	}
	for _, e := range c.Contact {
		state := contactSensorState{Name: e.Name, ID: e.ID, LastSeen: time.Now(), Contact: false, Conbee: e}
		full.contactSensors[state.ID] = state
	}

	return full
}

type signal_t int32

func (b *signal_t) set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32((*int32)(b), int32(i))
}
func (b *signal_t) is_true() bool {
	return atomic.LoadInt32((*int32)(b)) != 0
}
func (b *signal_t) is_not_true() bool {
	return atomic.LoadInt32((*int32)(b)) == 0
}

func async_conbee(ctx context.Context, sg *signal_t, cc *config.ConbeeSensorsConfig, mm *microservice.Service) {
	mm.Logger.LogInfo(mm.Name, fmt.Sprintf("Connecting to deCONZ at %s with API key %s", fmt.Sprintf("http://%s/api", cc.Host), cc.APIKey))

	defer sg.set(false)

	eventChan, err := eventChan(fmt.Sprintf("http://%s/api", cc.Host), cc.APIKey)
	if err != nil {
		return
	}
	defer close(eventChan)

	mm.Logger.LogInfo(mm.Name, fmt.Sprintf("Connected to deCONZ at %s", fmt.Sprintf("http://%s/api", cc.Host)))

	fullState := fullStateFromConfig(cc)
	for {

		select {
		case <-ctx.Done():
			return

		case ev := <-eventChan:
			if ev.State == nil || event.IsEmptyState(ev.State) {
				// An empty state for one of the devices
			} else {
				cstate, exist := fullState.contactSensors[ev.UniqueID]
				if exist {
					dstate := ev.State.(*event.ZHAOpenClose)
					if dstate != nil {
						mm.Logger.LogInfo(mm.Name, fmt.Sprintf("contact:  %s -> %v = %v", cstate.Name, cstate.Contact, dstate.Open))
						cstate.Contact = dstate.Open
						fullState.contactSensors[ev.UniqueID] = cstate
						if dstate.Open {
							msg := &microservice.Message{Topic: "state/sensor/", Payload: []byte(cstate.Conbee.Open)}
							mm.ProcessMessages <- msg
						} else {
							msg := &microservice.Message{Topic: "state/sensor/", Payload: []byte(cstate.Conbee.Close)}
							mm.ProcessMessages <- msg
						}
					}
				} else {
					mstate, exist := fullState.motionSensors[ev.UniqueID]
					if exist {
						dstate := ev.State.(*event.ZHAPresence)
						if dstate != nil {
							mm.Logger.LogInfo(mm.Name, fmt.Sprintf("motion:  %s -> %v = %v", mstate.Name, mstate.Motion, dstate.Presence))
							mstate.Motion = dstate.Presence
							fullState.motionSensors[ev.UniqueID] = mstate
							if dstate.Presence {
								msg := &microservice.Message{Topic: "state/sensor/", Payload: []byte(mstate.Conbee.On)}
								mm.ProcessMessages <- msg
							} else {
								msg := &microservice.Message{Topic: "state/sensor/", Payload: []byte(mstate.Conbee.Off)}
								mm.ProcessMessages <- msg
							}
						}
					} else {
						sstate, exist := fullState.switches[ev.UniqueID]
						if exist {
							dstate := ev.State.(*event.ZHASwitch)
							if dstate != nil {
								mm.Logger.LogInfo(mm.Name, fmt.Sprintf("switch:  %s -> %v = %v", sstate.Name, sstate.Button, dstate.ButtonEventAsString()))
								sstate.Button = dstate.ButtonEventAsString()
								fullState.switches[ev.UniqueID] = sstate
								if sstate.Button == config.SwitchSingleClick {
									msg := &microservice.Message{Topic: "state/sensor/", Payload: []byte(sstate.Conbee.SingleClick)}
									mm.ProcessMessages <- msg
								} else if sstate.Button == config.SwitchDoubleClick {
									msg := &microservice.Message{Topic: "state/sensor/", Payload: []byte(sstate.Conbee.DoubleClick)}
									mm.ProcessMessages <- msg
								} else if sstate.Button == config.SwitchTrippleClick {
									msg := &microservice.Message{Topic: "state/sensor/", Payload: []byte(sstate.Conbee.TrippleClick)}
									mm.ProcessMessages <- msg
								}
							}
						}
					}
				}
			}
		}
	}
}

func main() {
	var cc *config.ConbeeSensorsConfig = nil
	var nc *config.ConbeeSensorsConfig = nil

	var alive signal_t

	for {
		var err error

		cc = nc

		register := []string{"config/request/", "config/conbee/sensors/"}
		subscribe := []string{"config/conbee/sensors/"}

		if cc != nil {
			register = append(register, cc.SensorsOut)
		}

		// context.WithCancel returns a copy of parent with a new Done channel.
		// The returned context's Done channel is closed when the returned cancel function is called or when the parent
		// context's Done channel is closed, whichever happens first.
		ctx, cancel := context.WithCancel(context.Background())

		m := microservice.New("conbee/sensors")
		m.RegisterAndSubscribe(register, subscribe)

		m.RegisterHandler("config/conbee/sensors/", func(m *microservice.Service, topic string, msg []byte) bool {
			m.Logger.LogInfo(m.Name, "Received configuration, schedule restart")
			nc, err = config.ConbeeSensorsConfigFromJSON(msg)
			if err != nil {
				m.Logger.LogError(m.Name, err.Error())
			} else {
				cc = nil
				return false
			}
			return true
		})

		m.RegisterHandler("state/sensor/", func(m *microservice.Service, topic string, sensor_state_payload []byte) bool {
			m.Pubsub.PublishStr(cc.SensorsOut, string(sensor_state_payload))
			return true
		})

		tickCount := 0
		m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
			if (tickCount % 30) == 0 {
				if nc == nil {
					m.Logger.LogInfo(m.Name, "Requesting configuration..")
					m.Pubsub.PublishStr("config/request/", m.Name)
				}
			} else if (tickCount % 9) == 0 {
				if cc != nil {
					if alive.is_not_true() {
						m.Logger.LogInfo(m.Name, "Conbee routine is not running, schedule restart..")
						// seems that async_conbee go routine is not running
						// micro-service exit and restart
						return false
					}
				}
			}
			tickCount++
			return true
		})

		if cc != nil {
			alive.set(true)
			go async_conbee(ctx, &alive, cc, m)
		}

		m.Loop()

		// signal our (running) async conbee go routine to exit
		cancel()

		for alive.is_true() {
			time.Sleep(1 * time.Second)
		}

		// Sleep for a while before restarting
		m.Logger.LogInfo(m.Name, "Waiting 10 seconds before re-starting..")
		time.Sleep(10 * time.Second)

		m.Logger.LogInfo(m.Name, "Re-start..")
	}

}

func eventChan(addr string, APIkey string) (chan *deconz.DeviceEvent, error) {
	// get an event reader from the API
	d := deconz.API{Config: deconz.Config{Addr: addr, APIKey: APIkey}}
	reader, err := d.EventReader()
	if err != nil {
		return nil, err
	}

	// Dial the reader
	err = reader.Dial()
	if err != nil {
		return nil, err
	}

	// create a new reader, embedding the event reader
	eventReader := d.DeviceEventReader(reader)
	channel := make(chan *deconz.DeviceEvent)
	// start it, it starts its own thread
	eventReader.Start(channel)
	// return the channel
	return channel, nil
}
