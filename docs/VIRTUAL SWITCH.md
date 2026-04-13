# Virtual Switch

The VirtualSwitch is a specialized switch that turns on or off when the physical switch is turned on by an external
source or under certain conditions by motion.
This means that when presence is detected, but the physical switch is off that the
virtual switch will not turn on. When the physical switch is turned on, the virtual switch will turn on immediately.
When the physical switch is physically turned off, the virtual switch will turn off immediately. However, when the
physical switch is ON and the presence detects that there is no presence, the virtual switch will turn off
after a delay of 5 minutes and the state is such that when presence is again detected, the virtual switch will turn
on immediately. If in the meantime the physical switch is turned off, the virtual switch will turn off immediately
and the presence sensor will not have any effect until the physical switch is turned on again.


```go
type PresenceSensor interface {
	SetPresenceState(state bool)
}

type PhysicalSwitch interface {
	SetSwitchState(state bool)
}

// VirtualSwitch is a specialized switch that turns on or off based on the state of a physical switch and a presence sensor.
// Interfaces: PresenceSensor, PhysicalSwitch
type VirtualSwitch struct {
	Name                string
	PhysicalSwitch      bool
	PresenceSensor      bool
	VirtualState        bool // On / Off, the state of the virtual switch, which may differ from the physical switch based on the presence sensor
	ActualPhysicalState bool // On / Off,
	DelayOff            int  // Seconds, delay in seconds before turning off the physical switch
}

func (vs *VirtualSwitch) SetPresenceState(state bool) {
	vs.PresenceSensor = state
	if vs.PhysicalSwitch && state {
		vs.VirtualState = true
	} else if vs.PhysicalSwitch && !state {
		// Start a timer to turn off the virtual switch after the delay
		go func() {
			time.Sleep(time.Duration(vs.DelayOff) * time.Second)
			if !vs.PresenceSensor {
				vs.VirtualState = false
			}
		}()
	} else {
		vs.VirtualState = false
	}
}

func (vs *VirtualSwitch) SetSwitchState(state bool) {
	vs.PhysicalSwitch = state
	if state {
		if vs.PresenceSensor {
			vs.VirtualState = true
		}
	} else {
		vs.VirtualState = false
	}
}
```

