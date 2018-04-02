package calendar

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jurgen-kluft/go-home/config"
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
}

// New  ... create a new Calendar from the given JSON configuration
func New(jsonstr string) (*Calendar, error) {
	var err error

	c := &Calendar{}
	c.sensors = map[string]config.Csensor{}
	c.sensorStates = map[string]*config.SensorState{}
	c.config, err = config.CalendarConfigFromJSON(jsonstr)
	if err != nil {
		fmt.Printf("ERROR: '%s'\n", err.Error())
	}
	//c.ccal.print()
	for _, sn := range c.config.Sensors {
		ekey := strings.ToLower(sn.Domain) + ":" + strings.ToLower(sn.Product) + ":" + strings.ToLower(sn.Name)
		c.sensors[ekey] = sn
		sensor := &config.SensorState{Domain: sn.Domain, Product: sn.Product, Name: sn.Name, Type: sn.Type, Value: sn.State, Time: time.Now()}
		c.sensorStates[ekey] = sensor
	}

	for _, cal := range c.config.Calendars {
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

func (c *Calendar) updateSensorStates(when time.Time) error {
	//fmt.Printf("Update calendar events: '%d'\n", len(c.cals))

	for _, cal := range c.cals {
		//fmt.Printf("Update calendar events: '%s'\n", cal.Name)

		eventsForDay := cal.GetEventsByDate(when)
		for _, e := range eventsForDay {
			var domain string
			var dproduct string
			var dname string
			var dstate string
			title := strings.Replace(e.Summary, ":", " : ", 1)
			title = strings.Replace(title, "=", " = ", 1)
			fmt.Sscanf(title, "%s : %s : %s = %s", &domain, &dproduct, &dname, &dstate)
			//fmt.Printf("Parsed: '%s' - '%s' - '%s' - '%s'\n", domain, dproduct, dname, dstate)

			domain = strings.ToLower(domain)
			dproduct = strings.ToLower(dproduct)
			dname = strings.ToLower(dname)
			dstate = strings.ToLower(dstate)
			ekey := domain + ":" + dproduct + ":" + dname

			sensor, exists := c.sensorStates[ekey]
			if exists {
				sensor.Value = dstate
				sensor.Time = time.Now()
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
		if exists {
			ifthen, exists = c.sensorStates[p.IfThen.Key]
			if exists {
				if ifthen.Type == "string" && ifthen.Value == p.IfThen.State {
					sensor.Value = p.State
				}
			}
		}
	}
}

func publishSensorState(s *config.SensorState, client *pubsub.Context) {
	jsonbytes, err := json.Marshal(s)
	if err == nil {
		client.Publish("state/sensor/calendar", string(jsonbytes))
	}
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
			fmt.Printf("ERROR: '%s'\n", err.Error())
			return
		}
	}

	// Other general states
	weekend, _, _, _, _ := weekOrWeekEndStartEnd(now)
	sensor := &config.SensorState{Domain: "sensor", Product: "calendar", Name: "weekend", Type: "bool", Value: fmt.Sprintf("%v", weekend), Time: time.Now()}
	publishSensorState(sensor, client)

	// Update sensors and apply configured rules to sensors
	err = c.updateSensorStates(now)
	if err != nil {
		c.applyRulesToSensorStates()

		// Publish sensor states
		for _, ss := range c.sensorStates {
			publishSensorState(ss, client)
		}
	}
}

func main() {

	var calendar *Calendar

	for {
		connected := true
		for connected {
			client := pubsub.New()
			err := client.Connect("calendar")
			if err == nil {

				// Subscribe to the presence demo channel
				client.Subscribe("config/calendar")

				for connected {
					select {
					case msg := <-client.InMsgs:
						topic := msg.Topic()
						if topic == "calendar/config" {
							jsonmsg := string(msg.Payload())
							calendar, err = New(jsonmsg)
							if err != nil {
								calendar = nil
							}
						} else if topic == "client/disconnected" {
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
			} else {
				panic("Error on Client.Connect(): " + err.Error())
			}
		}

		// Wait for 10 seconds before retrying
		time.Sleep(10 * time.Second)
	}

}
