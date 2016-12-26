package main

import (
	"gopkg.in/redis.v5"
)

func main() {
	// Open REDIS and read all the configurations
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	presenceConfig, err := redisClient.Get("GO-HOME-PRESENCE-CONFIG").Result()
	if err != nil {
		panic(err)
	}

	gohomeConfigChannelName := "GO-HOME-CONFIG"
	redisPubSub, err := redisClient.Subscribe(gohomeConfigChannelName)
	if err != nil {
		panic(err)
	}
	defer redisPubSub.Close()

	// Block for message on channel
	for true {
		msg, err := redisPubSub.ReceiveMessage()
		if err == nil {
			// Check which config is requested
			break
		}
		switch msg.Payload {
		case "PRESENCE":
			redisClient.Publish(gohomeConfigChannelName, presenceConfig)
		}
	}

}
