package main

import (
	"fmt"
	"github.com/jurgen-kluft/go-xbase"
	"gopkg.in/redis.v5"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	// ghChannelName := "Go-Home"
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
	timeOfDayConfig := CreateTimeOfDayConfig([]byte(ghConfigJSON))

	sunsetSunriseResponse, _ := http.Get(timeOfDayConfig.URL)
	sunsetSunriseConfigBytes, _ := ioutil.ReadAll(sunsetSunriseResponse.Body)
	ssrc := CreateSunSetSunRiseConfig(sunsetSunriseConfigBytes)
	sunset := &xbase.TimeOfDay{}
	sunset.Parse(ssrc.Results.Sunset)
	sunset.Add(&xbase.TimeOfDay{Hours: 8, Minutes: 0, Seconds: 0})
	sunrise := &xbase.TimeOfDay{}
	sunrise.Parse(ssrc.Results.Sunrise)
	sunrise.Add(&xbase.TimeOfDay{Hours: 8, Minutes: 0, Seconds: 0})
	fmt.Printf("Sunrise: %v, Sunset: %v\n", sunrise, sunset)

	ticker := time.NewTicker(time.Second * time.Duration(timeOfDayConfig.UpdateEvery))
	for currentTime := range ticker.C {
		//
		//     Determine the TimeOfDay elements we are in using current time
		hours, minutes, seconds := currentTime.Clock()
		tods := timeOfDayConfig.Find(hours, minutes, seconds)
		names := make([]string, len(tods))
		for _, index := range tods {
			names = append(names, timeOfDayConfig.TimeOfDay[index].Name)
		}

		json := fmt.Sprintf("{ timeofday : \"%v\" }", strings.Join(names, "&"))
		fmt.Println(json)
		// Send as JSON to REDIS channel 'channel_name', like:
		// { timeofday : "MORNING&BREAKFAST" }
	}
}
