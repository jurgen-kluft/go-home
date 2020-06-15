// To parse and unparse this JSON data, add this code to your project and do:
//
//    r, err := UnmarshalCcalendar(bytes)
//    bytes, err = r.Marshal()

package config

import (
	"encoding/json"
	"fmt"
)

func CalendarConfigFromJSON(data []byte) (*CalendarConfig, error) {
	r := &CalendarConfig{}
	err := json.Unmarshal(data, r)
	return r, err
}

// FromJSON converts a json string to a CalendarConfig instance
func (r *CalendarConfig) FromJSON(data []byte) error {
	c := CalendarConfig{}
	err := json.Unmarshal(data, &c)
	*r = c
	return err
}

// ToJSON converts a CalendarConfig to a JSON string
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
	Name string      `json:"name"`
	URL  CryptString `json:"url"`
}

func (c Ccal) print() {
	fmt.Printf("calendar.name = %s\n", c.Name)
	fmt.Printf("calendar.url = %s\n", c.URL)
}

type Csensor struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	AttrType  string `json:"attrtype"`
	AttrValue string `json:"attrvalue"`
}

func (c Csensor) print() {
	fmt.Printf("sensor.name = %s\n", c.Name)
	fmt.Printf("sensor.type = %s\n", c.Type)
	fmt.Printf("sensor.attrvalue = %s\n", c.AttrValue)
	fmt.Printf("sensor.attrtype = %s\n", c.AttrType)
}

type Crule struct {
	Key    string   `json:"key"`
	Value  string   `json:"value"`
	IfThen []IfThen `json:"if"`
}

type IfThen struct {
	Key   string `json:"key"`
	State string `json:"state"`
}

func (c Crule) print() {
	fmt.Printf("rule.key = %s\n", c.Key)
	fmt.Printf("rule.value = %s\n", c.Value)
	fmt.Println("If:")
	for _, rule := range c.IfThen {
		fmt.Printf("    if key='%s' has state:'%s'\n", rule.Key, rule.State)
	}
}
