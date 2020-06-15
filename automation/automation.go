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
	microservice "github.com/jurgen-kluft/go-home/micro-service"
)

func main() {
	register := []string{}
	subscribe := []string{"config/automation/"}

	m := microservice.New("automation")
	m.RegisterAndSubscribe(register, subscribe)

	auto := new()
	auto.service = m

	m.RegisterHandler("config/automation/", func(m *microservice.Service, topic string, msg []byte) bool {
		// Register used channels and subscribe to channels we are interested in
		config, err := config.AutomationConfigFromJSON(msg)
		if err == nil {
			auto.config = config
			// Register used channels
			for _, ss := range auto.config.ChannelsToRegister {
				if err = m.Pubsub.Register(ss); err != nil {
					m.Logger.LogError(m.Name, err.Error())
				}
			}
			// Subscribe channels
			for _, ss := range auto.config.SubChannels {
				if err = m.Pubsub.Subscribe(ss); err != nil {
					m.Logger.LogError(m.Name, err.Error())
				}
			}
		} else {
			m.Logger.LogError(m.Name, err.Error())
		}
		return true
	})

	m.RegisterHandler("*", func(m *microservice.Service, topic string, msg []byte) bool {
		if strings.HasPrefix(topic, "state") {
			state, err := config.SensorStateFromJSON(msg)
			if err == nil {
				auto.handleEvent(topic, state)
			} else {
				m.Logger.LogError(m.Name, err.Error())
			}
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if (tickCount & 0x1) == 0 {
			// Every 10 seconds
			if auto.config != nil {
				auto.now = time.Now()
				auto.presenceDetection()
				auto.updateTimedActions()
			}
		}
		if (tickCount % 30) == 0 {
			if auto.config == nil {
				m.Pubsub.PublishStr("config/request/", m.Name)
			}
		}
		tickCount++
		return true
	})

	m.Loop()
}

type homePresence struct {
	peopleAreHome          bool
	detectionStamp         time.Time
	detectionState         string
	detectionDelayDuration time.Duration
	detectionEvalResult    bool
	detectionEvalDuration  time.Duration
	presence               map[string]bool
}

func newPresence() *homePresence {
	h := &homePresence{}
	h.peopleAreHome = true
	h.detectionState = "Open/Closed"
	h.detectionStamp = time.Now()
	h.detectionDelayDuration = time.Minute * 15
	h.detectionEvalResult = false
	h.detectionEvalDuration = time.Minute * 15
	h.presence = map[string]bool{}
	return h
}

// reset() should be called when the front-door is opened->closed because this
// can indicate people have left the house. This will start a procedure:
// - Wait for 10 minutes to start an evaluation window
// - In the evaluation window we determine if there are people home by Wifi and Motion
// - If after the evaluation window nothing is detected we mark 'peopleAreHome' as false
// - After the evaluation state a scan state is started that will keep looking at
//   Wifi and Motion etc..
func (h *homePresence) frontDoorOpenClosed() {
	h.detectionState = "Open/Closed"
	h.detectionStamp = time.Now()
}

// determineIfPeopleAreHome is the only function that is allowed to change
// the variable 'peopleAreHome'!
func (h *homePresence) determineIfPeopleAreHome(now time.Time) {
	if h.detectionState == "OpenClosed" {
		// The door was opened/closed so that means people are/where at home
		h.peopleAreHome = true

		if now.Sub(h.detectionStamp) > h.detectionDelayDuration {
			h.detectionState = "Evaluate"
			h.detectionStamp = time.Now()
			h.detectionEvalResult = false
		}
	} else if h.detectionState == "Evaluate" {
		if now.Sub(h.detectionStamp) > h.detectionEvalDuration {
			// After the evaluation window, did we detect any motion, switch presses or Wifi Presence ?
			// If we did not then it means nobody is home.
			if h.detectionEvalResult == false {
				h.peopleAreHome = false
			}
			h.detectionState = "Scan"
		} else {
			// Look at the Wifi presence?, no! (WIFI can stay connected until people are out of the building)
			if h.peopleAreHome == false {
				for _, prsnc := range h.presence {
					h.peopleAreHome = h.peopleAreHome || prsnc
				}
			}
		}
	} else if h.detectionState == "Scan" {
		// After the detection and evaluation window (door was closed + N minutes) we can look at the wifi-presence again.
		if h.peopleAreHome == false {
			for _, prsnc := range h.presence {
				h.peopleAreHome = h.peopleAreHome || prsnc
			}
		}
	}
}

// reportCausation() should be called when a causation is detected like a button press, light switch press or movement
func (h *homePresence) reportCausation(now time.Time) {
	if h.detectionState == "Evaluate" {
		// Any detected presence after the detection window but within the evaluation
		// window means that there are people home. Indicate this by setting the
		// evaluation windows result to true.
		h.detectionEvalResult = true
	} else if h.detectionState == "Scan" {
		if h.peopleAreHome == false {
			// Any detected presence after the detection window means that people are home
			// This kind of presence is a button, motion, light switch or plug press
			h.peopleAreHome = true
		}
	}
}

type automation struct {
	config                      *config.AutomationConfig
	sensors                     map[string]string
	timeofday                   string
	now                         time.Time
	lastseenMotionInHouse       time.Time
	lastseenMotionInKitchenArea time.Time
	lastseenMotionInBedroom     time.Time
	timedActions                map[string]*timedBasedAction
	motionBasedActions          map[string]*motionBasedAction
	presence                    *homePresence
	service                     *microservice.Service
}

func new() *automation {
	auto := &automation{}
	auto.sensors = map[string]string{}
	auto.presence = newPresence()
	return auto
}

func (a *automation) peopleAreHome() bool {
	return a.presence.peopleAreHome
}

func (a *automation) sensorHasValue(name string, value string) bool {
	v, e := a.sensors[name]
	return e && (v == value)
}

func (a *automation) turnOnDevice(name string) error {
	dc, exists := a.config.DeviceControlCache["name"]
	if exists {
		err := a.service.Pubsub.PublishStr(dc.Channel, dc.On)
		return err
	}
	return fmt.Errorf("device with name %s doesn't exist", name)
}
func (a *automation) turnOffDevice(name string) error {
	dc, exists := a.config.DeviceControlCache["name"]
	if exists {
		err := a.service.Pubsub.PublishStr(dc.Channel, dc.Off)
		return err
	}
	return fmt.Errorf("device with name %s doesn't exist", name)
}
func (a *automation) toggleDevice(name string) error {
	dc, exists := a.config.DeviceControlCache["name"]
	if exists {
		err := a.service.Pubsub.PublishStr(dc.Channel, dc.Toggle)
		return err
	}
	return fmt.Errorf("device with name %s doesn't exist", name)
}

func (a *automation) presenceDetection() {
	peopleWhereHome := a.peopleAreHome()
	a.presence.determineIfPeopleAreHome(a.now)
	peopleAreHome := a.peopleAreHome()

	if !peopleWhereHome && peopleAreHome {
		// Depending on time-of-day
		switch a.timeofday {
		case "lunch":
			a.sendNotification("Turning on kitchen lights since it is noon and someone came home")
			a.turnOnDevice(config.KitchenLights)
			break
		case "evening":
			a.sendNotification("Turning on kitchen and livingroom lights since it is evening and someone came home")
			a.turnOnDevice(config.KitchenLights)
			a.turnOnDevice(config.LivingroomLightMain)
			a.turnOnDevice(config.LivingroomLightStand)
			a.turnOnDevice(config.LivingroomLightChandelier)
			break
		case "bedtime":
			a.sendNotification("Turning on kitchen and livingroom lights since it is bedtime and someone came home")
			a.turnOnDevice(config.KitchenLights)
			a.turnOnDevice(config.LivingroomLightMain)
			a.turnOnDevice(config.LivingroomLightStand)
			a.turnOnDevice(config.LivingroomLightChandelier)
			break
		case "sleeptime":
			a.sendNotification("Turning on kitchen and livingroom lights since it is sleeptime and someone came home")
			a.turnOnDevice(config.KitchenLights)
			a.turnOnDevice(config.LivingroomLightMain)
			a.turnOnDevice(config.LivingroomLightStand)
			a.turnOnDevice(config.LivingroomLightChandelier)
			break
		}

	} else if peopleWhereHome && !peopleAreHome {
		// Turn off everything
		a.turnOffEverything()
	}
}

func (a *automation) handleEvent(channel string, state *config.SensorState) {
	sensortype := state.Type
	sensorname := state.Name
	if sensortype == "sensor" {
		a.sensors[sensorname] = state.GetValueAttr("state", "")
		if sensorname == "timeofday" {
			a.handleTimeOfDay(a.sensors[sensorname])
		}
	} else if sensortype == "switch" {
		if sensorname == config.SophiaRoomSwitch || sensorname == config.BedroomSwitch {
			a.handleSwitch(sensorname, state)
		} else if sensorname == config.KitchenMotionSensor || sensorname == config.LivingroomMotionSensor || sensorname == config.BedroomMotionSensor {
			a.handleMotionSensor(sensorname, state)
		} else if sensorname == config.FrontdoorMagnetSensor {
			a.handleMagnetSensor(sensorname, state)
		}
	} else if sensortype == "presence" {
		a.handlePresence(state)
	}
}
func wakeUpParentsForSophiaAndJennifer(ta *timedBasedAction, a *automation) {
	a.sendNotification("Waking up Parents for Sophia & Jennifer")
	a.turnOnDevice(config.BedroomLightStand)
}
func wakeUpParentsForJennifer(ta *timedBasedAction, a *automation) {
	a.sendNotification("Waking up Parents for Jennifer")
	a.turnOnDevice(config.BedroomLightStand)
}
func wakeUpParentsForSophia(ta *timedBasedAction, a *automation) {
	a.sendNotification("Waking up Parents for Sophia")
	a.turnOnDevice(config.BedroomLightStand)
}
func wakeUpParentsForWork(ta *timedBasedAction, a *automation) {
	a.sendNotification("Waking up Parents")
	a.turnOnDevice(config.BedroomLightStand)
}
func wakeUpJennifer(ta *timedBasedAction, a *automation) {
	a.sendNotification("Waking up Jennifer")
	a.turnOnDevice(config.JenniferRoomLightMain)
}
func wakeUpSophia(ta *timedBasedAction, a *automation) {
	a.sendNotification("Waking up Sophia")
	a.turnOnDevice(config.SophiaRoomLightStand)
}
func turnOnFrontdoorHallLight(ta *timedBasedAction, a *automation) {
	a.turnOnDevice(config.FrontdoorHallLight)
}

// HandleTimeOfDay deals with time-of-day transitions
func (a *automation) handleTimeOfDay(to string) {
	if to != a.timeofday {
		a.timeofday = to
		switch to {
		case "breakfast":
			jenniferHasSchool := a.sensorHasValue("jennifer", "school")
			sophiaHasSchool := a.sensorHasValue("sophia", "school")
			parentsHaveToWork := a.sensorHasValue("parents", "work")
			if parentsHaveToWork {
				a.setRealTimeAction("Waking up Parents", 8, 0, wakeUpParentsForWork)
			}
			if jenniferHasSchool && sophiaHasSchool {
				a.setRealTimeAction("Waking up Parents", 6, 20, wakeUpParentsForSophiaAndJennifer)
			} else if jenniferHasSchool {
				a.setRealTimeAction("Waking up Parents", 6, 20, wakeUpParentsForJennifer)
			} else if sophiaHasSchool {
				a.setRealTimeAction("Waking up Parents", 6, 20, wakeUpParentsForSophia)
			}
			if jenniferHasSchool {
				a.setRealTimeAction("Waking up Jennifer", 6, 30, wakeUpJennifer)
				a.setRealTimeAction("Turn on Hall Light", 7, 11, turnOnFrontdoorHallLight)
			}
			if sophiaHasSchool {
				a.setRealTimeAction("Waking up Sophia", 6, 30, wakeUpSophia)
				a.setRealTimeAction("Turn on Hall Light", 7, 11, turnOnFrontdoorHallLight)
			}
		case "morning":
			a.sendNotification("Turning off lights since it is morning")
			a.turnOffDevice(config.KitchenLights)
			a.turnOffDevice(config.LivingroomLightStand)
		case "lunch":
			if a.peopleAreHome() {
				a.sendNotification("Turning on lights since it is noon and someone is home")
				a.turnOnDevice(config.KitchenLights)
			}
		case "afternoon":
			a.sendNotification("Turning off lights since it is afternoon")
			a.turnOffDevice(config.KitchenLights)
		case "evening":
			a.turnOnDevice(config.KitchenLights)
			a.turnOnDevice(config.LivingroomLightStand)
		case "bedtime":
			if a.peopleAreHome() {
				if a.sensorHasValue("jennifer", "school") {
					a.turnOnDevice(config.JenniferRoomLightMain)
				}
				if a.sensorHasValue("sophia", "school") {
					a.turnOnDevice(config.SophiaRoomLightMain)
					a.turnOnDevice(config.SophiaRoomLightStand)
				}
			}
		case "sleeptime":
			if a.peopleAreHome() {
				a.turnOnDevice(config.BedroomLightStand)
				a.turnOnDevice(config.BedroomLightMain)
				if a.sensorHasValue("jennifer", "school") {
					a.turnOffDevice(config.JenniferRoomLightMain)
				}
				if a.sensorHasValue("sophia", "school") {
					a.turnOnDevice(config.SophiaRoomLightMain)
				}
			}
		case "night":
			if a.sensorHasValue("jennifer", "school") || a.sensorHasValue("sophia", "school") {
				a.turnOffDevice(config.BedroomLightMain)
			}
			a.turnOffDevice(config.KitchenLights)
			a.turnOffDevice(config.LivingroomLightMain)
			a.turnOffDevice(config.LivingroomLightStand)
			a.turnOffDevice(config.JenniferRoomLightMain)
			a.turnOffDevice(config.SophiaRoomLightStand)
			a.turnOffDevice(config.SophiaRoomLightMain)
			a.turnOffDevice(config.FrontdoorHallLight)
		}
	}
}

// HandleSwitch deals with switches being pressed
func (a *automation) handleSwitch(name string, state *config.SensorState) {
	if name == config.BedroomSwitch {
		value := state.GetValueAttr("click", "")
		if value == config.WirelessSwitchDoubleClick {
			a.presence.reportCausation(a.now)
			a.toggleDevice(config.BedroomLightMain)
		}
		if value == config.WirelessSwitchSingleClick {
			a.presence.reportCausation(a.now)
			a.toggleDevice(config.BedroomLightStand)
		}
		if value == config.WirelessSwitchLongPress {
			a.presence.reportCausation(a.now)
			a.toggleDevice(config.BedroomPowerPlug)
		}
	} else if name == config.SophiaRoomSwitch {
		value := state.GetValueAttr("click", "")
		if value == config.WirelessSwitchSingleClick {
			a.presence.reportCausation(a.now)
			a.toggleDevice(config.SophiaRoomLightStand)
		} else if value == config.WirelessSwitchDoubleClick {
			a.presence.reportCausation(a.now)
			a.toggleDevice(config.SophiaRoomLightMain)
		}
	}
}

// HandleMotionSensor deals with motion detected
func (a *automation) handleMotionSensor(name string, state *config.SensorState) {
	now := time.Now()
	if name == config.KitchenMotionSensor || name == config.LivingroomMotionSensor {
		value := state.GetValueAttr("motion", "")
		if value == "on" {
			a.presence.reportCausation(now)
			a.lastseenMotionInHouse = now // Update the time we last detected motion
			if name == config.KitchenMotionSensor {
				a.lastseenMotionInKitchenArea = now
				a.setDelayTimeAction("Turnoff front door hall light", 4*time.Minute, func(ta *timedBasedAction, a *automation) { a.turnOffDevice(config.FrontdoorHallLight) })
			}
			if a.timeofday == "breakfast" {
				a.turnOnDevice(config.KitchenLights)
				a.turnOnDevice(config.LivingroomLightStand)
			}

			if a.timeofday == "night" {
				if name == config.KitchenMotionSensor {
					a.sendNotification("Turning on kitchen and livingroom lights since it is night and there is movement in the kitchen area")
				} else if name == config.LivingroomMotionSensor {
					a.sendNotification("Turning on kitchen and livingroom lights since it is night and there is movement in the livingroom area")
				}
				a.turnOnDevice(config.KitchenLights)
				a.turnOnDevice(config.LivingroomLightStand)
				a.setDelayTimeAction(config.KitchenLights, 5*time.Minute, func(ta *timedBasedAction, a *automation) { a.turnOffDevice(config.KitchenLights) })
				a.setDelayTimeAction(config.KitchenLights, 5*time.Minute, func(ta *timedBasedAction, a *automation) { a.turnOffDevice(config.KitchenLights) })
			}
		}
	} else if name == config.BedroomMotionSensor {
		value := state.GetValueAttr("motion", "")
		lastseenDuration := now.Sub(a.lastseenMotionInBedroom)
		if value == "on" {
			a.presence.reportCausation(now)
			a.lastseenMotionInHouse = now // Update the time we last detected motion
			a.lastseenMotionInBedroom = now

			if a.timeofday == "evening" || a.timeofday == "bedtime" {
				if lastseenDuration > (time.Duration(15) * time.Minute) {
					a.turnOnDevice(config.BedroomLightStand)
					a.turnOnDevice(config.BedroomLightMain)
				}
			}
		} else if value == "off" {
			if a.timeofday != "night" && a.timeofday != "sleeptime" {
				if lastseenDuration > (time.Duration(30) * time.Minute) {
					a.turnOffDevice(config.BedroomLightMain)
				}
			}
		}
	}
}

func (a *automation) handleMagnetSensor(name string, state *config.SensorState) {
	if name == config.FrontdoorMagnetSensor {
		value := state.GetValueAttr("state", "")
		if value == "open" {
			a.sendNotification("Front door opened")
			a.turnOnDevice(config.FrontdoorHallLight)
			a.setDelayTimeAction("Turnoff front door hall light", 10*time.Minute, func(ta *timedBasedAction, a *automation) { a.turnOffDevice(config.FrontdoorHallLight) })
			a.presence.frontDoorOpenClosed()
			a.lastseenMotionInKitchenArea = time.Now()
		} else if value == "close" {
			a.sendNotification("Front door closed")
			a.setDelayTimeAction("Turnoff front door hall light", 5*time.Minute, func(ta *timedBasedAction, a *automation) { a.turnOffDevice(config.FrontdoorHallLight) })
			a.presence.frontDoorOpenClosed()
			a.lastseenMotionInKitchenArea = time.Now()
		}
	}
}

func (a *automation) sendNotification(message string) {
	a.service.Pubsub.PublishStr("shout/message/", message)
}

func (a *automation) updatePresence(name string, presence bool) (current bool, previous bool) {
	var exists bool
	previous, exists = a.presence.presence[name]
	if !exists {
		previous = false
		a.presence.presence[name] = presence
	}
	return presence, previous
}

func (a *automation) handlePresence(state *config.SensorState) {
	if state.StringAttrs != nil {
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
	}
}

func (a *automation) turnOffEverything() {
	a.turnOffDevice(config.KitchenLights)
	a.turnOffDevice(config.LivingroomLightStand)
	a.turnOffDevice(config.LivingroomLightMain)
	a.turnOffDevice(config.LivingroomLightChandelier)
	a.turnOffDevice(config.BedroomLightStand)
	a.turnOffDevice(config.BedroomLightMain)
	a.turnOffDevice(config.JenniferRoomLightMain)
	a.turnOffDevice(config.SophiaRoomLightStand)
	a.turnOffDevice(config.SophiaRoomLightMain)
	a.turnOffDevice(config.FrontdoorHallLight)

	a.turnOffDevice(config.BedroomSamsungTV)
	a.turnOffDevice(config.LivingroomBraviaTV)
}

// UpdateTimedActions will tick all actions and remove any that have ended
func (a *automation) updateTimedActions() {
	var toremove []string
	for name, ta := range a.timedActions {
		if ta.tick(a.now) == true {
			ta.Action(ta, a)
			toremove = append(toremove, name)
		}
	}
	for _, name := range toremove {
		delete(a.timedActions, name)
	}
}

// UpdateTimedActions will tick all actions and remove any that have ended
func (a *automation) updateMotionBasedActions(sensors map[string]bool) {
	var toremove []string
	for name, ta := range a.motionBasedActions {
		if ta.tick(a.now, sensors) == actionTriggered {
			ta.Action(ta, a)
			if ta.RemoveWhenTriggered {
				toremove = append(toremove, name)
			} else {
				ta.reset(a.now)
			}
		}
	}
	for _, name := range toremove {
		delete(a.motionBasedActions, name)
	}
}

// actionDelegate is the action delegate
type actionDelegate func(ta *timedBasedAction, a *automation)

// timedBasedAction holds the 'when' and 'action' necessary to trigger
type timedBasedAction struct {
	Name   string
	When   time.Time
	Action actionDelegate
}

// setRealTimeAction sets an action that will trigger at a specific time
func (a *automation) setRealTimeAction(name string, hour int, minute int, ad actionDelegate) {
	action, exists := a.timedActions[name]
	if !exists {
		when := time.Date(a.now.Year(), a.now.Month(), a.now.Day(), hour, minute, 0, 0, a.now.Location())
		action = &timedBasedAction{Name: name, When: when, Action: ad}
		a.timedActions[name] = action
	}
}

// setDelayTimeAction sets an action that trigger after 'duration'
func (a *automation) setDelayTimeAction(name string, duration time.Duration, actiondelegate actionDelegate) {
	action, exists := a.timedActions[name]
	if !exists {
		action = &timedBasedAction{Name: name, When: time.Now().Add(duration), Action: actiondelegate}
		a.timedActions[name] = action
	} else {
		action.When = time.Now().Add(duration)
	}
}

// tick returns true when the action has ended, false otherwise
func (ta *timedBasedAction) tick(now time.Time) bool {
	return now.After(ta.When)
}

// Examples:
// A) While there is motion every 10 minutes keep the lights ON, once there is not turn the
//    lights OFF and auto-delete
// B) Once there is no motion in the kitchen for 15 minutes turn OFF the lights and auto-delete
// C) For the next duration once there is motion turn ON the lights and auto-delete
//
type motionDelegate func(ta *motionBasedAction, a *automation)
type motionType int

const (
	whileMotionEveryPeriodThen motionType = iota
	whileNoMotionForCertainDurationThen
	whileNoMotionUntilMotion
)

type motionBasedAction struct {
	Name                string
	Sensors             []string
	From                time.Time
	Duration            time.Duration
	Count               int
	Action              motionDelegate
	Type                motionType
	RemoveWhenTriggered bool // When this action fired should we remove ourselves?
}

// setNoMotionAction sets an action that will trigger at a specific time
func (a *automation) setMotionAction(name string, motiontype motionType, sensors []string, hour int, minute int, actiondelegate motionDelegate, removeWhenTriggered bool) {
	action, exists := a.motionBasedActions[name]
	duration := time.Date(a.now.Year(), a.now.Month(), a.now.Day(), hour, minute, 0, 0, a.now.Location()).Sub(a.now)
	if !exists {
		action = &motionBasedAction{Name: name}
		action.Sensors = sensors
		action.From = a.now
		action.Duration = duration
		action.Count = 0
		action.Action = actiondelegate
		action.Type = motiontype
		action.RemoveWhenTriggered = removeWhenTriggered
		a.motionBasedActions[name] = action
	} else {
		action.Sensors = sensors
		action.From = a.now
		action.Duration = duration
		action.Count = 0
		action.Action = actiondelegate
		action.Type = motiontype
		action.RemoveWhenTriggered = removeWhenTriggered
	}
}

// resetMotionAction sets an action that trigger after 'motion'
func (ta *motionBasedAction) reset(now time.Time) {
	ta.From = now
	ta.Count = 0
}

type tickResult int

const (
	actionTriggered tickResult = iota
	actionEvaluating
	actionFailed
)

// tick returns the result when the action is still evaluating, has been triggered or has failed
func (ta *motionBasedAction) tick(now time.Time, sensors map[string]bool) tickResult {
	motion := false
	for _, name := range ta.Sensors {
		state, exists := sensors[name]
		if exists {
			motion = motion || state
		}
	}
	if ta.Type == whileMotionEveryPeriodThen {
		until := ta.From.Add(ta.Duration)
		if now.After(until) {
			if ta.Count > 0 {
				ta.From = now
				ta.Count = 0
				return actionEvaluating
			}
			return actionTriggered
		} else if motion {
			ta.Count++
		}
		return actionEvaluating
	} else if ta.Type == whileNoMotionForCertainDurationThen {
		when := ta.From.Add(ta.Duration)
		if now.After(when) {
			return actionTriggered
		}
		if !motion {
			return actionEvaluating
		}
	} else if ta.Type == whileNoMotionUntilMotion {
		if motion {
			return actionTriggered
		}
		return actionEvaluating
	}
	return actionFailed
}
