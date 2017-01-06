package main

import (
	"encoding/json"
	"fmt"
	"github.com/jurgen-kluft/go-xbase"
	"time"
)

// TimeOfDayConfig contains the user configuration for naming certain periods of the day
type TimeOfDayConfig struct {
	UpdateEvery int `json:"UpdateEvery"`
	Sunset      int
	Sunrise     int
	TimeOfDay   []struct {
		Name  string           `json:"name"`
		Start *xbase.TimeOfDay `json:"start,string"`
		End   *xbase.TimeOfDay `json:"end,string"`
		Wnd   *xbase.TimeOfDay `json:"wnd,string"`
	} `json:"TimeOfDay"`
}

// TimeOfDay is a time-period with a name and a time-range [Start,End]
// The field 'Wnd' is used like: [Start+Wnd, End-Wnd]

// CreateTimeOfDayConfig returns an instance of TimeOfDayConfig by unmarshalling a stream of json bytes
func CreateTimeOfDayConfig(jsondata []byte) (config *TimeOfDayConfig) {
	config = &TimeOfDayConfig{}
	json.Unmarshal(jsondata, config)
	for i := range config.TimeOfDay {
		if config.TimeOfDay[i].Name == "Sunset" {
			config.Sunset = i
		} else if config.TimeOfDay[i].Name == "Sunrise" {
			config.Sunrise = i
		}
	}
	return
}

// Find will return an array of indices that mark elements in the TimeOfDay array that match timeofday.IsBetween
func (t *TimeOfDayConfig) find(hours, minutes, seconds int) (result []int) {
	result = make([]int, 0, 2)
	for i, v := range t.TimeOfDay {
		if xbase.TimeIsBetween(hours, minutes, seconds, v.Start, v.End) {
			result = append(result, i)
		}
	}
	return
}

// Build constructs a JSON message about the current state
func (t *TimeOfDayConfig) Build(currentTime time.Time, latitude float64, longtitude float64) string {

	// Update sunrise and sunset
	sunrise := xbase.CalcSunrise(currentTime, latitude, longtitude)
	t.TimeOfDay[t.Sunrise].Start.Hours = int8(sunrise.Hour())
	t.TimeOfDay[t.Sunrise].Start.Minutes = int8(sunrise.Minute())
	t.TimeOfDay[t.Sunrise].Start.Seconds = int8(sunrise.Second())
	t.TimeOfDay[t.Sunrise].Start.Sub(t.TimeOfDay[t.Sunrise].Wnd)
	t.TimeOfDay[t.Sunrise].End.Hours = int8(sunrise.Hour())
	t.TimeOfDay[t.Sunrise].End.Minutes = int8(sunrise.Minute())
	t.TimeOfDay[t.Sunrise].End.Seconds = int8(sunrise.Second())
	t.TimeOfDay[t.Sunrise].End.Add(t.TimeOfDay[t.Sunrise].Wnd)

	sunset := xbase.CalcSunset(currentTime, latitude, longtitude)
	t.TimeOfDay[t.Sunset].Start.Hours = int8(sunset.Hour())
	t.TimeOfDay[t.Sunset].Start.Minutes = int8(sunset.Minute())
	t.TimeOfDay[t.Sunset].Start.Seconds = int8(sunset.Second())
	t.TimeOfDay[t.Sunset].Start.Sub(t.TimeOfDay[t.Sunset].Wnd)
	t.TimeOfDay[t.Sunset].End.Hours = int8(sunset.Hour())
	t.TimeOfDay[t.Sunset].End.Minutes = int8(sunset.Minute())
	t.TimeOfDay[t.Sunset].End.Seconds = int8(sunset.Second())
	t.TimeOfDay[t.Sunset].End.Add(t.TimeOfDay[t.Sunset].Wnd)

	//
	//     Determine the TimeOfDay elements we are in using current time
	//
	hours, minutes, seconds := currentTime.Clock()
	var strbuffer []byte
	for _, v := range t.TimeOfDay {
		if xbase.TimeIsBetween(hours, minutes, seconds, v.Start, v.End) {
			if len(strbuffer) != 0 {
				strbuffer = append(strbuffer, "&"...)
			}
			strbuffer = append(strbuffer, v.Name...)
		}
	}
	json := fmt.Sprintf("{ timeofday : \"%v\" }", string(strbuffer))
	return json
}
