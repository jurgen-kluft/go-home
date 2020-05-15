package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/micro-service"
	"github.com/jurgen-kluft/go-icloud-calendar"
)

// Calendar ...
type Calendar struct {
	name         string
	config       *config.CalendarConfig
	sensors      map[string]config.Csensor
	sensorStates map[string]*config.SensorState
	cals         []*icalendar.Calendar
	update       time.Time
	service      *microservice.Service
}

func new() *Calendar {
	c := &Calendar{}
	c.name = "calendar"
	c.sensors = map[string]config.Csensor{}
	c.sensorStates = map[string]*config.SensorState{}
	return c
}

// initialize  ... create a new Calendar from the given JSON configuration
func (c *Calendar) initialize(jsondata []byte) (err error) {
	c.config, err = config.CalendarConfigFromJSON(jsondata)
	if err != nil {
		return err
	}
	//c.ccal.print()
	for _, sn := range c.config.Sensors {
		ekey := strings.ToLower(sn.Name)
		c.sensors[ekey] = sn
		sensor := &config.SensorState{Name: ekey, Time: time.Now()}
		if sn.Type == "string" {
			sensor.AddStringAttr(ekey, sn.State)
		} else if sn.Type == "float" {
			value, _ := strconv.ParseFloat(sn.State, 64)
			sensor.AddFloatAttr(ekey, value)
		} else if sn.Type == "bool" {
			value, _ := strconv.ParseBool(sn.State)
			sensor.AddBoolAttr(ekey, value)
		}
		c.sensorStates[ekey] = sensor
	}

	for _, cal := range c.config.Calendars {
		if strings.HasPrefix(cal.URL.String, "http") {
			ical := icalendar.NewURLCalendar(cal.URL.String)
			c.cals = append(c.cals, ical)
		} else if strings.HasPrefix(cal.URL.String, "file") {
			c.cals = append(c.cals, icalendar.NewFileCalendar(cal.URL.String))
		} else {
			c.service.Logger.LogError(c.name, fmt.Sprintf("Unknown calendar source '%s'", cal.URL))
		}
		c.service.Logger.LogInfo(c.name, cal.URL.String)
	}

	c.update = time.Now()
	return c.load()
}

func (c *Calendar) updateSensorStates(when time.Time) error {
	fmt.Printf("Update calendar events: '%d'\n", len(c.cals))

	for _, cal := range c.cals {
		fmt.Printf("Update calendar events: '%s'\n", cal.Name)
		eventsForDay := cal.GetEventsFor(when)

		if len(eventsForDay) == 0 {
			fmt.Printf("Calendar '%s' has no events\n", cal.Name)
		}

		for _, event := range eventsForDay {
			var dname string
			var dstate string
			title := strings.Replace(event.Summary, ":", ".", 3)
			title = strings.Replace(title, "=", " = ", 1)
			n, err := fmt.Sscanf(title, "%s = %s", &dname, &dstate)
			fmt.Printf("Parsed: '%s' - '%s' - '%s' - '%s'\n", event.Summary, title, dname, dstate)
			if n == 2 && err == nil {
				dname = strings.ToLower(strings.Trim(dname, " "))
				dstate = strings.ToLower(strings.Trim(dstate, " "))
				ekey := dname

				sensor, exists := c.sensorStates[ekey]
				if exists {
					sensor.Time = time.Now()
					if sensor.StringAttrs != nil {
						sensor.StringAttrs[0].Value = dstate
					} else if sensor.IntAttrs != nil {
						value, _ := strconv.ParseInt(dstate, 10, 64)
						sensor.IntAttrs[0].Value = value
					} else if sensor.FloatAttrs != nil {
						value, _ := strconv.ParseFloat(dstate, 64)
						sensor.FloatAttrs[0].Value = value
					} else if sensor.BoolAttrs != nil {
						value, _ := strconv.ParseBool(dstate)
						sensor.BoolAttrs[0].Value = value
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

func (c *Calendar) applyRulesToSensorStates() {
	for _, p := range c.config.Rules {
		var sensor *config.SensorState
		var ifsensor *config.SensorState
		var exists bool
		sensor, exists = c.sensorStates[p.Key]
		if exists && sensor.StringAttrs != nil {
			for _, ifthen := range p.IfThen {
				ifsensor, exists = c.sensorStates[ifthen.Key]
				if exists && ifsensor.StringAttrs != nil {
					ifthenValue := ifsensor.StringAttrs[0].Value
					if ifthenValue == ifthen.State {
						sensor.StringAttrs[0].Value = p.State
					}
				} else if exists && ifsensor.BoolAttrs != nil {
					ifthenValue := ifsensor.BoolAttrs[0].Value
					if ifthenValue && ("true" == ifthen.State) {
						sensor.StringAttrs[0].Value = p.State
					} else if !ifthenValue && ("false" == ifthen.State) {
						sensor.StringAttrs[0].Value = p.State
					}
				} else {
					c.service.Logger.LogError(c.name, fmt.Sprintf("Logical error when applying rules to sensor states (%s)", p.Key+", "+p.State))
				}
			}
		}
	}
}

func (c *Calendar) publishSensorState(name string, sensorjsonbytes []byte) {
	c.service.Pubsub.Publish(fmt.Sprintf("state/sensor/%s/", name), sensorjsonbytes)
}

// Process will update 'events' from the calendar
func (c *Calendar) process() (err error) {
	now := time.Now()

	// Other general states
	weekendsensor, exists := c.sensorStates["weekend"]
	if exists {
		if weekendsensor.BoolAttrs != nil {
			weekend, _, _, _, _ := weekOrWeekEndStartEnd(now)
			weekendsensor.BoolAttrs[0].Value = weekend
			sensorjsonbytes, err := weekendsensor.ToJSON()
			if err == nil {
				c.publishSensorState("weekend", sensorjsonbytes)
			} else {
				return err
			}
		}
	}

	// Update sensors and apply configured rules to sensors
	err = c.updateSensorStates(now)
	if err == nil {
		c.applyRulesToSensorStates()

		// Publish sensor states
		for _, ss := range c.sensorStates {
			sensorjsonbytes, err := ss.ToJSON()
			if err == nil {
				c.publishSensorState(ss.Name, sensorjsonbytes)
			} else {
				c.service.Logger.LogError(c.name, err.Error())
			}
		}
	}

	return err
}

func main() {
	register := []string{"config/calendar/"}
	subscribe := []string{"config/calendar/", "config/request/"}

	m := microservice.New("calendar")
	m.RegisterAndSubscribe(register, subscribe)

	c := new()
	c.service = m

	m.RegisterHandler("config/calendar/", func(m *microservice.Service, topic string, msg []byte) bool {
		m.Logger.LogInfo(m.Name, "received configuration")
		err := c.initialize(msg)
		if err != nil {
			c = nil
			m.Logger.LogError(c.name, err.Error())
		} else {
			// Register emitter channel for every sensor
			for _, ss := range c.sensorStates {
				if err = m.Register(fmt.Sprintf("state/sensor/%s/", ss.Name)); err != nil {
					m.Logger.LogError("pubsub", err.Error())
				}
			}
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if tickCount%5 == 0 {
			if c != nil && c.config == nil {
				m.Pubsub.PublishStr("config/request/", m.Name)
			}
		}
		if tickCount%150 == 0 {
			if c != nil && c.config != nil {
				if err := c.load(); err != nil {
					m.Logger.LogError(m.Name, err.Error())
				} else {
					m.Logger.LogInfo(m.Name, "(re)loading calendars succesfully")
				}
			}
		}
		if tickCount%5 == 0 {
			if c != nil && c.config != nil {
				if err := c.process(); err != nil {
					m.Logger.LogError(m.Name, err.Error())
				}
			}
		}
		tickCount++
		return true
	})

	m.Loop()
}
