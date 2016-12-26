package main

import (
	"encoding/json"
)

// SHANGHAI :
// - Lattitude = 31.2222200
// - Longitude = 121.4580600

// http://api.sunrise-sunset.org/json?lat=31.2222200&lng=121.4580600

// Time is in UTC, response:
/*
{
	"results": {
		"sunrise": "10:45:01 PM",
		"sunset": "8:53:12 AM",
		"solar_noon": "3:49:06 AM",
		"day_length": "10:08:11",
		"civil_twilight_begin": "10:18:24 PM",
		"civil_twilight_end": "9:19:48 AM",
		"nautical_twilight_begin": "9:48:10 PM",
		"nautical_twilight_end": "9:50:03 AM",
		"astronomical_twilight_begin": "9:18:35 PM",
		"astronomical_twilight_end": "10:19:38 AM"
	},
	"status": "OK"
}
*/

// SunSetSunRiseConfig contains the sun-rise/sun-set information of a specific world position
type SunSetSunRiseConfig struct {
	Results struct {
		Sunrise                   string `json:"sunrise"`
		Sunset                    string `json:"sunset"`
		SolarNoon                 string `json:"solar_noon"`
		DayLength                 string `json:"day_length"`
		CivilTwilightBegin        string `json:"civil_twilight_begin"`
		CivilTwilightEnd          string `json:"civil_twilight_end"`
		NauticalTwilightBegin     string `json:"nautical_twilight_begin"`
		NauticalTwilightEnd       string `json:"nautical_twilight_end"`
		AstronomicalTwilightBegin string `json:"astronomical_twilight_begin"`
		AstronomicalTwilightEnd   string `json:"astronomical_twilight_end"`
	} `json:"results"`
	Status string `json:"status"`
}

// CreateSunSetSunRiseConfig returns an instance of SunSetSunRiseConfig by unmarshalling a stream of json bytes
func CreateSunSetSunRiseConfig(jsondata []byte) (config *SunSetSunRiseConfig) {
	json.Unmarshal(jsondata, config)
	return
}

// Config is a JSON configuration like this:
/*
{
    "URL": "http://api.sunrise-sunset.org/json?lat=31.2222200&lng=121.4580600",
    "UpdateEvery": "300",
    "TimeOfDay": [
        {
            "name": "SUN",
            "start": "$(SUNRISE)",
            "end": "$(SUNSET)"
        },
        {
            "name": "BREAKFAST",
            "start": "7:00",
            "end": "9:30"
        },
        {
            "name": "MORNING",
            "start": "6:00",
            "end": "12:00"
        },
        {
            "name": "NOON",
            "start": "12:00",
            "end": "13:00"
        },
        {
            "name": "LUNCH",
            "start": "11:45",
            "end": "12:45"
        },
        {
            "name": "AFTERNOON",
            "start": "13:00",
            "end": "18:00"
        },
        {
            "name": "DINNER",
            "start": "18:00",
            "end": "20:00"
        },
        {
            "name": "NIGHT",
            "start": "20:00",
            "end": "6:00"
        },
        {
            "name": "SLEEPING",
            "start": "22:00",
            "end": "6:00"
        },
        {
            "name": "EVENING",
            "start": "$(SUN.start) - 0:30",
            "end": "$(SLEEPING.start)"
        }
    ]
}
*/

// TimeOfDayConfig contains the user configuration for naming certain periods of the day
type TimeOfDayConfig struct {
	URL         string `json:"URL"`
	UpdateEvery string `json:"UpdateEvery"`
	TimeOfDay   []struct {
		Name  string `json:"name"`
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"TimeOfDay"`
}

// CreateTimeOfDayConfig returns an instance of TimeOfDayConfig by unmarshalling a stream of json bytes
func CreateTimeOfDayConfig(jsondata []byte) (config *TimeOfDayConfig) {
	json.Unmarshal(jsondata, config)
	return
}
