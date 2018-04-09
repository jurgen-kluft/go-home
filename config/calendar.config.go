// To parse and unparse this JSON data, add this code to your project and do:
//
//    r, err := UnmarshalCcalendar(bytes)
//    bytes, err = r.Marshal()

package config

import (
	"encoding/json"
	"fmt"
)

func CalendarConfigFromJSON(jsonstr string) (*CalendarConfig, error) {
	r := &CalendarConfig{}
	err := json.Unmarshal([]byte(jsonstr), &r)
	return r, err
}

func (r *CalendarConfig) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type CalendarConfig struct {
	Calendars []Ccal    `json:"calendars"`
	Sensors   []Csensor `json:"sensors"`
	Rules     []Crule   `json:"rules"`
}

func (c *CalendarConfig) Print() {
	for _, cal := range c.Calendars {
		cal.print()
	}
	for _, sensor := range c.Sensors {
		sensor.print()
	}
	for _, rule := range c.Rules {
		rule.print()
	}
}

type Ccal struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (c Ccal) print() {
	fmt.Printf("calendar.name = %s\n", c.Name)
	fmt.Printf("calendar.url = %s\n", c.URL)
}

type Csensor struct {
	Name  string `json:"name"`
	State string `json:"state"`
	Type  string `json:"type"`
}

func (c Csensor) print() {
	fmt.Printf("sensor.name = %s\n", c.Name)
	fmt.Printf("sensor.state = %s\n", c.State)
	fmt.Printf("sensor.type = %s\n", c.Type)
}

type Crule struct {
	Key    string `json:"key"`
	State  string `json:"state"`
	IfThen IfThen `json:"if"`
}

type IfThen struct {
	Key   string `json:"key"`
	State string `json:"state"`
}

func (c Crule) print() {
	fmt.Printf("rule.key = %s\n", c.Key)
	fmt.Printf("rule.state = %s\n", c.State)
	fmt.Printf("rule.if.key = %s\n", c.IfThen.Key)
	fmt.Printf("rule.if.state = %s\n", c.IfThen.State)
}
