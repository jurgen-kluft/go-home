package calendar

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-icloud-calendar"
	"github.com/jurgen-kluft/hass-go/state"
)

type Calendar struct {
	ccal   *Ccalendar
	events map[string]Cevent
	cals   []*icalendar.Calendar
	update time.Time
}

func (c *Calendar) readConfig() (*Ccalendar, error) {
	jsonBytes, err := ioutil.ReadFile("config/calendar.json")
	if err != nil {
		return nil, fmt.Errorf("ERROR: Failed to read calendar config ( %s )", err)
	}
	ccal, err := unmarshalccalendar(jsonBytes)
	return ccal, err
}

func New() (*Calendar, error) {
	var err error

	c := &Calendar{}
	c.events = map[string]Cevent{}
	c.ccal, err = c.readConfig()
	if err != nil {
		fmt.Printf("ERROR: '%s'\n", err.Error())
	}
	//c.ccal.print()
	for _, ev := range c.ccal.Event {

		c.events[strings.ToLower(ev.Domain)+":"+strings.ToLower(ev.Name)] = ev
	}
	for _, cal := range c.ccal.Calendars {
		if strings.HasPrefix(cal.URL, "http") {
			c.cals = append(c.cals, icalendar.NewURLCalendar(cal.URL))
		} else if strings.HasPrefix(cal.URL, "file") {
			c.cals = append(c.cals, icalendar.NewFileCalendar(cal.URL))
		} else if strings.HasPrefix(cal.URL, "file") {
			filepath := strings.Replace(cal.URL, "file://", "", 1)
			fmt.Printf("ERROR: Unknown Calendar source: '%s'\n", filepath)
		}
	}

	c.update = time.Now()

	return c, err
}

func (c *Calendar) updateEvents(when time.Time, states *state.Instance) error {
	//fmt.Printf("Update calendar events: '%d'\n", len(c.cals))
	for _, cal := range c.cals {
		//fmt.Printf("Update calendar events: '%s'\n", cal.Name)

		eventsForDay := cal.GetEventsByDate(when)
		for _, e := range eventsForDay {
			var domain string
			var dname string
			var dstate string
			title := strings.Replace(e.Summary, ":", " : ", 1)
			title = strings.Replace(title, "=", " = ", 1)
			fmt.Sscanf(title, "%s : %s = %s", &domain, &dname, &dstate)
			//fmt.Printf("Parsed: '%s' - '%s' - '%s'\n", domain, dname, dstate)

			domain = strings.ToLower(domain)
			dname = strings.ToLower(dname)
			dstate = strings.ToLower(dstate)

			ekey := domain + ":" + dname
			ce, exists := c.events[ekey]
			if exists {
				if domain == "report" {
					states.SetStringState(domain+"."+dname, e.GenerateUUID())
					states.SetStringState(domain+"."+dname+".ID", e.GenerateUUID())
					states.SetTimeState(domain+"."+dname+".from", e.Start)
					states.SetTimeState(domain+"."+dname+".until", e.End)
				}

				if ce.Typeof == "string" {
					states.SetStringState(domain+"."+dname, dstate)
				} else if ce.Typeof == "float" {
					fstate, err := strconv.ParseFloat(dstate, 64)
					if err == nil {
						states.SetFloatState(domain+"."+dname, fstate)
					}
				}
			}
		}
	}
	return nil
}

func (c *Calendar) load() (err error) {
	for _, cal := range c.cals {
		err = cal.Load()
	}
	return err
}

func timeInRange(when time.Time, rangeBegin time.Time, rangeEnd time.Time) bool {
	t := when.Unix()
	rb := rangeBegin.Unix()
	re := rangeEnd.Unix()
	return t >= rb && t < re
}

func weekOrWeekEndStartEnd(now time.Time) (weekend bool, westart, weend, wdstart, wdend time.Time) {
	day := now.Day()

	westart = now
	weend = now

	wdstart = now
	wdend = now

	if now.Weekday() == time.Friday || now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		if now.Weekday() == time.Saturday {
			day--
		} else if now.Weekday() == time.Sunday {
			day -= 2
		}
		westart = time.Date(now.Year(), now.Month(), day, 18, 0, 0, 0, now.Location())
		weend = time.Date(now.Year(), now.Month(), day+2, 18, 0, 0, 0, now.Location())
		if timeInRange(now, westart, weend) {
			weekend = true
			wdstart = weend
			wdend = time.Date(wdstart.Year(), wdstart.Month(), wdstart.Day()+5, 18, 0, 0, 0, now.Location())
		} else {
			weekend = false
			if now.Weekday() == time.Friday {
				wdend = westart
				wdstart = time.Date(wdend.Year(), wdend.Month(), wdend.Day()-4, 18, 0, 0, 0, now.Location())
			} else {
				wdstart = weend
				wdend = time.Date(wdstart.Year(), wdstart.Month(), wdstart.Day()+5, 18, 0, 0, 0, now.Location())
			}
		}
	} else {
		if now.Weekday() == time.Monday {
			day--
		} else if now.Weekday() == time.Tuesday {
			day -= 2
		} else if now.Weekday() == time.Wednesday {
			day -= 3
		} else if now.Weekday() == time.Thursday {
			day -= 4
		} else if now.Weekday() == time.Friday {
			day -= 5
		}

		wdstart = time.Date(now.Year(), now.Month(), day, 18, 0, 0, 0, now.Location())
		wdend = time.Date(now.Year(), now.Month(), day+4, 18, 0, 0, 0, now.Location())

		weekend = false
		westart = wdend
		weend = time.Date(westart.Year(), westart.Month(), westart.Day()+2, 18, 0, 0, 0, now.Location())
	}

	return weekend, westart, weend, wdstart, wdend
}

func (c *Calendar) findPolicy(domain string, name string, policy string) (bool, string) {
	for _, p := range c.ccal.Policy {
		if p.Domain == domain && p.Name == name {
			pname := ""
			pvalue := ""
			n, err := fmt.Sscanf(p.Policy, "%s = %s", &pname, &pvalue)
			if n == 2 && err == nil {
				if pname == policy {
					return true, pvalue
				}
			}
		}
	}
	return false, ""
}

// Process will update 'events' from the calendar
func (c *Calendar) Process(states *state.Instance) time.Duration {
	var err error
	now := states.GetTimeState("time.now", time.Now())

	if now.Unix() >= c.update.Unix() {
		// Download again after 15 minutes
		c.update = time.Unix(now.Unix()+15*60, 0)

		// Download calendar
		// fmt.Println("CALENDAR: LOAD")
		err := c.load()
		if err != nil {
			fmt.Printf("ERROR: '%s'\n", err.Error())
			return 15 * time.Minute
		}
	}

	// Other general states
	weekend, weStart, weEnd, wdStart, wdEnd := weekOrWeekEndStartEnd(now)

	//fmt.Println("CALENDAR: DEFAULT")

	// Default all states before updating them
	for _, eevent := range c.events {
		if eevent.Typeof == "string" {
			states.SetStringState(eevent.Domain+"."+eevent.Name, eevent.State)
			if weekend {
				policyOk, policyValue := c.findPolicy(eevent.Domain, eevent.Name, "weekend")
				if policyOk {
					states.SetStringState(eevent.Domain+"."+eevent.Name, policyValue)
				}
			} else {
				policyOk, policyValue := c.findPolicy(eevent.Domain, eevent.Name, "!weekend")
				if policyOk {
					states.SetStringState(eevent.Domain+"."+eevent.Name, policyValue)
				}
			}
		} else if eevent.Typeof == "float" {
			fstate, err := strconv.ParseFloat(eevent.State, 64)
			if err == nil {
				states.SetFloatState(eevent.Domain+"."+eevent.Name, fstate)
			}
		}
	}

	// Update events
	//fmt.Println("CALENDAR: UPDATE EVENTS")
	err = c.updateEvents(now, states)
	if err != nil {
		fmt.Printf("ERROR: '%s'\n", err.Error())
		return 1 * time.Minute
	}

	states.SetBoolState("calendar.weekend", weekend)
	states.SetBoolState("calendar.weekday", !weekend)
	states.SetTimeState("calendar.weekend.start", weStart)
	states.SetTimeState("calendar.weekend.end", weEnd)
	states.SetTimeState("calendar.weekday.start", wdStart)
	states.SetTimeState("calendar.weekday.end", wdEnd)
	states.SetStringState("calendar.weekend.title", "Weekend")
	states.SetStringState("calendar.weekday.title", "Weekday")
	states.SetStringState("calendar.weekend.description", "Saturday and Sunday")
	states.SetStringState("calendar.weekday.description", "Monday to Friday")

	return 1 * time.Minute
}
