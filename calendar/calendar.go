package calendar

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/jurgen-kluft/go-icloud-calendar"
)

// Calendar ...
type Calendar struct {
	config       *config.CalendarConfig
	sensors      map[string]config.Csensor
	sensorStates map[string]*config.SensorState
	cals         []*icalendar.Calendar
	update       time.Time
	log          *logpkg.Logger
}

// New  ... create a new Calendar from the given JSON configuration
func New(jsonstr string, log *logpkg.Logger) (*Calendar, error) {
	var err error

	c := &Calendar{}
	c.log = log
	c.sensors = map[string]config.Csensor{}
	c.sensorStates = map[string]*config.SensorState{}
	c.config, err = config.CalendarConfigFromJSON(jsonstr)
	if err != nil {
		log.LogError("calendar", err.Error())
	}
	//c.ccal.print()
	for _, sn := range c.config.Sensors {
		ekey := strings.ToLower(sn.Name)
		c.sensors[ekey] = sn
		sensor := &config.SensorState{Name: ekey, Time: time.Now()}
		if sn.Type == "string" {
			sensor.AddStringAttr(sn.Name, sn.State)
		} else if sn.Type == "float" {
			value, _ := strconv.ParseFloat(sn.State, 64)
			sensor.AddFloatAttr(sn.Name, value)
		} else if sn.Type == "bool" {
			value, _ := strconv.ParseBool(sn.State)
			sensor.AddBoolAttr(sn.Name, value)
		}
		c.sensorStates[ekey] = sensor
	}

	for _, cal := range c.config.Calendars {
		if strings.HasPrefix(cal.URL, "http") {
			c.cals = append(c.cals, icalendar.NewURLCalendar(cal.URL))
		} else if strings.HasPrefix(cal.URL, "file") {
			c.cals = append(c.cals, icalendar.NewFileCalendar(cal.URL))
		} else {
			log.LogError("calendar", fmt.Sprintf("Unknown calendar source '%s'", cal.URL))
		}
	}

	c.update = time.Now()

	return c, err
}

func (c *Calendar) updateSensorStates(when time.Time) error {
	//fmt.Printf("Update calendar events: '%d'\n", len(c.cals))

	for _, cal := range c.cals {
		//fmt.Printf("Update calendar events: '%s'\n", cal.Name)

		eventsForDay := cal.GetEventsByDate(when)
		for _, e := range eventsForDay {
			var dname string
			var dstate string
			title := strings.Replace(e.Summary, ":", ".", 3)
			title = strings.Replace(title, "=", " = ", 1)
			n, err := fmt.Sscanf(title, "%s = %s", &dname, &dstate)
			if n == 2 && err == nil {
				//fmt.Printf("Parsed: '%s' - '%s' - '%s' - '%s'\n", domain, dproduct, dname, dstate)
				dname = strings.ToLower(strings.Trim(dname, " "))
				dstate = strings.ToLower(strings.Trim(dstate, " "))
				ekey := dname

				sensor, exists := c.sensorStates[ekey]
				if exists {
					sensor.Time = time.Now()
					if sensor.StringAttrs != nil {
						(*sensor.StringAttrs)[0].Value = dstate
					} else if sensor.IntAttrs != nil {
						value, _ := strconv.ParseInt(dstate, 10, 64)
						(*sensor.IntAttrs)[0].Value = value
					} else if sensor.FloatAttrs != nil {
						value, _ := strconv.ParseFloat(dstate, 64)
						(*sensor.FloatAttrs)[0].Value = value
					} else if sensor.BoolAttrs != nil {
						value, _ := strconv.ParseBool(dstate)
						(*sensor.BoolAttrs)[0].Value = value
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
		var ifthen *config.SensorState
		var exists bool
		sensor, exists = c.sensorStates[p.Key]
		if exists && sensor.StringAttrs != nil {
			ifthen, exists = c.sensorStates[p.IfThen.Key]
			if exists && ifthen.StringAttrs != nil {
				ifthenValue := (*ifthen.StringAttrs)[0].Value
				if ifthenValue == p.IfThen.State {
					(*sensor.StringAttrs)[0].Value = p.State
				}
			}
		}
	}
}

func publishSensorState(sensorjson string, client *pubsub.Context) {
	client.Publish("state/sensor/calendar/", string(sensorjson))
}

// Process will update 'events' from the calendar
func (c *Calendar) Process(client *pubsub.Context) {
	var err error
	now := time.Now()

	if now.Unix() >= c.update.Unix() {
		// Download again after 15 minutes
		c.update = time.Unix(now.Unix()+15*60, 0)

		// Download calendars
		// fmt.Println("CALENDAR: LOAD")
		err := c.load()
		if err != nil {
			c.log.LogError("calendar", err.Error())
			return
		}
	}

	// Other general states
	weekend, _, _, _, _ := weekOrWeekEndStartEnd(now)
	sensorjson, err := config.StringAttrAsJSON("sensor.calendar.weekend", "weekend", fmt.Sprintf("%v", weekend))
	publishSensorState(sensorjson, client)

	// Update sensors and apply configured rules to sensors
	err = c.updateSensorStates(now)
	if err != nil {
		c.applyRulesToSensorStates()

		// Publish sensor states
		for _, ss := range c.sensorStates {
			jsonstr, err := ss.ToJSON()
			if err == nil {
				publishSensorState(jsonstr, client)
			}
		}
	}
}

func main() {

	var calendar *Calendar

	logger := logpkg.New("calendar")
	logger.AddEntry("emitter")
	logger.AddEntry("calendar")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/calendar/", "state/sensor/calendar/"}
		subscribe := []string{"config/calendar/"}
		err := client.Connect("calendar", register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/calendar/" {
						logger.LogInfo("calendar", "received configuration")
						jsonmsg := string(msg.Payload())
						calendar, err = New(jsonmsg, logger)
						if err != nil {
							calendar = nil
						}
					} else if topic == "client/disconnected/" {
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}
					break
				case <-time.After(time.Second * 60):
					if calendar != nil && calendar.config != nil {
						calendar.Process(client)
					}
					break

				}
			}
		}

		if err != nil {
			logger.LogError("calendar", err.Error())
		}
		time.Sleep(5 * time.Second)
	}

}
