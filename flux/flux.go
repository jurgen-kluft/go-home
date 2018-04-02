package flux

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jurgen-kluft/go-home/config"
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
	suncalc *config.SuncalcState
	season  *config.Season
	clouds  *config.SensorState
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
	for _, sm := range s.suncalc.Moments {
		sunmoments[sm.Name] = sm.Begin
	}

	for i, lt := range s.config.Lighttime {
		start := sunmoments[lt.StartMoment]
		end := sunmoments[lt.EndMoment]
		s.config.Lighttime[i].Start = start
		s.config.Lighttime[i].End = end
	}
}

// Process will update 'string'states and 'float'states
// States are both input and output, for example as input
// there are Season/Weather states like 'Season':'Winter'
// and 'Clouds':0.5
func Process(f *instance, client *pubsub.Context) {
	if f.config == nil || f.suncalc == nil || f.season == nil || f.clouds == nil {
		return
	}

	now := time.Now()

	current := config.Lighttime{}
	for _, sm := range f.config.Lighttime {
		t0 := sm.Start
		t1 := sm.End
		if inTimeSpan(t0, t1, now) {
			current = sm
		}
	}

	// Time interpolation factor, where are we between startMoment - endMoment
	currentx := computeTimeSpanX(current.Start, current.End, now)
	currentx = float64(int64(currentx*100.0)) / 100.0

	clouds := config.Weather{Clouds: config.MinMax{0.0, 0.001}, CTPct: 0.0, BriPct: 0.0}
	cloudFac, err := strconv.ParseFloat(f.clouds.Value, 64)
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
	CT := current.CT[0] + currentx*(current.CT[1]-current.CT[0])
	if current.Darkorlight != "dark" {
		if clouds.CTPct >= 0 {
			CT = CT + clouds.CTPct*(1.0-CT)
		} else {
			CT = CT - clouds.CTPct*CT
		}
	}
	CT = f.season.CT.Min + (CT * (f.season.CT.Max - f.season.CT.Min))

	// Full cloud cover will increase brightness by 10% of (Max - Current)
	// BRI = 0 -> Very dim light
	// BRI = 1 -> Very bright light
	BRI := current.Bri[0] + currentx*(current.Bri[1]-current.Bri[0])
	BRI = BRI + cloudFac*0.1*(1.0-BRI)
	if current.Darkorlight != "dark" {
		// A bit brighter lights when there are clouds during the day.
		if clouds.CTPct >= 0 {
			BRI = BRI + clouds.BriPct*(1.0-BRI)
		} else {
			BRI = BRI - clouds.BriPct*BRI
		}
	}
	BRI = f.season.BRI.Min + (BRI * (f.season.BRI.Max - f.season.BRI.Min))

	// Publishing the following sensors:
	//  - Sensor.Light.HUE_CT = float64(100.0)
	//  - Sensor.Light.HUE_BRI = float64(100.0)
	//  - Sensor.Light.YEE_CT = float64(100.0)
	//  - Sensor.Light.YEE_BRI = float64(100.0)
	//  - Sensor.Light.DarkOrLight = string(Dark)

	for _, ltype := range f.config.Lighttype {
		lct := ltype.CT.Min + CT*(ltype.CT.Max-ltype.CT.Min)
		sensorCT := config.SensorState{Domain: "sensor", Product: "light", Name: ltype.Name + "_CT", Type: "float", Value: fmt.Sprintf("%f", lct), Time: time.Now()}
		publishSensor(sensorCT, client)
		lbri := ltype.BRI.Min + BRI*(ltype.BRI.Max-ltype.BRI.Min)
		sensorBRI := config.SensorState{Domain: "sensor", Product: "light", Name: ltype.Name + "_BRI", Type: "float", Value: fmt.Sprintf("%f", lbri), Time: time.Now()}
		publishSensor(sensorBRI, client)

		//	json += fmt.Sprintf("\"type\": %s, ", ltype.Name)
		//	json += fmt.Sprintf("\"ct\": %f, ", math.Floor(lct))
		//	json += fmt.Sprintf("\"bri\": %f ", math.Floor(lbri))
	}

	sensorDOL := config.SensorState{Domain: "sensor", Product: "light", Name: "DarkOrLight", Type: "string", Value: string(current.Darkorlight), Time: time.Now()}
	publishSensor(sensorDOL, client)
}

func publishSensor(sensor config.SensorState, client *pubsub.Context) {
	data, err := json.Marshal(sensor)
	if err == nil {
		jsonstr := string(data)
		client.Publish(fmt.Sprintf("%s/%s/%s", sensor.Domain, sensor.Product, sensor.Name), jsonstr)
	}
}

func main() {
	flux := &instance{}
	for {
		client := pubsub.New()
		err := client.Connect("flux")
		if err == nil {
			for {
				select {
				case msg := <-client.InMsgs:
					if msg.Topic() == "config/flux" {
						if flux.config == nil {
							flux.config, err = config.FluxConfigFromJSON(string(msg.Payload()))
						}
					} else if msg.Topic() == "sensor/weather/clouds" {
						flux.clouds, err = config.SensorStateFromJSON(string(msg.Payload()))
					} else if msg.Topic() == "sensor//clouds" {
						flux.suncalc, err = NewSuncalc(msg.Data)
					} else if msg.Topic() == "sensor/weather/clouds" {
						flux.updateSeasonFromName(msg.Data)
					}
					break
				case <-time.After(time.Second * 10):
					// do something if messages are taking too long
					// or if we haven't received enough state info.
					Process(flux, client)
					break
				}
			}
		}

		// Wait for 10 seconds before retrying
		time.Sleep(10 * time.Second)
	}
}
