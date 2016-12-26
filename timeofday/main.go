package main

import (
	"fmt"
	"gopkg.in/redis.v5"
	"io/ioutil"
	"net/http"
)

func main() {
	// Open REDIS and read all the configurations
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Create and initialize
	msgTimeOfDayConfigJSON, err := redisClient.Get("GO-HOME-TIMEOFDAY-CONFIG").Result()
	if err != nil {
		panic(err)
	}
	timeOfDayConfig := CreateTimeOfDayConfig([]byte(msgTimeOfDayConfigJSON))

	// timeofday_config_key := "Go-Home-TimeOfDay-Config"
	// timeofday_state_key := "Go-Home-TimeOfDay-State"

	// get TimeOfDayConfig from REDIS
	// Unmarshal the JSON data into our Golang struct

	sunsetSunriseResponse, _ := http.Get(timeOfDayConfig.URL)
	sunsetSunriseConfigBytes, _ := ioutil.ReadAll(sunsetSunriseResponse.Body)
	ssrc := CreateSunSetSunRiseConfig(sunsetSunriseConfigBytes)
	fmt.Printf("SunSet: %v, SunRise: %v\n", ssrc.Results.Sunset, ssrc.Results.Sunrise)

	// channel_name := "Go-Home"

	// Every N seconds
	//
	//     Determine the TimeOfDay elements we are in using current time
	//     Send as JSON to REDIS channel 'channel_name', like:
	//
	//          ---- TIMEOFDAY
	//          timeofday : Morning & Breakfast

}
