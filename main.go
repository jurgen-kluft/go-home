package main

import (
	"time"
)

const (
	_10seconds = 10
	_30seconds = 30
)

type TimeOfDayLapse struct {
	Name       string
	Start, End time.Time
}

func (t TimeOfDayLapse) Is(name string) bool {
	return t.Name == name
}

type TimeOfDayConfig struct {
	Current TimeOfDayLapse
	Lapses  map[string]TimeOfDayLapse
}

type HomeLights struct {
}

func (lights *HomeLights) IsOff(lightgroup string) bool {
	return false
}

func (lights *HomeLights) TurnAllOffIn(seconds int) {

}

func (lights *HomeLights) TurnOffIn(seconds int, lightgroups ...string) {

}

func (lights *HomeLights) TurnOnIn(seconds int, lightgroups ...string) {

}

type HomeSwitches struct {
}

func (switches *HomeSwitches) TurnAllOffIn(seconds int) {

}

func (switches *HomeSwitches) TurnOnIn(seconds int, switchgroups ...string) {

}

type Home struct {
	timeofday *TimeOfDayConfig
	presence  *PresenceProcessor
	Lights    *HomeLights
	Switches  *HomeSwitches
}

type Database interface {
	Save(key string, value string)
	Load(key string, value string)
}

type DynamicObject struct {
	Key          string
	BoolFields   map[string]bool
	IntFields    map[string]int32
	FloatFields  map[string]float32
	StringFields map[string]string
	ObjectFields map[string]*DynamicObject
}

type Presence uint32

const (
	Absent Presence = iota
	Present
	Leaving
	Arriving
)

func (p Presence) IsPresent() bool {
	return p == Present
}
func (p Presence) IsAbsent() bool {
	return p == Absent
}
func (p Presence) IsLeaving() bool {
	return p == Leaving
}
func (p Presence) IsArriving() bool {
	return p == Arriving
}

type PresenceProcessor struct {
}

// Example: p.OnEvent("Grandpa=Arriving;Grandma=Arriving;Faith=Leaving;Jurgen=Home")
func (p PresenceProcessor) OnEvent(event string) {

}

func (p PresenceProcessor) PresenceOf(member ...string) Presence {
	return 0
}

// Save
func (p PresenceProcessor) Save(db Database) {

}

// Load
func (p PresenceProcessor) Load(db Database) {

}

func main() {

	for true {
		// START

		// Connect to REDIS, if unable to connect, wait for N seconds then restart
		// Subscribe to chennel "Go-Home"

		// Create Home
		home := &Home{}

		for true {

			// Wait for message from "Go-Home" channel, this can be any of the following
			// - "Update"; this is the tick
			// - "Presence"
			// - "TimeOfDay"
			// - "Lights";
			// - "Switches"

			// Process message
			home.Logic()

			// Any error, like from REDIS -> break from this loop
		}

		// Disconnect from anything
	}
}

// Logic is the main logic procedure
func (home *Home) Logic() {
	time_of_day := home.timeofday.Current

	family := home.presence.PresenceOf("Faith", "Jurgen", "GrandPa", "GrandMa")
	mumdad := home.presence.PresenceOf("Faith", "Jurgen")

	if family.IsAbsent() {
		// Do nothing for now

	} else if family.IsLeaving() {
		home.Lights.TurnAllOffIn(_30seconds)
		home.Switches.TurnAllOffIn(_30seconds)
	} else if family.IsArriving() {
		home.Lights.TurnOnIn(_10seconds, "Kitchen")
		home.Lights.TurnOnIn(_10seconds, "Living Room")
		home.Switches.TurnOnIn(_10seconds, "Christmas Tree")
	} else if family.IsPresent() {
		if time_of_day.Is("Lunch") {
			home.Lights.TurnOnIn(_10seconds, "Kitchen")
		}
	} else if time_of_day.Is("Evening") {
		if mumdad.IsPresent() {
			if home.Lights.IsOff("Kitchen") && home.Lights.IsOff("Living Room") && home.Lights.IsOff("Bed Room") {
				home.Lights.TurnOffIn(_10seconds, "Christmas Tree")
				home.Lights.TurnOnIn(_10seconds, "Bedroom")
			}
		}
	} else if time_of_day.Is("Morning") {
		if mumdad.IsArriving() {
			if home.Lights.IsOff("Kitchen") && home.Lights.IsOff("Living Room") {
				home.Lights.TurnOnIn(_10seconds, "Christmas Tree")
				home.Lights.TurnOnIn(_10seconds, "Kitchen")
				home.Lights.TurnOnIn(_10seconds, "Living Room")
			}
		} else if time_of_day.Is("Breakfast") {
			home.Lights.TurnOffIn(_10seconds, "Kitchen")
			home.Lights.TurnOffIn(_10seconds, "Living Room")
		}
	}
}
