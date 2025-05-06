package yeel

import (
	"fmt"
	"time"
)

var validProperties = make(map[string]bool, 0)

func init() {
	for _, v := range []string{
		"power", "bright",
		"ct", "rgb", "hue",
		"sat", "color_mode",
		"flowing", "delayoff",
		"flow_params", "music_on",
		"name"} {
		validProperties[v] = true
	}
}

func validateDuration(field string, value time.Duration) error {
	if value < 0 || (value > 0 && value < 50*time.Millisecond) {
		return IllegalArgumentError{field: field, value: fmt.Sprintf("%d", durationMillis(value)), supported: "0,>50"}
	}
	return nil
}

func validateDurationShort(field string, value time.Duration) error {
	if value < 0 || (value > 0 && value < 30*time.Millisecond) {
		return IllegalArgumentError{field: field, value: fmt.Sprintf("%d", durationMillis(value)), supported: "0,>30"}
	}
	return nil
}
