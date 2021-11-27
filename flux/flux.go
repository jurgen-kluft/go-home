package main

import (
	"fmt"
	"math"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/metrics"
	microservice "github.com/jurgen-kluft/go-home/micro-service"
)

// Color Temperature
// URL: https://panasonic.net/es/solution-works/jiyugaoka/

func clockToSeconds(t time.Time) float64 {
	sh, sm, sc := t.Clock()
	return float64(sh*60*60) + float64(sm*60) + float64(sc)
}

// Return the factor 0.0 - 1.0 that indicates where we are in between start - end
func computeTimeSpanX(start, end, t time.Time) float64 {
	sx := clockToSeconds(start)
	ex := clockToSeconds(end)
	tx := clockToSeconds(t)
	return (tx - sx) / (ex - sx)
}

func inTimeSpan(start, end, t time.Time) bool {
	in := computeTimeSpanX(start, end, t)
	return in >= 0.0 && in <= 1.0
}

type MovingAverage struct {
	history []float64
	index   int
}

func NewFilter(sizeOfHistory int) *MovingAverage {
	filter := &MovingAverage{history: make([]float64, sizeOfHistory), index: -1}
	return filter
}

func (m *MovingAverage) Sample(sample float64) float64 {
	if m.index == -1 {
		for i := range m.history {
			m.history[i] = sample
		}
		m.index = 0
	}

	m.history[m.index] = sample
	m.index = (m.index + 1) % len(m.history)

	sum := 0.0
	for _, s := range m.history {
		sum += s
	}
	return sum / float64(len(m.history))
}

type context struct {
	name       string
	cfgch      string
	config     *config.FluxConfig
	metrics    *metrics.Metrics
	suncalc    *config.SensorState
	seasonName string
	season     *config.Season
	weather    *config.SensorState
	averageCT  *MovingAverage
	averageBRI *MovingAverage
	service    *microservice.Service
}

func new() *context {
	c := &context{}
	c.name = "flux"
	c.cfgch = "config/flux/"
	c.metrics, _ = metrics.New()
	c.metrics.Register("hue", map[string]string{"CT": "Color Temperature", "BRI": "Brightness"}, map[string]interface{}{"CT": 200.0, "BRI": 200.0})
	c.metrics.Register("yee", map[string]string{"CT": "Color Temperature", "BRI": "Brightness"}, map[string]interface{}{"CT": 200.0, "BRI": 200.0})
	c.seasonName = "spring"
	c.averageCT = NewFilter(8)
	c.averageBRI = NewFilter(8)
	return c
}

func (c *context) updateSeasonFromName(season string) {
	for _, e := range c.config.Seasons {
		if e.Name == season {
			c.season = &config.Season{}
			*c.season = e
		}
	}
}

// Process will update 'string'states and 'float'states
// States are both input and output, for example as input
// there are Season/Weather states like 'Season':'Winter'
// and 'Clouds':0.5
func (c *context) Process() error {
	if c.config == nil || c.suncalc == nil {
		return fmt.Errorf("No configuration or no suncalc information received")
	}

	now := time.Now()

	// Update our season
	c.updateSeasonFromName(c.seasonName)

	// First build our sun moments map
	sunmoments := map[string]time.Time{}
	for _, tss := range c.suncalc.TimeWndAttrs {
		sunmoments[tss.Name+".begin"] = tss.Begin
		sunmoments[tss.Name+".end"] = tss.End
	}

	// Add our custom time-points to the sun moments map
	for _, at := range c.config.SuncalcMoments {
		moment, exists := sunmoments[at.Name]
		if exists {
			moment = moment.Add(time.Duration(at.Shift) * time.Minute)
			sunmoments[at.Name+at.Tag] = moment
		}
	}

	// Update our Lighttime start and end time from the sun moments map
	for i, lt := range c.config.Lighttime {
		start, exists := sunmoments[lt.TimeSlot.StartMoment]
		if exists {
			end, exists := sunmoments[lt.TimeSlot.EndMoment]
			if exists {
				c.config.Lighttime[i].TimeSlot.StartTime = start
				c.config.Lighttime[i].TimeSlot.EndTime = end
			}
		}
	}

	// Figure out in which light time moment we are now
	bcurrent := false
	current := config.Lighttime{}
	for _, sm := range c.config.Lighttime {
		t0 := sm.TimeSlot.StartTime
		t1 := sm.TimeSlot.EndTime
		if inTimeSpan(t0, t1, now) {
			fmt.Println("Current light time from", sm.TimeSlot.StartMoment, "to", sm.TimeSlot.EndMoment)
			current = sm
			bcurrent = true
			break
		}
	}
	if !bcurrent {
		return fmt.Errorf("Unable to identify the current light time window")
	}

	// Time interpolation factor, where are we between startMoment - endMoment
	currentx := computeTimeSpanX(current.TimeSlot.StartTime, current.TimeSlot.EndTime, now)
	currentx = float64(int64(currentx*100.0)) / 100.0

	clouds := config.Weather{Clouds: config.MinMax{Min: 0.0, Max: 0.001}, CTPct: 0.0, BriPct: 0.0}
	cloudFac := float64(0.0)
	if c.weather != nil {
		cloudFac = c.weather.GetFloatAttr("clouds", 0.0)
	}
	for _, w := range c.config.Weather {
		if cloudFac >= w.Clouds.Min && cloudFac < w.Clouds.Max {
			clouds = w
			break
		}
	}

	// Full cloud cover will increase color-temperature by 10% of (Max - Current)
	// NOTE: Only during the day (twilight + light)
	// TODO: when the moon is shining in the night the amount
	//       of blue-light is also higher than normal.
	// CT = 0.0 -> Coldest (>6500K)
	// CT = 1.0 -> Warmest (2000K)
	CT := current.CT.LinearInterpolated(currentx)
	if current.Darkorlight != "dark" {
		if clouds.CTPct >= 0 {
			CT = CT + clouds.CTPct*(1.0-CT)
		} else {
			CT = CT + clouds.CTPct*CT
		}
	}
	CT = c.season.CT.LinearInterpolated(CT)
	CT = c.averageCT.Sample(CT)

	// Full cloud cover will increase brightness by 10% of (Max - Current)
	// BRI = 0 -> Very dim light
	// BRI = 1 -> Very bright light
	BRI := current.Bri.LinearInterpolated(currentx)
	BRI = BRI + cloudFac*0.1*(1.0-BRI)
	if current.Darkorlight != "dark" {
		// A bit brighter lights when there are clouds during the day.
		if clouds.BriPct >= 0 {
			BRI = BRI + clouds.BriPct*(1.0-BRI)
		} else {
			BRI = BRI + clouds.BriPct*BRI
		}
	}
	BRI = c.season.BRI.LinearInterpolated(BRI)
	BRI = c.averageBRI.Sample(BRI)

	// Publishing the following sensors:
	//  - Sensor.Light.HUE, Name = CT, Value = float64(100.0)
	//  - Sensor.Light.HUE, Name = BRI, Value = float64(100.0)
	//  - Sensor.Light.YEE, Name = CT, Value = float64(100.0)
	//  - Sensor.Light.YEE, Name = BRI, Value = float64(100.0)
	//  - Sensor.Light.DarkOrLight = string(Dark)

	for _, ltype := range c.config.Lighttype {
		sensor := config.NewSensorState("all", "flux")

		c.metrics.Begin(ltype.Name)

		lct := ltype.CT.LinearInterpolated(CT)
		sensor.AddFloatAttr("CT", math.Floor(lct))
		c.metrics.Set(ltype.Name, "CT", lct)

		lbri := ltype.BRI.LinearInterpolated(BRI)
		sensor.AddFloatAttr("BRI", math.Floor(lbri))
		c.metrics.Set(ltype.Name, "BRI", lbri)

		c.metrics.Send(ltype.Name)

		jsonbytes, err := sensor.ToJSON()
		if err == nil {
			c.publishSensor(ltype.Channel, string(jsonbytes))
		}
	}

	sensorDOL, err := config.StringAttrAsJSON("darkorlight", "illumination", "DarkOrLight", string(current.Darkorlight))
	if err == nil {
		c.publishSensor("state/sensor/darkorlight/", sensorDOL)
	}

	return nil
}

func (c *context) publishSensor(channel string, json string) {
	c.service.Logger.LogInfo(c.service.Name, "Publish at '"+channel+"' JSON ["+json+"]")
	c.service.Pubsub.Publish(channel, []byte(json))
}

func main() {
	register := []string{"config/request/", "config/flux/", "state/sensor/weather/", "state/sensor/sun/", "state/sensor/season/"}
	subscribe := []string{"config/flux/", "state/sensor/weather/", "state/sensor/sun/", "state/sensor/season/"}

	m := microservice.New("flux")
	m.RegisterAndSubscribe(register, subscribe)

	c := new()
	c.service = m

	m.RegisterHandler("config/flux/", func(m *microservice.Service, topic string, msg []byte) bool {
		var err error
		c.config, err = config.FluxConfigFromJSON(msg)
		if err == nil {
			m.Logger.LogInfo(m.Name, "received configuration")
			for _, ltype := range c.config.Lighttype {
				m.Register(ltype.Channel)
				if err == nil {
					m.Logger.LogInfo(c.name, fmt.Sprintf("registered pubsub channel %s for lighttype %s", ltype.Channel, ltype.Name))
				} else {
					m.Logger.LogError(c.name, err.Error())
				}
			}
		} else {
			m.Logger.LogError(m.Name, err.Error())
		}
		return true
	})

	m.RegisterHandler("state/sensor/weather/", func(m *microservice.Service, topic string, msg []byte) bool {
		var err error
		c.weather, err = config.SensorStateFromJSON(msg)
		if err == nil {
			m.Logger.LogInfo(c.name, "received weather state")
		} else {
			m.Logger.LogError(c.name, err.Error())
		}
		return true
	})

	m.RegisterHandler("state/sensor/sun/", func(m *microservice.Service, topic string, msg []byte) bool {
		var err error
		c.suncalc, err = config.SensorStateFromJSON(msg)
		if err == nil {
			m.Logger.LogInfo(c.name, "received sun state")
		} else {
			m.Logger.LogError(c.name, err.Error())
		}
		return true
	})

	m.RegisterHandler("state/sensor/season/", func(m *microservice.Service, topic string, msg []byte) bool {
		seasonSensorState, err := config.SensorStateFromJSON(msg)
		if err == nil {
			m.Logger.LogInfo(c.name, "received season state")
			c.seasonName = seasonSensorState.GetValueAttr("season", "winter")
		} else {
			m.Logger.LogError(c.name, err.Error())
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if (tickCount % 30) == 0 {
			err := c.Process()
			if err != nil {
				m.Logger.LogError(c.name, err.Error())
			}
		} else if (tickCount % 59) == 0 {
			if c.config == nil {
				m.Pubsub.PublishStr("config/request/", m.Name)
			}
		}
		tickCount++
		return true
	})

	m.Loop()
}
