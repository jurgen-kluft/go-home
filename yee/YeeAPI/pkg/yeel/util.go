package yeel

import "time"

func round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}

func scale(min, max int, val float64) int {
	scaled := val*(float64(max)-float64(min)) + float64(min)
	return round(scaled)

}

func durationMillis(d time.Duration) int {
	return int(d / time.Millisecond)
}
