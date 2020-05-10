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
	auto := new()

	module := "automation"
	logger := logpkg.New(module)
	logger.AddEntry("emitter")
	logger.AddEntry(module)

	for {
		auto.pubsub = pubsub.New(config.PubSubCfg)
		err := auto.pubsub.Connect(module, []string{}, []string{"config/automation/", "config/request/"})

		if err == nil {
			logger.LogInfo("emitter", "connected")
			connected := true
			for connected {
				select {
				case msg := <-auto.pubsub.InMsgs:
					topic := msg.Topic()
					if topic == "config/automation/" {
						// Register used channels and subscribe to channels we are interested in
						config, err := config.AutomationConfigFromJSON(msg.Payload())
						if err == nil {
							auto.config = config
							// Register used channels
							for _, ss := range auto.config.ChannelsToRegister {
								if err = auto.pubsub.Register(ss); err != nil {
									logger.LogError(module, err.Error())
								}
							}
							// Subscribe channels
							for _, ss := range auto.config.SubChannels {
								if err = auto.pubsub.Subscribe(ss); err != nil {
									logger.LogError(module, err.Error())
								}
							}
						} else {
							logger.LogError(module, err.Error())
						}
					} else if topic == "client/disconnected/" {
						connected = false
						logger.LogInfo("emitter", "disconnected")
					} else if strings.HasPrefix(topic, "state/") {
						state, err := config.SensorStateFromJSON(msg.Payload())
						if err == nil {
							auto.handleEvent(topic, state)
						} else {
							logger.LogError(module, err.Error())
						}
					}
				case <-time.After(time.Second * 5):
					if auto.config != nil {
						auto.now = time.Now()
						auto.presenceDetection()
						auto.updateTimedActions()
					}

				case <-time.After(time.Minute * 1): // Try and request our configuration
					if auto.config == nil {
						auto.pubsub.Publish("config/request/", "automation")
					}
				}
			}
		}
		if err != nil {
			logger.LogError(module, err.Error())
		}

		// Wait for 5 seconds before retrying
		time.Sleep(5 * time.Second)
	}
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
	h.detectionEvalDuration = time.Minute * 10
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
			// Look at the Wifi presence. (Should we actually do this, is motion not enough ?)
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

// report() should be called when a causation is detected like a button press, light switch press or movement
func (h *homePresence) report(now time.Time) {
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
	pubsub                      *pubsub.Context
	config                      *config.AutomationConfig
	sensors                     map[string]string
	timeofday                   string
	now                         time.Time
	lastseenMotionInHouse       time.Time
	lastseenMotionInKitchenArea time.Time
	lastseenMotionInBedroom     time.Time
	timedActions                map[string]*timedAction
	motionBasedActions          map[string]*motionBasedAction
	presence                    *homePresence
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
		err := a.pubsub.Publish(dc.Channel, dc.On)
		return err
	}
	return fmt.Errorf("device with name %s doesn't exist", name)
}
func (a *automation) turnOffDevice(name string) error {
	dc, exists := a.config.DeviceControlCache["name"]
	if exists {
		err := a.pubsub.Publish(dc.Channel, dc.Off)
		return err
	}
	return fmt.Errorf("device with name %s doesn't exist", name)
}
func (a *automation) toggleDevice(name string) error {
	dc, exists := a.config.DeviceControlCache["name"]
	if exists {
		err := a.pubsub.Publish(dc.Channel, dc.Toggle)
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
		case "evening":
			a.sendNotification("Turning on kitchen and livingroom lights since it is evening and someone came home")
			a.turnOnDevice(config.KitchenLights)
			a.turnOnDevice(config.LivingroomLights)
		case "bedtime":
			a.sendNotification("Turning on kitchen and livingroom lights since it is bedtime and someone came home")
			a.turnOnDevice(config.KitchenLights)
			a.turnOnDevice(config.LivingroomLights)
		case "sleeptime":
			a.sendNotification("Turning on kitchen and livingroom lights since it is sleeptime and someone came home")
			a.turnOnDevice(config.KitchenLights)
			a.turnOnDevice(config.LivingroomLights)
		}

	} else if peopleWhereHome && !peopleAreHome {
		// Turn off everything
		a.turnOffEverything()
	}
}

func (a *automation) handleEvent(channel string, state *config.SensorState) {
	sensortype := ""
	sensorname := ""
	parts := strings.Split(channel, "/")
	if len(parts) >= 2 {
		sensortype = parts[1]
		if len(parts) == 3 {
			sensorname = parts[2]
		}
	}

	if sensorname != "" && sensortype == "sensor" {
		a.sensors[sensorname] = state.GetValueAttr(sensorname, "")
		if sensorname == "timeofday" {
			a.handleTimeOfDay(a.sensors[sensorname])
		}
	} else if sensortype == "xiaomi" {
		name := state.Name
		if name == config.SophiaRoomSwitch || name == config.BedroomSwitch {
			a.handleSwitch(name, state)
		} else if name == config.KitchenMotionSensor || name == config.LivingroomMotionSensor || name == config.BedroomMotionSensor {
			a.handleMotionSensor(name, state)
		} else if name == config.FrontdoorMagnetSensor {
			a.handleMagnetSensor(name, state)
		}
	} else if sensortype == "presence" {
		a.handlePresence(state)
	}
}

func wakeUpParentsForJennifer(ta *timedAction, a *automation) {
	a.sendNotification("Waking up Parents for Jennifer")
	a.turnOnDevice(config.BedroomLights)
}
func wakeUpParentsForSophia(ta *timedAction, a *automation) {
	a.sendNotification("Waking up Parents for Sophia")
	a.turnOnDevice(config.BedroomLights)
}
func wakeUpParentsForWork(ta *timedAction, a *automation) {
	a.sendNotification("Waking up Parents")
	a.turnOnDevice(config.BedroomLights)
}
func wakeUpJennifer(ta *timedAction, a *automation) {
	a.sendNotification("Waking up Jennifer")
	a.turnOnDevice(config.JenniferRoomLights)
}
func wakeUpSophia(ta *timedAction, a *automation) {
	a.sendNotification("Waking up Sophia")
	a.turnOnDevice(config.SophiaRoomLights)
}
func turnOnFrontdoorHallLight(ta *timedAction, a *automation) {
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
			if jenniferHasSchool {
				a.setRealTimeAction("Waking up Parents for Jennifer", 6, 20, wakeUpParentsForJennifer)
			} else if sophiaHasSchool {
				a.setRealTimeAction("Waking up Parents for Sophia", 7, 0, wakeUpParentsForSophia)
			} else if parentsHaveToWork {
				a.setRealTimeAction("Waking up Parents", 7, 30, wakeUpParentsForWork)
			}
			if jenniferHasSchool {
				a.setRealTimeAction("Waking up Jennifer", 6, 30, wakeUpJennifer)
				a.setRealTimeAction("Turn on Hall Light", 7, 11, turnOnFrontdoorHallLight)
			}
			if sophiaHasSchool {
				a.setRealTimeAction("Waking up Sophia", 7, 10, wakeUpSophia)
			}
		case "morning":
			a.sendNotification("Turning off lights since it is morning")
			a.turnOffDevice(config.KitchenLights)
			a.turnOffDevice(config.LivingroomLights)
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
			a.turnOnDevice(config.LivingroomLights)
		case "bedtime":
			if a.peopleAreHome() {
				if a.sensorHasValue("jennifer", "school") {
					a.turnOnDevice(config.JenniferRoomLights)
				}
				if a.sensorHasValue("sophia", "school") {
					a.turnOnDevice(config.SophiaRoomLights)
				}
			}
		case "sleeptime":
			if a.peopleAreHome() {
				a.turnOnDevice(config.BedroomLights)
				if a.sensorHasValue("jennifer", "school") {
					a.turnOffDevice(config.JenniferRoomLights)
				}
				if a.sensorHasValue("sophia", "school") {
					a.turnOffDevice(config.SophiaRoomLights)
				}
			}
		case "night":
			if a.sensorHasValue("jennifer", "school") {
				a.turnOffDevice(config.BedroomLights)
			}
			a.turnOffDevice(config.KitchenLights)
			a.turnOffDevice(config.LivingroomLights)
			a.turnOffDevice(config.JenniferRoomLights)
			a.turnOffDevice(config.SophiaRoomLights)
			a.turnOffDevice(config.FrontdoorHallLight)
		}
	}
}

// HandleSwitch deals with switches being pressed
func (a *automation) handleSwitch(name string, state *config.SensorState) {
	if name == config.BedroomSwitch {
		value := state.GetValueAttr("click", "")
		if value == config.WirelessSwitchDoubleClick {
			a.presence.report(a.now)
			a.toggleDevice(config.BedroomLights)
		}
		if value == config.WirelessSwitchSingleClick {
			a.presence.report(a.now)
			a.turnOffDevice(config.BedroomCeilingLightSwitch)
			a.turnOffDevice(config.BedroomChandelierLightSwitch)
		}
		if value == config.WirelessSwitchLongPress {
			a.presence.report(a.now)
			a.toggleDevice(config.BedroomPowerPlug)
		}
	} else if name == config.SophiaRoomSwitch {
		value := state.GetValueAttr("click", "")
		if value == config.WirelessSwitchSingleClick {
			a.presence.report(a.now)
			a.toggleDevice(config.SophiaRoomLights)
		}
	}
}

// HandleMotionSensor deals with motion detected
func (a *automation) handleMotionSensor(name string, state *config.SensorState) {
	now := time.Now()
	if name == config.KitchenMotionSensor || name == config.LivingroomMotionSensor {
		value := state.GetValueAttr("motion", "")
		if value == "on" {
			a.presence.report(now)
			a.lastseenMotionInHouse = now // Update the time we last detected motion
			if name == config.KitchenMotionSensor {
				a.lastseenMotionInKitchenArea = now
				a.setDelayTimeAction("Turnoff front door hall light", 4*time.Minute, func(ta *timedAction, a *automation) { a.turnOffDevice("Front door hall light") })
			}
			if a.timeofday == "breakfast" {
				a.turnOnDevice(config.KitchenLights)
				a.turnOnDevice(config.LivingroomLights)
			}

			if a.timeofday == "night" {
				if name == config.KitchenMotionSensor {
					a.sendNotification("Turning on kitchen and livingroom lights since it is night and there is movement in the kitchen area")
				} else if name == config.LivingroomMotionSensor {
					a.sendNotification("Turning on kitchen and livingroom lights since it is night and there is movement in the livingroom area")
				}
				a.turnOnDevice(config.KitchenLights)
				a.turnOnDevice(config.LivingroomLights)
				a.setDelayTimeAction(config.KitchenLights, 5*time.Minute, func(ta *timedAction, a *automation) { a.turnOffDevice(config.KitchenLights) })
				a.setDelayTimeAction(config.KitchenLights, 5*time.Minute, func(ta *timedAction, a *automation) { a.turnOffDevice(config.KitchenLights) })
			}
		}
	} else if name == config.BedroomMotionSensor {
		value := state.GetValueAttr("motion", "")
		lastseenDuration := now.Sub(a.lastseenMotionInBedroom)
		if value == "on" {
			a.presence.report(now)
			a.lastseenMotionInHouse = now // Update the time we last detected motion
			a.lastseenMotionInBedroom = now

			if a.timeofday == "evening" || a.timeofday == "bedtime" {
				if lastseenDuration > (time.Duration(15) * time.Minute) {
					a.turnOnDevice(config.BedroomLights)
					a.turnOnDevice(config.BedroomChandelierLightSwitch)
				}
			}
		} else if value == "off" {
			if a.timeofday != "night" && a.timeofday != "sleeptime" {
				if lastseenDuration > (time.Duration(30) * time.Minute) {
					a.turnOffDevice(config.BedroomLights)
					a.turnOffDevice(config.BedroomChandelierLightSwitch)
					a.turnOffDevice(config.BedroomCeilingLightSwitch)
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
			a.setDelayTimeAction("Turnoff front door hall light", 10*time.Minute, func(ta *timedAction, a *automation) { a.turnOffDevice("Front door hall light") })

			a.lastseenMotionInKitchenArea = time.Now()
		} else if value == "close" {
			a.sendNotification("Front door closed")
			a.setDelayTimeAction("Turnoff front door hall light", 5*time.Minute, func(ta *timedAction, a *automation) { a.turnOffDevice("Front door hall light") })
			a.presence.frontDoorOpenClosed()
		}
	}
}

func (a *automation) sendNotification(message string) {
	a.pubsub.Publish("shout/message/", message)
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
	a.turnOffDevice(config.LivingroomLights)
	a.turnOffDevice(config.BedroomLights)
	a.turnOffDevice(config.JenniferRoomLights)
	a.turnOffDevice(config.SophiaRoomLights)
	a.turnOffDevice(config.FrontdoorHallLight)
	a.turnOffDevice(config.BedroomPowerPlug)
	a.turnOffDevice(config.BedroomChandelierLightSwitch)
	a.turnOffDevice(config.BedroomCeilingLightSwitch)

	a.turnOffDevice(config.BedroomSamsungTV)
	a.turnOffDevice(config.LivingroomBraviaTV)
}

// actionDelegate is the action delegate
type actionDelegate func(ta *timedAction, a *automation)

// timedAction holds the 'when' and 'action' necessary to trigger
type timedAction struct {
	Name   string
	When   time.Time
	Action actionDelegate
}

// setRealTimeAction sets an action that will trigger at a specific time
func (a *automation) setRealTimeAction(name string, hour int, minute int, ad actionDelegate) {
	action, exists := a.timedActions[name]
	if !exists {
		when := time.Date(a.now.Year(), a.now.Month(), a.now.Day(), hour, minute, 0, 0, a.now.Location())
		action = &timedAction{Name: name, When: when, Action: ad}
		a.timedActions[name] = action
	}
}

// setDelayTimeAction sets an action that trigger after 'duration'
func (a *automation) setDelayTimeAction(name string, duration time.Duration, actiondelegate actionDelegate) {
	action, exists := a.timedActions[name]
	if !exists {
		action = &timedAction{Name: name, When: time.Now().Add(duration), Action: actiondelegate}
		a.timedActions[name] = action
	} else {
		action.When = time.Now().Add(duration)
	}
}

// tick returns true when the action has ended, false otherwise
func (ta *timedAction) tick(now time.Time) bool {
	return now.After(ta.When)
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
		delete(a.motionBasedActions, name)
	}
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
