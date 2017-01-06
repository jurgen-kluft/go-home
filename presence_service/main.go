package main

import (
	"encoding/json"
	"github.com/jurgen-kluft/go-home/presence"
	"gopkg.in/redis.v5"
	"time"
)

func main() {
	ghChannelName := "Go-Home"
	ghConfigKey := "GO-HOME-PRESENCE-CONFIG"

	// Open REDIS and read all the configurations
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Create and initialize
	ghConfigValue, err := client.Get(ghConfigKey).Result()
	if err != nil {
		panic(err)
	}

	home := presence.Create(ghConfigValue)
	presence := &presence.Presence{}

	// Every N seconds
	ticker := time.NewTicker(time.Second * time.Duration(1.0/home.UpdateFrequency))
	for currentTime := range ticker.C {
		home.Presence(currentTime, presence)
		presenceInfo, _ := json.Marshal(presence)
		client.Publish(ghChannelName, string(presenceInfo)).Result()
	}
}
