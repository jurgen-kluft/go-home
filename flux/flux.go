package flux

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jurgen-kluft/go-home/suncalc"
	"github.com/nanopack/mist/clients"
)

// Color Temperature
// URL: https://panasonic.net/es/solution-works/jiyugaoka/

func NewConfig(json string) (*Config, error) {
	data := []byte(json)
	l, err := UnmarshalConfig(data)
	return l, err
}

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
	config  *Config
	suncalc *suncalc.State
	season  *Season
	weather *WeatherState
}

func (s *instance) updateSeasonFromName(season string) {
	for _, e := range s.config.Seasons {
		if e.Name == season {
			s.season = &Season{}
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
func Process(f *instance, client *clients.TCP) {
	if f.config == nil || f.suncalc == nil || f.season == nil || f.weather == nil {
		return
	}

	now := time.Now()

	current := Lighttime{}
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

	clouds := Weather{Clouds: MinMax{0.0, 0.001}, CTPct: 0.0, BriPct: 0.0}
	cloudFac := f.weather.current.Clouds
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
		sensorCT := SensorState{Domain: "sensor", Product: "light", Name: ltype.Name + "_CT", Type: "float", Value: fmt.Sprintf("%f", lct), Time: time.Now()}
		publishSensor(sensorCT, client)
		lbri := ltype.BRI.Min + BRI*(ltype.BRI.Max-ltype.BRI.Min)
		sensorBRI := SensorState{Domain: "sensor", Product: "light", Name: ltype.Name + "_BRI", Type: "float", Value: fmt.Sprintf("%f", lbri), Time: time.Now()}
		publishSensor(sensorBRI, client)

		//	json += fmt.Sprintf("\"type\": %s, ", ltype.Name)
		//	json += fmt.Sprintf("\"ct\": %f, ", math.Floor(lct))
		//	json += fmt.Sprintf("\"bri\": %f ", math.Floor(lbri))
	}

	sensorDOL := SensorState{Domain: "sensor", Product: "light", Name: "DarkOrLight", Type: "string", Value: string(current.Darkorlight), Time: time.Now()}
	publishSensor(sensorDOL, client)
}

func publishSensor(sensor SensorState, client *clients.TCP) {
	data, err := json.Marshal(sensor)
	if err == nil {
		jsonstr := string(data)
		client.Publish([]string{"sensor", "light"}, jsonstr)
	}
}

type SensorState struct {
	Domain  string    `json:"domain"`
	Product string    `json:"product"`
	Name    string    `json:"name"`
	Type    string    `json:"type"`
	Value   string    `json:"value"`
	Time    time.Time `json:"time"`
}

func tagsContains(tag string, tags []string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func main() {
	flux := &instance{}
	for {
		client, err := clients.New("127.0.0.1:1445", "authtoken.wicked")
		if err != nil {
			fmt.Println(err)
			continue
		}

		client.Ping()

		// 'flux' OR (('suncalc' OR 'weather') AND 'state')
		client.Subscribe([]string{"suncalc weather | state & flux |"})
		client.Publish([]string{"request", "config", "weather", "suncalc", "season"}, "flux")

		client.ListAll()

		for {
			select {
			case msg := <-client.Messages():
				if tagsContains("config", msg.Tags) {
					if flux.config == nil {
						flux.config, err = NewConfig(msg.Data)
					}
				} else if tagsContains("state", msg.Tags) {
					if tagsContains("weather", msg.Tags) {
						flux.weather, err = NewWeatherState(msg.Data)
					} else if tagsContains("suncalc", msg.Tags) {
						flux.suncalc, err = NewSuncalc(msg.Data)
					} else if tagsContains("season", msg.Tags) {
						flux.updateSeasonFromName(msg.Data)
					}
				}
				break
			case <-time.After(time.Second * 10):
				// do something if messages are taking too long
				// or if we haven't received enough state info.
				Process(flux, client)
				break
			}
		}

		// Disconnect from Mist
	}
}
