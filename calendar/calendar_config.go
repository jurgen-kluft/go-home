// To parse and unparse this JSON data, add this code to your project and do:
//
//    r, err := UnmarshalCcalendar(bytes)
//    bytes, err = r.Marshal()

package calendar

import (
	"encoding/json"
	"fmt"
)

func unmarshalccalendar(data []byte) (*Ccalendar, error) {
	r := &Ccalendar{}
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Ccalendar) marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Ccalendar struct {
	Calendars []Ccal    `json:"calendars"`
	Event     []Cevent  `json:"event"`
	Policy    []Cpolicy `json:"policy"`
}

func (c *Ccalendar) print() {
	for _, cal := range c.Calendars {
		cal.print()
	}
	for _, evn := range c.Event {
		evn.print()
	}
	for _, pol := range c.Policy {
		pol.print()
	}
}

type Ccal struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (c Ccal) print() {
	fmt.Printf("ccal.name = %s\n", c.Name)
	fmt.Printf("ccal.url = %s\n", c.URL)
}

type Cevent struct {
	Calendar string   `json:"calendar"`
	Domain   string   `json:"domain"`
	Name     string   `json:"name"`
	State    string   `json:"state"`
	Typeof   string   `json:"typeof"`
	Values   []string `json:"values"`
}

func (c Cevent) print() {
	fmt.Printf("cevent.calendar = %s\n", c.Calendar)
	fmt.Printf("cevent.domain = %s\n", c.Domain)
	fmt.Printf("cevent.name = %s\n", c.Name)
	fmt.Printf("cevent.state = %s\n", c.State)
	fmt.Printf("cevent.typeof = %s\n", c.Typeof)
}

type Cpolicy struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Policy string `json:"policy"`
}

func (c Cpolicy) print() {
	fmt.Printf("cpolicy.domain = %s\n", c.Domain)
	fmt.Printf("cpolicy.name = %s\n", c.Name)
	fmt.Printf("cpolicy.policy = %s\n", c.Policy)
}
