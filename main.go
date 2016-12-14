package go-home

import (
	"time"
)

type TimeOfDayLapse struct {
	Name string
	Start, End Time
}

type TimeOfDay struct {
	Current TimeOfDayLapse
	Lapses map[string]TimeOfDayLapse
}

type Home struct {
	timeofday *TimeOfDay
}

type Database interface {
    Save(key string, value string)
    Load(key string, value string)
}

type DynamicObject struct {
	Key string
	BoolFields map[string] bool
	IntFields map[string] int32
	FloatFields map[string] float
	StringFields map[string] string
	ObjectFields map[string] *DynamicObject
}


type PresenceProcessor struct {
	

}

// Example: p.OnEvent("Grandpa=Arriving;Grandma=Arriving;Faith=Leaving;Jurgen=Home")
func (p PresenceProcessor) OnEvent(event string) {

}

// Example: p.IsPresent("family")
func (p PresenceProcessor) IsPresent(member string) bool {
	return false
}

// Save
func (p PresenceProcessor) Save(db Database) {

}

// Load
func (p PresenceProcessor) Load(db Database) {

}



type TimeOfDay struct {
	Time start, end
	Name string
}


func HomeLogic() {
    time_of_day := home.timeofday.current;
    
    family = presence state of family members
	mumdad = presence state of mum and dad

	IF (family == ABSENT) {
	   // Do nothing for now

	} else IF (family == LEAVING) {
	   home.Lights.TurnAllOffIn(30_seconds);
	   home.Switches.TurnAllOffIn(30_seconds);
	} else IF (family == ARRIVING) {
	   home.Lights.TurnOnIn(10_seconds, "Kitchen");
	   home.Lights.TurnOnIn(10_seconds, "Living Room");
	   home.Switches.TurnOnIn(10_seconds, "Christmas Tree");
	} else IF (family == PRESENT) {
	   IF (time_of_day.IsLunch()) {
	      home.Lights.TurnOnIn(10_seconds, "Kitchen");
	   }
	} else IF (time_of_day.Is("Evening") {
	   IF (mumdad == PRESENT) {
	      IF (home.Lights.IsOff("Kitchen") && home.Lights.IsOff("Living Room") && home.Lights.IsOff("Bed Room")) {
	         home.Lights.TurnOffIn(10_seconds, "Christmas Tree");
	         home.Lights.TurnOnIn(10_seconds, "Bedroom");
	      }
	   }
	} else IF (time_of_day.Is("Morning")) {
	   IF (mumdad == ARRIVING) {
	      IF (home.Lights.IsOff("Kitchen") && home.Lights.IsOff("Living Room")) {
	         home.Lights.TurnOnIn(10_seconds, "Christmas Tree");
	         home.Lights.TurnOnIn(10_seconds, "Kitchen");
	         home.Lights.TurnOnIn(10_seconds, "Living Room");
	      }
	   } else IF (time_of_day.IsBreakfast()) {
	         home.Lights.TurnOffIn(10_seconds, "Kitchen");
	         home.Lights.TurnOffIn(10_seconds, "Living Room");
	   }
	}
}