package main

import (
	"encoding/json"
	"github.com/jurgen-kluft/go-xbase"
)

// Config is a JSON configuration like this:
/*
{
    "UpdateEvery": "300",
    "TimeOfDay": [
        {
            "name": "BREAKFAST",
            "start": 7:00,
            "end": 9:30
        },
        {
            "name": "MORNING",
            "start": 6:00,
            "end": 12:00
        },
        {
            "name": "NOON",
            "start": 12:00,
            "end": 13:00
        },
        {
            "name": "LUNCH",
            "start": 11:45,
            "end": 12:45
        },
        {
            "name": "AFTERNOON",
            "start": 13:00,
            "end": 18:00
        },
        {
            "name": "DINNER",
            "start": 18:00,
            "end": 20:00
        },
        {
            "name": "NIGHT",
            "start": 20:00,
            "end": 6:00
        },
        {
            "name": "SLEEPING",
            "start": 22:00,
            "end": 6:00
        },
        {
            "name": "EVENING",
            "start": 16:30,
            "end": 22:00
        }
    ]
}
*/

// TimeOfDayConfig contains the user configuration for naming certain periods of the day
type TimeOfDayConfig struct {
	UpdateEvery int `json:"UpdateEvery"`
	TimeOfDay   []struct {
		Name  string           `json:"name"`
		Start *xbase.TimeOfDay `json:"start,string"`
		End   *xbase.TimeOfDay `json:"end,string"`
	} `json:"TimeOfDay"`
}

// CreateTimeOfDayConfig returns an instance of TimeOfDayConfig by unmarshalling a stream of json bytes
func CreateTimeOfDayConfig(jsondata []byte) (config *TimeOfDayConfig) {
	config = &TimeOfDayConfig{}
	json.Unmarshal(jsondata, config)
	return
}

// Find will return an array of indices that mark elements in the TimeOfDay array that match timeofday.IsBetween
func (t *TimeOfDayConfig) Find(hours, minutes, seconds int) (result []int) {
	result = make([]int, 0, 2)
	for i, v := range t.TimeOfDay {
		if xbase.TimeIsBetween(hours, minutes, seconds, v.Start, v.End) {
			result = append(result, i)
		}
	}
	return
}
