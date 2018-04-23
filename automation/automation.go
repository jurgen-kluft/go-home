package main

// All automation logic is in this package
// Here we react to:
// - presence (people arriving/leaving)
// - switches (pressed)
// - events (timeofday, calendar)
// - time-based logic (morning 6:20 turn on bedroom lights)

import (
	"strings"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
)

func main() {
	auto := New()

	logger := logpkg.New("automation")
	logger.AddEntry("emitter")
	logger.AddEntry("automation")

	for {
		// The state channels we are interested in
		stateChannels := []string{"state/xiaomi/", "state/presence/", "state/hue/", "state/yee/", "state/xiaomi/", "state/bravia.tv/", "state/samsung.tv/",
			"state/sensor/sophia", "state/sensor/jennifer", "state/sensor/parents",
		}

		auto.pubsub = pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/automation/"}
		register = append(register, stateChannels...)
		subscribe := []string{"config/automation/"}
		subscribe = append(subscribe, stateChannels...)
		err := auto.pubsub.Connect("automation", register, subscribe)

		if err == nil {
			logger.LogInfo("emitter", "connected")
			connected := true
			for connected {
				select {
				case msg := <-auto.pubsub.InMsgs:
					topic := msg.Topic()
					if topic == "config/automation/" {
					} else if topic == "client/disconnected/" {
						connected = false
						logger.LogInfo("emitter", "disconnected")
					} else if strings.HasPrefix(topic, "state/") {
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							auto.HandleEvent(state)
						} else {
							logger.LogError("automation", err.Error())
						}
					}
				case <-time.After(time.Second * 30):
					auto.HandleTime(time.Now())
				}
			}
		}
		if err != nil {
			logger.LogError("automation", err.Error())
		}

		// Wait for 5 seconds before retrying
		time.Sleep(5 * time.Second)
	}
}

type Automation struct {
	pubsub              *pubsub.Context
	sensors             map[string]string
	presence            map[string]bool
	timeofday           string
	lastmotion          time.Time
	lastmotionFrontDoor time.Time
	timedActions        map[string]TimedAction
}

func New() *Automation {
	auto := &Automation{}
	auto.sensors = map[string]string{}
	auto.presence = map[string]bool{}
	return auto
}

func (a *Automation) UpdateTimedActions(now time.Time) {
	for name, ta := range a.timedActions {
		if ta.Tick(now) == true {
			ta.Action(a)
			delete(a.timedActions, name)
		}
	}
}

func (a *Automation) FamilyIsHome() bool {
	// This can be a lot more complex:

	// Motion sensors will mark if people are home and this will be reset when the front-door openened.
	// When presence shows no devices on the network but there is still motion detected we still should
	// regard that family is at home.
	if len(a.presence) == 0 {
		if time.Now().Sub(a.lastmotion) > (time.Minute * 12) {
			return false
		}
	}

	return len(a.presence) > 0
}

func (a *Automation) IsSensor(name string, value string) bool {
	v, e := a.sensors[name]
	return e && (v == value)
}

func (a *Automation) TurnOnLight(name string) error {
	state := config.NewSensorState(name)
	state.AddStringAttr("power", "on")
	jsonstr, err := state.ToJSON()
	a.pubsub.Publish("state/light/hue/", jsonstr)
	return err
}
func (a *Automation) TurnOffLight(name string) {

}
func (a *Automation) ToggleLight(name string) {

}
func (a *Automation) TurnOffSwitch(name string) {

}
func (a *Automation) ToggleSwitch(name string) {

}
func (a *Automation) TurnOffTV(name string) {

}

func (a *Automation) HandleEvent(state *config.SensorState) {
	name := state.Name
	what := state.GetValueAttr("type", "")
	if name == "timeofday" {
		value := state.GetValueAttr("timeofday", "")
		a.HandleTimeOfDay(value)
	} else if what == "switch" {
		a.HandleSwitch(name, state)
	} else if what == "sensor" {
		a.HandleSensor(name, state)
	} else if what == "presence" {
		a.HandlePresence(name, state)
	}
}

func (a *Automation) HandleTimeOfDay(to string) {
	if to != a.timeofday {
		a.timeofday = to
		switch to {
		case "morning":
			a.TurnOffLight("Kitchen")
			a.TurnOffLight("Living Room")
		case "lunch":
			if a.FamilyIsHome() {
				a.TurnOnLight("Kitchen")
			}
		case "afternoon":
			a.TurnOffLight("Kitchen")
		case "evening":
			a.TurnOnLight("Kitchen")
			a.TurnOnLight("Living Room")
		case "bedtime":
			if a.FamilyIsHome() {
				if a.IsSensor("jennifer", "school") {
					a.TurnOnLight("Jennifer")
				}
				if a.IsSensor("sophia", "school") {
					a.TurnOnLight("Sophia")
				}
			}
		case "sleeptime":
			if a.FamilyIsHome() {
				a.TurnOnLight("Bedroom")

				if a.IsSensor("jennifer", "school") {
					a.TurnOffLight("Jennifer")
				}
				if a.IsSensor("sophia", "school") {
					a.TurnOffLight("Sophia")
				}
			}
		case "night":
			if a.IsSensor("jennifer", "school") {
				a.TurnOffLight("Bedroom")
				a.TurnOffLight("Bedroom")
				a.TurnOffLight("Bedroom")
			}
			a.TurnOffLight("Kitchen")
			a.TurnOffLight("Living Room")
			a.TurnOffLight("Jennifer")
			a.TurnOffLight("Sophia")
			a.TurnOffLight("Front door hall light")
		}
	}
}

func (a *Automation) HandleSensor(name string, state *config.SensorState) {
	if name == "motion_sensor_kitchen" || name == "motion_sensor_livingroom" {
		value := state.GetValueAttr("motion", "")
		if value == "on" {
			a.lastmotion = time.Now() // Update the time we last detected motion
			if name == "motion_sensor_kitchen" {
				a.lastmotionFrontDoor = time.Now()
			}
			if a.timeofday == "breakfast" {
				a.TurnOnLight("Kitchen")
				a.TurnOnLight("Living Room")
			}
		} else if value == "off" {
			if time.Now().Sub(a.lastmotionFrontDoor) > time.Minute*5 {
				a.TurnOffLight("Front door hall light")
			}
		}
	}
	if name == "magnet_sensor_frontdoor" {
		value := state.GetValueAttr("state", "")
		if value == "open" {
			a.TurnOnLight("Front door hall light")
			a.lastmotionFrontDoor = time.Now()
		}
	}
}

func (a *Automation) HandleSwitch(name string, state *config.SensorState) {
	if name == "Bedroom Switch" {
		value := state.GetValueAttr("click", "")
		if value == "double click" {
			a.ToggleLight("Bedroom")
		}
		if value == "single click" {
			a.TurnOffSwitch("Bedroom ceiling light")
			a.TurnOffSwitch("Bedroom chandelier")
		}
		if value == "press release" {
			a.ToggleSwitch("Bedroom power plug")
		}
	} else if name == "Sophia Switch" {
		value := state.GetValueAttr("click", "")
		if value == "single click" {
			a.ToggleLight("Sophia")
		}
	}
}

func (a *Automation) HandlePresence(name string, state *config.SensorState) {
	if state.StringAttrs != nil {
		anyone_home := (len(a.presence) == 0)
		for _, attr := range *state.StringAttrs {
			name := attr.Name
			value := attr.Value
			if value == "away" {
				delete(a.presence, name)
			} else {
				a.presence[name] = true
			}
		}
		if anyone_home && len(a.presence) == 0 {
			a.HandlePresenceLeaving()
		} else if !anyone_home && len(a.presence) != 0 {
			a.HandlePresenceArriving()
		}
	}
}

func (a *Automation) HandlePresenceLeaving() {
	// Turn off everything
	a.TurnOffEverything()
}

func (a *Automation) HandlePresenceArriving() {
	// Depending on time-of-day
	// Turn on Kitchen
	// Turn on Living-Room

}

func (a *Automation) TurnOffEverything() {
	a.TurnOffLight("Kitchen")
	a.TurnOffLight("Living Room")
	a.TurnOffLight("Bedroom")
	a.TurnOffLight("Jennifer")
	a.TurnOffLight("Sophia")
	a.TurnOffLight("Front door hall light")

	a.TurnOffSwitch("Bedroom power plug")
	a.TurnOffSwitch("Bedroom chandelier")
	a.TurnOffSwitch("Bedroom ceiling")

	a.TurnOffTV("Samsung bedroom")
	a.TurnOffTV("Sony livingroom")
}

func (a *Automation) HandleTime(now time.Time) {
	if a.IsSensor("sensor.calendar.jennifer", "school") {
		if now.Hour() == 6 && now.Minute() == 20 {
			a.TurnOnLight("Bedroom")
		}
		if now.Hour() == 6 && now.Minute() == 30 {
			a.TurnOnLight("Jennifer")
		}
	}
	if a.IsSensor("sensor.calendar.sophia", "school") {
		if now.Hour() == 7 && now.Minute() == 10 {
			a.TurnOnLight("Bedroom")
		}
		if now.Hour() == 7 && now.Minute() == 20 {
			a.TurnOnLight("Sophia")
		}
	}
	if a.IsSensor("sensor.calendar.parents", "work") {
		if !a.IsSensor("sensor.calendar.jennifer", "school") && !a.IsSensor("sensor.calendar.sophia", "school") {
			if now.Hour() == 7 && now.Minute() == 45 {
				a.TurnOnLight("Bedroom")
			}
		}
	}

	a.UpdateTimedActions(now)
}

type TimedAction interface {
	Tick(time.Time) bool
	Action(*Automation)
}

type TurnOffLightTimedAction struct {
	When        time.Time
	NameOfLight string
}

func (ta *TurnOffLightTimedAction) Tick(now time.Time) bool {
	return now.After(ta.When)
}
func (ta *TurnOffLightTimedAction) Action(auto *Automation) {
	auto.TurnOnLight(ta.NameOfLight)
}
