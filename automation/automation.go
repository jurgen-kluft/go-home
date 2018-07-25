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

	module := "automation"
	logger := logpkg.New(module)
	logger.AddEntry("emitter")
	logger.AddEntry(module)

	for {
		auto.pubsub = pubsub.New(config.EmitterSecrets["host"])
		err := auto.pubsub.Connect(module, []string{}, []string{"config/automation/"})

		if err == nil {
			logger.LogInfo("emitter", "connected")
			connected := true
			for connected {
				select {
				case msg := <-auto.pubsub.InMsgs:
					topic := msg.Topic()
					if topic == "config/automation/" {
						// Register used channels and subscribe to channels we are interested in
						config, err := config.AutomationConfigFromJSON(string(msg.Payload()))
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
						state, err := config.SensorStateFromJSON(string(msg.Payload()))
						if err == nil {
							auto.HandleEvent(topic, state)
						} else {
							logger.LogError(module, err.Error())
						}
					}
				case <-time.After(time.Second * 5):
					if auto.config != nil {
						auto.now = time.Now()
						auto.UpdateTimedActions()
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

type HomePresence struct {
	occupied        bool
	detection       bool
	detectionResult bool
	detectionStamp  time.Time
	detectionWindow time.Duration
	presence        map[string]bool
}

func NewPresence() *HomePresence {
	h := &HomePresence{}
	h.occupied = true
	h.detection = true
	h.detectionResult = false
	h.detectionStamp = time.Now()
	h.detectionWindow = time.Minute * 15
	h.presence = map[string]bool{}
	return h
}

// reset() should be called when the front-door is opened->closed because this
// can indicate people have left the house.
func (h *HomePresence) reset() {
	h.detection = true
	h.detectionResult = false
	h.detectionStamp = time.Now()
}

func (h *HomePresence) detect(now time.Time) {
	if h.detection {
		if now.Sub(h.detectionStamp) > h.detectionWindow {
			h.occupied = h.detectionResult
			h.detection = false
		}
	} else {
		// After the detection window (door was closed + N minutes) we can look at the wifi-presence again.
		if h.occupied == false {
			for _, prsnc := range h.presence {
				h.occupied = h.occupied || prsnc
			}
		}
	}
}

// report() should be called when a causation is detected like a button press or light switch press
func (h *HomePresence) report(now time.Time) {
	if h.detection {
		if now.Sub(h.detectionStamp) < h.detectionWindow {
			h.detectionResult = true
		}
	} else {
		// Any detected presence after the detection window means that people are home
		// This kind of presence is a button, light switch or plug press
		h.occupied = true
	}
}

type Automation struct {
	pubsub                      *pubsub.Context
	config                      *config.AutomationConfig
	sensors                     map[string]string
	timeofday                   string
	now                         time.Time
	lastseenMotionInHouse       time.Time
	lastseenMotionInKitchenArea time.Time
	lastseenMotionInBedroom     time.Time
	timedActions                map[string]*TimedAction
	presence                    *HomePresence
}

func New() *Automation {
	auto := &Automation{}
	auto.sensors = map[string]string{}
	auto.presence = NewPresence()
	return auto
}

func (a *Automation) FamilyIsHome() bool {
	return a.presence.occupied
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

func wakeUpParentsForJennifer(ta *TimedAction, a *Automation) {
	a.sendNotification("Waking up Parents for Jennifer")
	a.TurnOnDevice(config.BedroomLights)
}
func wakeUpParentsForSophia(ta *TimedAction, a *Automation) {
	a.sendNotification("Waking up Parents for Sophia")
	a.TurnOnDevice(config.BedroomLights)
}
func wakeUpParentsForWork(ta *TimedAction, a *Automation) {
	a.sendNotification("Waking up Parents")
	a.TurnOnDevice(config.BedroomLights)
}
func wakeUpJennifer(ta *TimedAction, a *Automation) {
	a.sendNotification("Waking up Jennifer")
	a.TurnOnDevice(config.JenniferRoomLights)
}
func wakeUpSophia(ta *TimedAction, a *Automation) {
	a.sendNotification("Waking up Sophia")
	a.TurnOnDevice(config.SophiaRoomLights)
}
func turnOnFrontdoorHallLight(ta *TimedAction, a *Automation) {
	a.TurnOnDevice(config.FrontdoorHallLight)
}

// HandleTimeOfDay deals with time-of-day transitions
func (a *Automation) HandleTimeOfDay(to string) {
	if to != a.timeofday {
		a.timeofday = to
		switch to {
		case "breakfast":
			jenniferHasSchool := a.SensorHasValue("jennifer", "school")
			sophiaHasSchool := a.SensorHasValue("sophia", "school")
			parentsHaveToWork := a.SensorHasValue("parents", "work")
			if jenniferHasSchool {
				a.SetRealTimeAction("Waking up Parents for Jennifer", 6, 20, wakeUpParentsForJennifer)
			} else if sophiaHasSchool {
				a.SetRealTimeAction("Waking up Parents for Sophia", 7, 0, wakeUpParentsForSophia)
			} else if parentsHaveToWork {
				a.SetRealTimeAction("Waking up Parents", 7, 30, wakeUpParentsForWork)
			}
			if jenniferHasSchool {
				a.SetRealTimeAction("Waking up Jennifer", 6, 30, wakeUpJennifer)
				a.SetRealTimeAction("Turn on Hall Light", 7, 11, turnOnFrontdoorHallLight)
			}
			if sophiaHasSchool {
				a.SetRealTimeAction("Waking up Sophia", 7, 10, wakeUpSophia)
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

// HandleSwitch deals with switches being pressed
func (a *Automation) HandleSwitch(name string, state *config.SensorState) {
	if name == config.BedroomSwitch {
		value := state.GetValueAttr("click", "")
		if value == config.WirelessSwitchDoubleClick {
			a.presence.report(a.now)
			a.ToggleDevice(config.BedroomLights)
		}
		if value == config.WirelessSwitchSingleClick {
			a.presence.report(a.now)
			a.TurnOffDevice(config.BedroomCeilingLightSwitch)
			a.TurnOffDevice(config.BedroomChandelierLightSwitch)
		}
		if value == config.WirelessSwitchLongPress {
			a.presence.report(a.now)
			a.ToggleDevice(config.BedroomPowerPlug)
		}
	} else if name == config.SophiaRoomSwitch {
		value := state.GetValueAttr("click", "")
		if value == config.WirelessSwitchSingleClick {
			a.presence.report(a.now)
			a.ToggleDevice(config.SophiaRoomLights)
		}
	}
}

// HandleMotionSensor deals with motion detected
func (a *Automation) HandleMotionSensor(name string, state *config.SensorState) {
	now := time.Now()
	if name == config.KitchenMotionSensor || name == config.LivingroomMotionSensor {
		value := state.GetValueAttr("motion", "")
		if value == "on" {
			a.presence.report(now)
			a.lastseenMotionInHouse = now // Update the time we last detected motion
			if name == config.KitchenMotionSensor {
				a.lastseenMotionInKitchenArea = now
				a.SetDelayTimeAction("Turnoff front door hall light", 4*time.Minute, func(ta *TimedAction, a *Automation) { a.TurnOffDevice("Front door hall light") })
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
				a.SetDelayTimeAction(config.KitchenLights, 5*time.Minute, func(ta *TimedAction, a *Automation) { a.TurnOffDevice(config.KitchenLights) })
				a.SetDelayTimeAction(config.KitchenLights, 5*time.Minute, func(ta *TimedAction, a *Automation) { a.TurnOffDevice(config.KitchenLights) })
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
					a.TurnOnDevice(config.BedroomLights)
					a.TurnOnDevice(config.BedroomChandelierLightSwitch)
				}
			}
		} else if value == "off" {
			if a.timeofday != "night" && a.timeofday != "sleeptime" {
				if lastseenDuration > (time.Duration(30) * time.Minute) {
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
			a.sendNotification("Front door opened")
			a.TurnOnDevice(config.FrontdoorHallLight)
			a.SetDelayTimeAction("Turnoff front door hall light", 10*time.Minute, func(ta *TimedAction, a *Automation) { a.TurnOffDevice("Front door hall light") })

			a.lastseenMotionInKitchenArea = time.Now()
		} else if value == "close" {
			a.sendNotification("Front door closed")
			a.SetDelayTimeAction("Turnoff front door hall light", 5*time.Minute, func(ta *TimedAction, a *Automation) { a.TurnOffDevice("Front door hall light") })
			a.presence.reset()
		}
	}
}

func (a *Automation) sendNotification(message string) {
	a.pubsub.Publish("shout/message/", message)
}

func (a *Automation) updatePresence(name string, presence bool) (current bool, previous bool) {
	var exists bool
	previous, exists = a.presence.presence[name]
	if !exists {
		previous = false
		a.presence.presence[name] = presence
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

// ActionDelegate is the action delegate
type ActionDelegate func(ta *TimedAction, a *Automation)

// TimedAction holds the 'when' and 'action' necessary to trigger
type TimedAction struct {
	Name   string
	When   time.Time
	Action ActionDelegate
}

// SetRealTimeAction sets an action that will trigger at a specific time
func (a *Automation) SetRealTimeAction(name string, hour int, minute int, ad ActionDelegate) {
	action, exists := a.timedActions[name]
	if !exists {
		when := time.Date(a.now.Year(), a.now.Month(), a.now.Day(), hour, minute, 0, 0, a.now.Location())
		action = &TimedAction{Name: name, When: when, Action: ad}
		a.timedActions[name] = action
	}
}

// SetDelayTimeAction sets an action that trigger after 'duration'
func (a *Automation) SetDelayTimeAction(name string, duration time.Duration, actiondelegate ActionDelegate) {
	action, exists := a.timedActions[name]
	if !exists {
		action = &TimedAction{Name: name, When: time.Now().Add(duration), Action: actiondelegate}
		a.timedActions[name] = action
	} else {
		action.When = time.Now().Add(duration)
	}
}

// Tick returns true when the action has ended, false otherwise
func (ta *TimedAction) Tick(now time.Time) bool {
	return now.After(ta.When)
}

// UpdateTimedActions will tick all actions and remove any that have ended
func (a *Automation) UpdateTimedActions() {
	for name, ta := range a.timedActions {
		if ta.Tick(a.now) == true {
			ta.Action(ta, a)
			delete(a.timedActions, name)
		}
	}
}
