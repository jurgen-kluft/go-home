package main

import (
	"fmt"
	"github.com/jurgen-kluft/go-home/timeofday"
	"github.com/jurgen-kluft/go-xbase"
	"gopkg.in/redis.v5"
	"time"
)

func main() {
	ghChannelName := "Go-Home"
	ghConfigKey := "GO-HOME-TIMEOFDAY-CONFIG"

	// Open REDIS and read all the configurations
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Create and initialize
	ghConfigJSON, err := redisClient.Get(ghConfigKey).Result()
	if err != nil {
		panic(err)
	}
	timeOfDayConfig := timeofday.CreateTimeOfDayConfig([]byte(ghConfigJSON))

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
		redisClient.Publish(ghChannelName, json)
	}
}
