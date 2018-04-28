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
					auto.now = time.Now()
					auto.HandleTime()
					auto.UpdateTimedActions()
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
	config                      *config.AutomationConfig
	sensors                     map[string]string
	presence                    map[string]bool
	timeofday                   string
	now                         time.Time
	lastseenMotionInHouse       time.Time
	lastseenMotionInKitchenArea time.Time
	lastseenMotionInBedroom     time.Time
	timedActions                map[string]*TimedAction
}

func New() *Automation {
	auto := &Automation{}
	auto.sensors = map[string]string{}
	auto.presence = map[string]bool{}
	return auto
}

func (a *Automation) UpdateTimedActions() {
	for name, ta := range a.timedActions {
		if ta.Tick(a.now) == true {
			ta.Action(ta, a)
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

func (a *Automation) SensorHasValue(name string, value string) bool {
	v, e := a.sensors[name]
	return e && (v == value)
}

func (a *Automation) TurnOnDevice(name string) error {
	dc, exists := a.config.DeviceControlCache["name"]
	if exists {
		err := a.pubsub.Publish(dc.Channel, dc.On)
		return err
	}
	return fmt.Errorf("device with name %s doesn't exist", name)
}
func (a *Automation) TurnOffDevice(name string) error {
	dc, exists := a.config.DeviceControlCache["name"]
	if exists {
		err := a.pubsub.Publish(dc.Channel, dc.Off)
		return err
	}
	return fmt.Errorf("device with name %s doesn't exist", name)
}
func (a *Automation) ToggleDevice(name string) error {
	dc, exists := a.config.DeviceControlCache["name"]
	if exists {
		err := a.pubsub.Publish(dc.Channel, dc.Toggle)
		return err
	}
	return fmt.Errorf("device with name %s doesn't exist", name)
}

func (a *Automation) HandleEvent(channel string, state *config.SensorState) {
	var sensortype string
	var sensorname string
	n, _ := fmt.Sscanf(channel, "state/%s/%s/", &sensortype, &sensorname)
	if n == 2 && sensortype == "sensor" {
		a.sensors[sensorname] = state.GetValueAttr("value", "")
		if sensorname == "timeofday" {
			a.HandleTimeOfDay(a.sensors[sensorname])
		}
	} else if sensortype == "xiaomi" {
		name := state.Name
		if name == config.SophiaRoomSwitch || name == config.BedroomSwitch {
			a.HandleSwitch(name, state)
		} else if name == config.KitchenMotionSensor || name == config.LivingroomMotionSensor || name == config.BedroomMotionSensor {
			a.HandleMotionSensor(name, state)
		} else if name == config.FrontdoorMagnetSensor {
			a.HandleMagnetSensor(name, state)
		}
	} else if sensortype == "presence" {
		a.HandlePresence(state)
	}
}

func WakeUpParentsForJennifer(ta *TimedAction, a *Automation) {
	a.sendNotification("Waking up Parents for Jennifer")
	a.TurnOnDevice(config.BedroomLights)
}
func WakeUpParentsForSophia(ta *TimedAction, a *Automation) {
	a.sendNotification("Waking up Parents for Sophia")
	a.TurnOnDevice(config.BedroomLights)
}
func WakeUpParentsForWork(ta *TimedAction, a *Automation) {
	a.sendNotification("Waking up Parents")
	a.TurnOnDevice(config.BedroomLights)
}
func WakeUpJennifer(ta *TimedAction, a *Automation) {
	a.sendNotification("Waking up Jennifer")
	a.TurnOnDevice(config.JenniferRoomLights)
}
func WakeUpSophia(ta *TimedAction, a *Automation) {
	a.sendNotification("Waking up Sophia")
	a.TurnOnDevice(config.SophiaRoomLights)
}
func TurnOnFrontdoorHallLight(ta *TimedAction, a *Automation) {
	a.TurnOnDevice(config.FrontdoorHallLight)
}

func (a *Automation) AddTimedAction(name string, hour int, minute int, ad ActionDelegate) {
	action, exists := a.timedActions[name]
	if !exists {
		when := time.Date(a.now.Year(), a.now.Month(), a.now.Day(), hour, minute, 0, 0, a.now.Location())
		action = &TimedAction{Name: name, When: when, Action: ad}
		a.timedActions[name] = action
	}
}

func (a *Automation) HandleTimeOfDay(to string) {
	if to != a.timeofday {
		a.timeofday = to
		switch to {
		case "breakfast":
			jenniferHasSchool := a.SensorHasValue("sensor.calendar.jennifer", "school")
			sophiaHasSchool := a.SensorHasValue("sensor.calendar.sophia", "school")
			parentsHaveToWork := a.SensorHasValue("sensor.calendar.parents", "work")
			if jenniferHasSchool {
				a.AddTimedAction("Waking up Parents for Jennifer", 6, 20, WakeUpParentsForJennifer)
			} else if sophiaHasSchool {
				a.AddTimedAction("Waking up Parents for Sophia", 7, 0, WakeUpParentsForSophia)
			} else if parentsHaveToWork {
				a.AddTimedAction("Waking up Parents", 7, 30, WakeUpParentsForWork)
			}
			if jenniferHasSchool {
				a.AddTimedAction("Waking up Jennifer", 6, 30, WakeUpJennifer)
				a.AddTimedAction("Turn on Hall Light", 7, 11, TurnOnFrontdoorHallLight)
			}
			if sophiaHasSchool {
				a.AddTimedAction("Waking up Sophia", 7, 10, WakeUpSophia)
			}
		case "morning":
			a.sendNotification("Turning off lights since it is morning")
			a.TurnOffDevice(config.KitchenLights)
			a.TurnOffDevice(config.LivingroomLights)
		case "lunch":
			if a.FamilyIsHome() {
				a.sendNotification("Turning on lights since it is noon and someone is home")
				a.TurnOnDevice(config.KitchenLights)
			}
		case "afternoon":
			a.sendNotification("Turning off lights since it is afternoon")
			a.TurnOffDevice(config.KitchenLights)
		case "evening":
			a.TurnOnDevice(config.KitchenLights)
			a.TurnOnDevice(config.LivingroomLights)
		case "bedtime":
			if a.FamilyIsHome() {
				if a.SensorHasValue("jennifer", "school") {
					a.TurnOnDevice(config.JenniferRoomLights)
				}
				if a.SensorHasValue("sophia", "school") {
					a.TurnOnDevice(config.SophiaRoomLights)
				}
			}
		case "sleeptime":
			if a.FamilyIsHome() {
				a.TurnOnDevice(config.BedroomLights)

				if a.SensorHasValue("jennifer", "school") {
					a.TurnOffDevice(config.JenniferRoomLights)
				}
				if a.SensorHasValue("sophia", "school") {
					a.TurnOffDevice(config.SophiaRoomLights)
				}
			}
		case "night":
			if a.SensorHasValue("jennifer", "school") {
				a.TurnOffDevice(config.BedroomLights)
			}
			a.TurnOffDevice(config.KitchenLights)
			a.TurnOffDevice(config.LivingroomLights)
			a.TurnOffDevice(config.JenniferRoomLights)
			a.TurnOffDevice(config.SophiaRoomLights)
			a.TurnOffDevice(config.FrontdoorHallLight)
		}
	}
}

func (a *Automation) HandleSwitch(name string, state *config.SensorState) {
	if name == config.BedroomSwitch {
		value := state.GetValueAttr("click", "")
		if value == "double click" {
			a.ToggleDevice(config.BedroomLights)
		}
		if value == "single click" {
			a.TurnOffDevice(config.BedroomCeilingLightSwitch)
			a.TurnOffDevice(config.BedroomChandelierLightSwitch)
		}
		if value == "press release" {
			a.ToggleDevice(config.BedroomPowerPlug)
		}
	} else if name == config.SophiaRoomSwitch {
		value := state.GetValueAttr("click", "")
		if value == "single click" {
			a.ToggleDevice(config.SophiaRoomLights)
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
				a.ActivateTurnOffDeviceTimedAction("Turnoff front door hall light", "Front door hall light", 400*time.Second)
			}
			if a.timeofday == "breakfast" {
				a.TurnOnDevice(config.KitchenLights)
				a.TurnOnDevice(config.LivingroomLights)
			}

			if a.timeofday == "night" {
				if name == config.KitchenMotionSensor {
					a.sendNotification("Turning on kitchen and livingroom lights since it is night and there is movement in the kitchen area")
				} else if name == config.LivingroomMotionSensor {
					a.sendNotification("Turning on kitchen and livingroom lights since it is night and there is movement in the livingroom area")
				}
				a.TurnOnDevice(config.KitchenLights)
				a.TurnOnDevice(config.LivingroomLights)
				a.ActivateTurnOffDeviceTimedAction(config.KitchenLights, config.KitchenLights, 5*time.Minute)
				a.ActivateTurnOffDeviceTimedAction(config.LivingroomLights, config.LivingroomLights, 5*time.Minute)
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
					a.TurnOnDevice(config.BedroomLights)
					a.TurnOnDevice(config.BedroomChandelierLightSwitch)
				}
			}
		} else if value == "off" {
			if a.timeofday != "night" && a.timeofday != "sleeptime" {
				if lastseenDuration > (time.Duration(10) * time.Minute) {
					a.TurnOffDevice(config.BedroomLights)
					a.TurnOffDevice(config.BedroomChandelierLightSwitch)
					a.TurnOffDevice(config.BedroomCeilingLightSwitch)
				}
			}
		}
	}

}

func (a *Automation) HandleMagnetSensor(name string, state *config.SensorState) {
	if name == config.FrontdoorMagnetSensor {
		value := state.GetValueAttr("state", "")
		if value == "open" {
			a.TurnOnDevice(config.FrontdoorHallLight)
			a.ActivateTurnOffDeviceTimedAction("Turnoff front door hall light", "Front door hall light", 400*time.Second)
			a.lastseenMotionInKitchenArea = time.Now()
		} else if value == "close" {
			a.ActivateTurnOffDeviceTimedAction("Turnoff front door hall light", "Front door hall light", 400*time.Second)
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
		for _, attr := range state.StringAttrs {
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
				a.TurnOnDevice(config.KitchenLights)
			case "evening":
				a.sendNotification("Turning on kitchen and livingroom lights since it is evening and someone came home")
				a.TurnOnDevice(config.KitchenLights)
				a.TurnOnDevice(config.LivingroomLights)
			case "bedtime":
				a.sendNotification("Turning on kitchen and livingroom lights since it is bedtime and someone came home")
				a.TurnOnDevice(config.KitchenLights)
				a.TurnOnDevice(config.LivingroomLights)
			case "sleeptime":
				a.sendNotification("Turning on kitchen, livingroom and bedroom lights since it is bedtime and someone came home")
				a.TurnOnDevice(config.KitchenLights)
				a.TurnOnDevice(config.LivingroomLights)
				a.TurnOnDevice(config.BedroomLights)
			}

		} else if anyonehome && !a.FamilyIsHome() {
			// Turn off everything
			a.TurnOffEverything()
		}
	}
}

func (a *Automation) TurnOffEverything() {
	a.TurnOffDevice(config.KitchenLights)
	a.TurnOffDevice(config.LivingroomLights)
	a.TurnOffDevice(config.BedroomLights)
	a.TurnOffDevice(config.JenniferRoomLights)
	a.TurnOffDevice(config.SophiaRoomLights)
	a.TurnOffDevice(config.FrontdoorHallLight)
	a.TurnOffDevice(config.BedroomPowerPlug)
	a.TurnOffDevice(config.BedroomChandelierLightSwitch)
	a.TurnOffDevice(config.BedroomCeilingLightSwitch)

	a.TurnOffDevice(config.BedroomSamsungTV)
	a.TurnOffDevice(config.LivingroomBraviaTV)
}

func (a *Automation) HandleTime() {

}

type ActionDelegate func(ta *TimedAction, a *Automation)

func TurnOffDeviceTimedAction(ta *TimedAction, a *Automation) {
	a.TurnOnDevice(ta.Name)
}

type TimedAction struct {
	Name   string
	When   time.Time
	Action ActionDelegate
}

func (a *Automation) ActivateTurnOffDeviceTimedAction(nameOfAction string, nameOfLight string, duration time.Duration) {
	ta, exists := a.timedActions[nameOfAction]
	if !exists {
		ta = &TimedAction{Name: nameOfLight, When: time.Now().Add(duration), Action: TurnOffDeviceTimedAction}
		a.timedActions[nameOfAction] = ta
	} else {
		turnOffLightAction := ta
		turnOffLightAction.When = time.Now().Add(duration)
	}
}

func (ta *TimedAction) Tick(now time.Time) bool {
	return now.After(ta.When)
}
