package main

// All automation logic is in this package
// Here we react to:
// - presence (people arriving/leaving)
// - switches (pressed)
// - events (timeofday, calendar)
// - time-based logic (morning 6:20 turn on bedroom lights)

import (
	"fmt"
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
			"state/sensor/timeofday", "state/sensor/sophia", "state/sensor/jennifer", "state/sensor/parents",
		}

		auto.pubsub = pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/automation/", "shout/message/"}
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
							auto.HandleEvent(topic, state)
						} else {
							logger.LogError("automation", err.Error())
						}
					}
				case <-time.After(time.Second * 30):
					now := time.Now()
					auto.HandleTime(now)
					auto.UpdateTimedActions(now)
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
	pubsub                      *pubsub.Context
	sensors                     map[string]string
	presence                    map[string]bool
	timeofday                   string
	lastseenMotionInHouse       time.Time
	lastseenMotionInKitchenArea time.Time
	lastseenMotionInBedroom     time.Time
	timedActions                map[string]TimedAction
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
	familyIsHome := false
	for _, pres := range a.presence {
		familyIsHome = familyIsHome || pres
	}
	if !familyIsHome {
		if time.Now().Sub(a.lastseenMotionInHouse) > (time.Minute * 5) {
			return false
		}
		return true
	}

	return familyIsHome
}

func (a *Automation) IsSensor(name string, value string) bool {
	v, e := a.sensors[name]
	return e && (v == value)
}

func (a *Automation) TurnOnLight(name string) {

}
func (a *Automation) TurnOffLight(name string) {

}
func (a *Automation) ToggleLight(name string) {

}
func (a *Automation) TurnOnSwitch(name string) {

}
func (a *Automation) TurnOffSwitch(name string) {

}
func (a *Automation) ToggleSwitch(name string) {

}
func (a *Automation) TurnOffPlug(name string) {

}
func (a *Automation) TogglePlug(name string) {

}
func (a *Automation) TurnOffTV(name string) {

}

func (a *Automation) HandleEvent(channel string, state *config.SensorState) {
	if channel == "state/sensor/timeofday/" {
		value := state.GetValueAttr("timeofday", "")
		a.HandleTimeOfDay(value)
	} else if channel == "state/xiaomi/" {
		name := state.Name
		if name == config.SophiaRoomSwitch || name == config.BedroomSwitch {
			a.HandleSwitch(name, state)
		} else if name == config.KitchenMotionSensor || name == config.LivingroomMotionSensor || name == config.BedroomMotionSensor {
			a.HandleMotionSensor(name, state)
		} else if name == config.FrontdoorMagnetSensor {
			a.HandleMagnetSensor(name, state)
		}
	} else if channel == "state/presence/" {
		a.HandlePresence(state)
	}

}

func (a *Automation) HandleTimeOfDay(to string) {
	if to != a.timeofday {
		a.timeofday = to
		switch to {
		case "morning":
			a.sendNotification("Turning off lights since it is morning")
			a.TurnOffLight(config.KitchenLights)
			a.TurnOffLight(config.LivingroomLights)
		case "lunch":
			if a.FamilyIsHome() {
				a.sendNotification("Turning on lights since it is noon and someone is home")
				a.TurnOnLight(config.KitchenLights)
			}
		case "afternoon":
			a.sendNotification("Turning off lights since it is afternoon")
			a.TurnOffLight(config.KitchenLights)
		case "evening":
			a.TurnOnLight(config.KitchenLights)
			a.TurnOnLight(config.LivingroomLights)
		case "bedtime":
			if a.FamilyIsHome() {
				if a.IsSensor("jennifer", "school") {
					a.TurnOnLight(config.JenniferRoomLights)
				}
				if a.IsSensor("sophia", "school") {
					a.TurnOnLight(config.SophiaRoomLights)
				}
			}
		case "sleeptime":
			if a.FamilyIsHome() {
				a.TurnOnLight("Bedroom")

				if a.IsSensor("jennifer", "school") {
					a.TurnOffLight(config.JenniferRoomLights)
				}
				if a.IsSensor("sophia", "school") {
					a.TurnOffLight(config.SophiaRoomLights)
				}
			}
		case "night":
			if a.IsSensor("jennifer", "school") {
				a.TurnOffLight(config.BedroomLights)
			}
			a.TurnOffLight(config.KitchenLights)
			a.TurnOffLight(config.LivingroomLights)
			a.TurnOffLight(config.JenniferRoomLights)
			a.TurnOffLight(config.SophiaRoomLights)
			a.TurnOffLight(config.FrontdoorHallLight)
		}
	}
}

func (a *Automation) HandleSwitch(name string, state *config.SensorState) {
	if name == "Bedroom Switch" {
		value := state.GetValueAttr("click", "")
		if value == "double click" {
			a.ToggleLight(config.BedroomLights)
		}
		if value == "single click" {
			a.TurnOffSwitch(config.BedroomCeilingLightSwitch)
			a.TurnOffSwitch(config.BedroomChandelierLightSwitch)
		}
		if value == "press release" {
			a.TogglePlug(config.BedroomPowerPlug)
		}
	} else if name == "Sophia Switch" {
		value := state.GetValueAttr("click", "")
		if value == "single click" {
			a.ToggleLight(config.SophiaRoomLights)
		}
	}
}

func (a *Automation) HandleMotionSensor(name string, state *config.SensorState) {
	now := time.Now()
	if name == config.KitchenMotionSensor || name == config.LivingroomMotionSensor {
		value := state.GetValueAttr("motion", "")
		if value == "on" {
			a.lastseenMotionInHouse = now // Update the time we last detected motion
			if name == config.KitchenMotionSensor {
				a.lastseenMotionInKitchenArea = now
				a.ActivateTurnOffLightTimedAction("Turnoff front door hall light", "Front door hall light", 400*time.Second)
			}
			if a.timeofday == "breakfast" {
				a.TurnOnLight(config.KitchenLights)
				a.TurnOnLight(config.LivingroomLights)
			}

			if a.timeofday == "night" {
				if name == config.KitchenMotionSensor {
					a.sendNotification("Turning on kitchen and livingroom lights since it is night and there is movement in the kitchen area")
				} else if name == config.LivingroomMotionSensor {
					a.sendNotification("Turning on kitchen and livingroom lights since it is night and there is movement in the livingroom area")
				}
				a.TurnOnLight(config.KitchenLights)
				a.TurnOnLight(config.LivingroomLights)
				a.ActivateTurnOffLightTimedAction(config.KitchenLights, config.KitchenLights, 5*time.Minute)
				a.ActivateTurnOffLightTimedAction(config.LivingroomLights, config.LivingroomLights, 5*time.Minute)
			}

		}
	} else if name == config.BedroomMotionSensor {
		value := state.GetValueAttr("motion", "")
		lastseenDuration := now.Sub(a.lastseenMotionInBedroom)
		if value == "on" {
			a.lastseenMotionInHouse = now // Update the time we last detected motion
			a.lastseenMotionInBedroom = now

			if a.timeofday == "evening" || a.timeofday == "bedtime" {
				if lastseenDuration > (time.Duration(5) * time.Minute) {
					a.TurnOnLight(config.BedroomLights)
					a.TurnOnSwitch(config.BedroomChandelierLightSwitch)
				}
			}
		} else if value == "off" {
			if a.timeofday != "night" && a.timeofday != "sleeptime" {
				if lastseenDuration > (time.Duration(10) * time.Minute) {
					a.TurnOffLight(config.BedroomLights)
					a.TurnOffSwitch(config.BedroomChandelierLightSwitch)
					a.TurnOffSwitch(config.BedroomCeilingLightSwitch)
				}
			}
		}
	}

}

func (a *Automation) HandleMagnetSensor(name string, state *config.SensorState) {
	if name == "Front Door Magnet" {
		value := state.GetValueAttr("state", "")
		if value == "open" {
			a.TurnOnLight("Front door hall light")
			a.ActivateTurnOffLightTimedAction("Turnoff front door hall light", "Front door hall light", 400*time.Second)
			a.lastseenMotionInKitchenArea = time.Now()
		} else if value == "close" {
			a.ActivateTurnOffLightTimedAction("Turnoff front door hall light", "Front door hall light", 400*time.Second)
		}
	}
}

func (a *Automation) sendNotification(message string) {
	a.pubsub.Publish("shout/message/", message)
}

func (a *Automation) updatePresence(name string, presence bool) (current bool, previous bool) {
	var exists bool
	previous, exists = a.presence[name]
	if !exists {
		previous = false
		a.presence[name] = presence
	}
	return presence, previous
}

func (a *Automation) HandlePresence(state *config.SensorState) {
	if state.StringAttrs != nil {
		anyonehome := a.FamilyIsHome()
		for _, attr := range *state.StringAttrs {
			presence, previous := a.updatePresence(attr.Name, attr.Value == "home")
			if presence != previous {
				if presence == false {
					a.sendNotification(fmt.Sprintf("%s is not home", attr.Name))
				} else {
					a.sendNotification(fmt.Sprintf("%s is home", attr.Name))
				}
			}
		}
		if !anyonehome && a.FamilyIsHome() {
			// Depending on time-of-day
			switch a.timeofday {
			case "lunch":
				a.sendNotification("Turning on kitchen lights since it is noon and someone came home")
				a.TurnOnLight(config.KitchenLights)
			case "evening":
				a.sendNotification("Turning on kitchen and livingroom lights since it is evening and someone came home")
				a.TurnOnLight(config.KitchenLights)
				a.TurnOnLight(config.LivingroomLights)
			case "bedtime":
				a.sendNotification("Turning on kitchen and livingroom lights since it is bedtime and someone came home")
				a.TurnOnLight(config.KitchenLights)
				a.TurnOnLight(config.LivingroomLights)
			case "sleeptime":
				a.sendNotification("Turning on kitchen, livingroom and bedroom lights since it is bedtime and someone came home")
				a.TurnOnLight(config.KitchenLights)
				a.TurnOnLight(config.LivingroomLights)
				a.TurnOnLight(config.BedroomLights)
			}

		} else if anyonehome && !a.FamilyIsHome() {
			// Turn off everything
			a.TurnOffEverything()
		}
	}
}

func (a *Automation) TurnOffEverything() {
	a.TurnOffLight(config.KitchenLights)
	a.TurnOffLight(config.LivingroomLights)
	a.TurnOffLight(config.BedroomLights)
	a.TurnOffLight(config.JenniferRoomLights)
	a.TurnOffLight(config.SophiaRoomLights)
	a.TurnOffLight(config.FrontdoorHallLight)

	a.TurnOffPlug(config.BedroomPowerPlug)
	a.TurnOffSwitch(config.BedroomChandelierLightSwitch)
	a.TurnOffSwitch(config.BedroomCeilingLightSwitch)

	a.TurnOffTV("Samsung bedroom")
	a.TurnOffTV("Sony livingroom")
}

func (a *Automation) HandleTime(now time.Time) {

	jenniferHasSchool := a.IsSensor("sensor.calendar.jennifer", "school")
	sophiaHasSchool := a.IsSensor("sensor.calendar.jennifer", "school")
	parentsHaveToWork := a.IsSensor("sensor.calendar.parents", "work")

	if jenniferHasSchool {
		if now.Hour() == 6 && now.Minute() == 20 {
			a.sendNotification("Waking up Parents for Jennifer")
			a.TurnOnLight(config.BedroomLights)
		}
		if now.Hour() == 6 && now.Minute() == 30 {
			a.sendNotification("Waking up Jennifer")
			a.TurnOnLight(config.JenniferRoomLights)
		}
	}
	if sophiaHasSchool && !jenniferHasSchool {
		if now.Hour() == 7 && now.Minute() == 10 {
			a.sendNotification("Waking up Parents for Sophia")
			a.TurnOnLight(config.BedroomLights)
		}
		if now.Hour() == 7 && now.Minute() == 20 {
			a.sendNotification("Waking up Sophia")
			a.TurnOnLight(config.SophiaRoomLights)
		}
	}
	if parentsHaveToWork {
		if !sophiaHasSchool && !jenniferHasSchool {
			if now.Hour() == 7 && now.Minute() == 30 {
				a.sendNotification("Waking up Parents")
				a.TurnOnLight(config.BedroomLights)
			}
		}
	}
}

type TimedAction interface {
	Tick(time.Time) bool
	Action(*Automation)
}

type TurnOffLightTimedAction struct {
	When        time.Time
	NameOfLight string
}

func (a *Automation) ActivateTurnOffLightTimedAction(nameOfAction string, nameOfLight string, duration time.Duration) {
	ta, exists := a.timedActions[nameOfAction]
	if !exists {
		ta = &TurnOffLightTimedAction{NameOfLight: nameOfLight, When: time.Now().Add(duration)}
		a.timedActions[nameOfAction] = ta
	} else {
		turnOffLightAction := ta.(*TurnOffLightTimedAction)
		turnOffLightAction.When = time.Now().Add(duration)
	}
}

func (ta *TurnOffLightTimedAction) Tick(now time.Time) bool {
	return now.After(ta.When)
}
func (ta *TurnOffLightTimedAction) Action(auto *Automation) {
	auto.TurnOnLight(ta.NameOfLight)
}
