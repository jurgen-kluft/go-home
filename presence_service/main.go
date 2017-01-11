package main

import (
	"encoding/json"
	"github.com/jurgen-kluft/go-home/com"
	"github.com/jurgen-kluft/go-home/presence"
	"time"
)

func main() {
	ghChannelName := "Go-Home"
	ghConfigKey := "GO-HOME-PRESENCE-CONFIG"
	ghCom := com.New()

	ghCom.Open()
	ghChannel := ghCom.Subscribe(ghChannelName)

	// Create and initialize
	ghConfigValue, err := ghCom.GetKV(ghConfigKey)
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
		ghCom.SubSend(ghChannel, string(presenceInfo))
	}

	ghCom.Close()
}
