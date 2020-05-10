package main

import (
	"fmt"
	"math"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/metrics"
	"github.com/jurgen-kluft/go-home/pubsub"
)

// Color Temperature
// URL: https://panasonic.net/es/solution-works/jiyugaoka/

func inTimeSpan(start, end, t time.Time) bool {
	sh, sm, sc := start.Clock()
	sx := float64(sh*60*60) + float64(sm*60) + float64(sc)
	eh, em, ec := end.Clock()
	ex := float64(eh*60*60) + float64(em*60) + float64(ec) + 1.0
	th, tm, tc := t.Clock()
	tx := float64(th*60*60) + float64(tm*60) + float64(tc)
	return tx >= sx && tx < ex
}

// Return the factor 0.0 - 1.0 that indicates where we are in between start - end
func computeTimeSpanX(start, end, t time.Time) float64 {
	sh, sm, sc := start.Clock()
	sx := float64(sh*60*60) + float64(sm*60) + float64(sc)
	eh, em, ec := end.Clock()
	ex := float64(eh*60*60) + float64(em*60) + float64(ec) + 1.0
	th, tm, tc := t.Clock()
	tx := float64(th*60*60) + float64(tm*60) + float64(tc)
	x := (tx - sx) / (ex - sx)
	return x
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
}

func new() *context {
	c := &context{}
	c.name = c.name
	c.cfgch = "config/flux/"
	c.metrics, _ = metrics.New()
	c.metrics.Register("hue", map[string]string{"CT": "Color Temperature", "BRI": "Brightness"}, map[string]interface{}{"CT": 200.0, "BRI": 200.0})
	c.metrics.Register("yee", map[string]string{"CT": "Color Temperature", "BRI": "Brightness"}, map[string]interface{}{"CT": 200.0, "BRI": 200.0})
	c.seasonName = "spring"
	c.averageCT = NewFilter(30)
	c.averageBRI = NewFilter(30)
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
func (c *context) Process(client *pubsub.Context) {
	if c.config == nil || c.suncalc == nil {
		return
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
	current := config.Lighttime{}
	for _, sm := range c.config.Lighttime {
		t0 := sm.TimeSlot.StartTime
		t1 := sm.TimeSlot.EndTime
		if inTimeSpan(t0, t1, now) {
			fmt.Println("Current light time from", sm.TimeSlot.StartMoment, "to", sm.TimeSlot.EndMoment)
			current = sm
			break
		}
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
		sensor := config.NewSensorState("all")

		c.metrics.Begin(ltype.Name)

		lct := ltype.CT.LinearInterpolated(CT)
		sensor.AddFloatAttr("CT", math.Floor(lct))
		c.metrics.Set(ltype.Name, "CT", lct)

		lbri := ltype.BRI.LinearInterpolated(BRI)
		sensor.AddFloatAttr("BRI", math.Floor(lbri))
		c.metrics.Set(ltype.Name, "BRI", lbri)

		c.metrics.Send(ltype.Name)

		jsonstr, err := sensor.ToJSON()
		if err == nil {
			publishSensor(fmt.Sprintf("state/sensor/%s/", ltype.Name), jsonstr, client)
		}
	}

	sensorDOL, err := config.StringAttrAsJSON("darkorlight", "DarkOrLight", string(current.Darkorlight))
	if err == nil {
		publishSensor("state/sensor/darkorlight/", sensorDOL, client)
	}
}

func publishSensor(channel string, sensorjson string, client *pubsub.Context) {
	fmt.Println("Publish at", channel, "JSON [", sensorjson, "]")
	client.Publish(channel, sensorjson)
}

func main() {
	c := new()

	logger := logpkg.New(c.name)
	logger.AddEntry("emitter")
	logger.AddEntry(c.name)

	for {
		client := pubsub.New(config.PubSubCfg)
		register := []string{"config/flux/", "state/sensor/weather/", "state/sensor/sun/", "state/sensor/season/", "state/light/hue/", "state/light/yee/"}
		subscribe := []string{"config/flux/", "state/sensor/weather/", "state/sensor/sun/", "state/sensor/season/", "config/request/"}
		err := client.Connect(c.name, register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == c.cfgch {
						c.config, err = config.FluxConfigFromJSON(msg.Payload())
						if err == nil {
							logger.LogInfo(c.name, "received configuration")
						} else {
							logger.LogError(c.name, err.Error())
						}
					} else if topic == "state/sensor/weather/" {
						c.weather, err = config.SensorStateFromJSON(msg.Payload())
						if err == nil {
							logger.LogInfo(c.name, "received weather state")
						} else {
							logger.LogError(c.name, err.Error())
						}
					} else if topic == "state/sensor/sun/" {
						c.suncalc, err = config.SensorStateFromJSON(msg.Payload())
						if err == nil {
							logger.LogInfo(c.name, "received sun state")
						} else {
							logger.LogError(c.name, err.Error())
						}
					} else if topic == "state/sensor/season/" {
						seasonSensorState, err := config.SensorStateFromJSON(msg.Payload())
						if err == nil {
							logger.LogInfo(c.name, "received season state")
							c.seasonName = seasonSensorState.GetValueAttr("season", "winter")
						} else {
							logger.LogError(c.name, err.Error())
						}
					} else if topic == "client/disconnected/" {
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}
				case <-time.After(time.Second * 10):
					c.Process(client)

				case <-time.After(time.Minute * 1): // Try and request our configuration
					if c.config == nil {
						client.Publish("config/request/", "flux")
					}

				}

			}
		}

		if err != nil {
			logger.LogError(c.name, err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}
