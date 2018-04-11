package main

import (
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
)

// Color Temperature
// URL: https://panasonic.net/es/solution-works/jiyugaoka/

func inTimeSpan(start, end, t time.Time) bool {
	return t.After(start) && t.Before(end)
}

// Return the factor 0.0 - 1.0 that indicates where we are in between start - end
func computeTimeSpanX(start, end, t time.Time) float64 {
	sh, sm, sc := start.Clock()
	sx := float64(sh*60*60) + float64(sm*60) + float64(sc)
	eh, em, ec := end.Clock()
	ex := float64(eh*60*60) + float64(em*60) + float64(ec)
	th, tm, tc := t.Clock()
	tx := float64(th*60*60) + float64(tm*60) + float64(tc)
	x := (tx - sx) / (ex - sx)
	return x
}

type instance struct {
	config  *config.FluxConfig
	suncalc *config.SensorState
	season  *config.Season
	weather *config.SensorState
}

func (s *instance) updateSeasonFromName(season string) {
	for _, e := range s.config.Seasons {
		if e.Name == season {
			s.season = &config.Season{}
			*s.season = e
		}
	}
}

func (s *instance) updateLighttimes() {
	sunmoments := map[string]time.Time{}
	for _, tss := range *s.suncalc.TimeWndAttrs {
		sunmoments[tss.Name] = tss.Begin
	}

	for i, lt := range s.config.Lighttime {
		start := sunmoments[lt.TimeSlot.StartMoment]
		end := sunmoments[lt.TimeSlot.EndMoment]
		s.config.Lighttime[i].TimeSlot.StartTime = start
		s.config.Lighttime[i].TimeSlot.EndTime = end
	}
}

// Process will update 'string'states and 'float'states
// States are both input and output, for example as input
// there are Season/Weather states like 'Season':'Winter'
// and 'Clouds':0.5
func Process(f *instance, client *pubsub.Context) {
	if f.config == nil || f.suncalc == nil || f.season == nil {
		return
	}

	now := time.Now()

	current := config.Lighttime{}
	for _, sm := range f.config.Lighttime {
		t0 := sm.TimeSlot.StartTime
		t1 := sm.TimeSlot.EndTime
		if inTimeSpan(t0, t1, now) {
			current = sm
		}
	}

	// Time interpolation factor, where are we between startMoment - endMoment
	currentx := computeTimeSpanX(current.TimeSlot.StartTime, current.TimeSlot.EndTime, now)
	currentx = float64(int64(currentx*100.0)) / 100.0

	clouds := config.Weather{Clouds: config.MinMax{Min: 0.0, Max: 0.001}, CTPct: 0.0, BriPct: 0.0}
	cloudFac := float64(0.0)
	if f.weather != nil {
		cloudFac = f.weather.GetFloatAttr("clouds", 0.0)
	}
	for _, w := range f.config.Weather {
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
			CT = CT - clouds.CTPct*CT
		}
	}
	CT = f.season.CT.LinearInterpolated(CT)

	// Full cloud cover will increase brightness by 10% of (Max - Current)
	// BRI = 0 -> Very dim light
	// BRI = 1 -> Very bright light
	BRI := current.Bri.LinearInterpolated(currentx)
	BRI = BRI + cloudFac*0.1*(1.0-BRI)
	if current.Darkorlight != "dark" {
		// A bit brighter lights when there are clouds during the day.
		if clouds.CTPct >= 0 {
			BRI = BRI + clouds.BriPct*(1.0-BRI)
		} else {
			BRI = BRI - clouds.BriPct*BRI
		}
	}
	BRI = f.season.BRI.LinearInterpolated(BRI)

	// Publishing the following sensors:
	//  - Sensor.Light.HUE, Name = CT, Value = float64(100.0)
	//  - Sensor.Light.HUE, Name = BRI, Value = float64(100.0)
	//  - Sensor.Light.YEE, Name = CT, Value = float64(100.0)
	//  - Sensor.Light.YEE, Name = BRI, Value = float64(100.0)
	//  - Sensor.Light.DarkOrLight = string(Dark)

	for _, ltype := range f.config.Lighttype {
		sensor := config.NewSensorState("light." + ltype.Name)

		lct := ltype.CT.LinearInterpolated(CT)
		sensor.AddFloatAttr("CT", lct)

		lbri := ltype.BRI.LinearInterpolated(BRI)
		sensor.AddFloatAttr("BRI", lbri)

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
	client.Publish(channel, sensorjson)
}

func main() {
	flux := &instance{}

	logger := logpkg.New("flux")
	logger.AddEntry("emitter")
	logger.AddEntry("flux")

	for {
		client := pubsub.New(config.EmitterSecrets["host"])
		register := []string{"config/flux/", "state/sensor/weather/", "state/sensor/sun/", "state/sensor/season/", "state/light/hue/", "state/light/yee/"}
		subscribe := []string{"config/flux/", "state/sensor/weather/", "state/sensor/sun/", "state/sensor/season/"}
		err := client.Connect("flux", register, subscribe)
		if err == nil {
			logger.LogInfo("emitter", "connected")

			connected := true
			for connected {
				select {
				case msg := <-client.InMsgs:
					topic := msg.Topic()
					if topic == "config/flux/" {
						logger.LogInfo("flux", "received configuration")
						flux.config, err = config.FluxConfigFromJSON(string(msg.Payload()))
						if err != nil {
							logger.LogError("flux", err.Error())
						}
					} else if topic == "state/sensor/weather/" {
						flux.weather, err = config.SensorStateFromJSON(string(msg.Payload()))
						if err != nil {
							logger.LogError("flux", err.Error())
						}
					} else if topic == "state/sensor/sun/" {
						flux.suncalc, err = config.SensorStateFromJSON(string(msg.Payload()))
						if err != nil {
							logger.LogError("flux", err.Error())
						}
					} else if topic == "state/sensor/season/" {
						flux.updateSeasonFromName(string(msg.Payload()))
					} else if topic == "client/disconnected/" {
						logger.LogInfo("emitter", "disconnected")
						connected = false
					}
				case <-time.After(time.Second * 10):
					Process(flux, client)
				}
			}
		}

		if err != nil {
			logger.LogError("flux", err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}
