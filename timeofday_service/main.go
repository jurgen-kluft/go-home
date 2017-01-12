package main

import (
	"fmt"
	"github.com/jurgen-kluft/go-home/com"
	"github.com/jurgen-kluft/go-home/timeofday"
	"github.com/jurgen-kluft/go-xbase"
	"time"
)

func main() {
	ghChannelName := "Go-Home"
	ghConfigKey := "GO-HOME-TIMEOFDAY-CONFIG"
	ghCom := com.New()

	ghCom.Open()
	ghChannel := ghCom.Subscribe(ghChannelName)

	// Create and initialize
	ghConfigValue, err := ghCom.GetKV(ghConfigKey)
	if err != nil {
		panic(err)
	}
	timeOfDayConfig := timeofday.CreateTimeOfDayConfig([]byte(ghConfigValue))

	// Shanghai latitude and longtitude
	localLatitude := 31.2222200
	localLongtitude := -121.4580600

	sunrise := xbase.CalcSunrise(time.Now(), localLatitude, localLongtitude)
	sunset := xbase.CalcSunset(time.Now(), localLatitude, localLongtitude)
	fmt.Printf("Sunrise: %v, Sunset: %v\n", sunrise, sunset)

	ticker := time.NewTicker(time.Second * time.Duration(timeOfDayConfig.UpdateEvery))
	for currentTime := range ticker.C {
		//
		//     Determine the TimeOfDay elements we are in using current time
		json := timeOfDayConfig.Build(currentTime, localLatitude, localLongtitude)

		// Send as JSON to REDIS channel 'channel_name', like:
		// { timeofday : "MORNING&BREAKFAST" }
		ghCom.SubSend(ghChannel, json)
	}

	ghCom.Close()
}
