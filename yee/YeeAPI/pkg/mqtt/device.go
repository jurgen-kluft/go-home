package mqtt

import (
	"fmt"
	"time"

	"github.com/thomasf/lg"

	"github.com/thomasf/yeelight/pkg/yeel"
)

// Device .
type Device struct {
	yeel.Device
	LastUpdated time.Time
	Transition  time.Duration
	state       map[string]string
}

func (t *Device) Updates() []PropUpdate {
	if t.state == nil {
		t.state = make(map[string]string)
	}
	prev := t.state
	next := t.stateMap()
	var updates []PropUpdate
loop:
	for prop, nextValue := range next {
		prevValue, prevOk := prev[prop]
		if prevOk && prevValue == nextValue {
			continue loop
		}
		updates = append(updates, PropUpdate{
			DeviceID: t.Device.ID,
			Prop:     prop,
			Value:    []byte(nextValue),
		})
		t.state[prop] = nextValue
	}

	return updates
}

func (d *Device) Command(msg Command) (yeel.Commander, error) {
	// if v, ok := propCommandMap[msg.Command]; ok {
	if m, ok := models[d.Model]; ok {
		if c, ok := m.commanders[msg.Command]; ok {
			return c(d, &msg)
		}
	}
	lg.V(5).Infof("prop printer not found for prop type: %s", msg.Command)
	return nil, fmt.Errorf("%v", msg.Command)
}

func (d *Device) Update(n yeel.Notification) {
	p := n.Params

	for _, v := range []struct{ next, prev *int }{
		{p.RGB, &d.RGB},
		{p.Bright, &d.Brightness},
		{p.Hue, &d.Hue},
		{p.ColorMode, &d.ColorMode},
		{p.CT, &d.ColorTemprature},
		{p.Sat, &d.Saturation},
		{p.Flowing, &d.Flowing},
	} {
		if v.next != nil {
			*v.prev = *v.next
		}
	}
	for _, v := range []struct{ next, prev *string }{
		{p.Power, &d.Power},
		{p.Name, &d.Name},
		{p.FlowParams, &d.FlowParams},
	} {
		if v.next != nil {
			*v.prev = *v.next
		}
	}
	if p.RGB != nil {
		d.ColorMode = 1
	}
	if p.CT != nil {
		d.ColorMode = 2
	}
	if p.Hue != nil || p.Sat != nil {
		d.ColorMode = 3
	}
}

func (d *Device) Value(prop string) string {
	if m, ok := models[d.Model]; ok {
		if p, ok := m.printers[d.Model]; ok {
			return p(d)
		}
	}

	lg.V(5).Infof("prop printer not found for prop type: %s", prop)
	return ""
}

func (d *Device) stateMap() map[string]string {
	m := make(map[string]string, 0)
	if mds, ok := models[d.Model]; ok {
		for prop, fn := range mds.printers {
			m[prop] = fn(d)
		}
	} else {
		lg.Errorln(d.Model)
	}

	return m
}
