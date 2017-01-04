package config

import (
	"gopkg.in/redis.v5"
)

const (
	ghPresenceConfigKey   = "GO-HOME-PRESENCE-CONFIG"
	ghPresenceConfigValue = `{
	"host": "192.168.1.3",
	"port": "5000",
	"user": "admin",
	"password": "*********",
	"detectionwnd": "200",
	"detectionfreq": "60",
	"devices": [{
		"name": "Faith",
		"mac": "D8:1D:72:97:51:94"
	}, {
		"name": "Jurgen",
		"mac": "90:3C:92:72:0D:C8"
	}, {
		"name": "GrandPa",
		"mac": "BC:44:86:5A:CD:D4"
	}, {
		"name": "GrandMa",
		"mac": "94:94:26:B5:E6:1C"
	}]}`

	ghTimeOfDayConfigKey   = "GO-HOME-TIMEOFDAY-CONFIG"
	ghTimeOfDayConfigValue = `{
    "UpdateEvery": "300",
    "TimeOfDay": [
        {
            "name": "BREAKFAST",
            "start": 7:00,
            "end": 9:30
        },
        {
            "name": "MORNING",
            "start": 6:00,
            "end": 12:00
        },
        {
            "name": "NOON",
            "start": 12:00,
            "end": 13:00
        },
        {
            "name": "LUNCH",
            "start": 11:45,
            "end": 12:45
        },
        {
            "name": "AFTERNOON",
            "start": 13:00,
            "end": 18:00
        },
        {
            "name": "DINNER",
            "start": 18:00,
            "end": 20:00
        },
        {
            "name": "NIGHT",
            "start": 20:00,
            "end": 6:00
        },
        {
            "name": "SLEEPING",
            "start": 22:00,
            "end": 6:00
        },
        {
            "name": "EVENING",
            "start": 16:30,
            "end": 22:00
        }
    ]}`

	ghEventLogConfigKey   = "GO-HOME-EVENTLOG-CONFIG"
	ghEventLogConfigValue = `{
	  "Filename": "go-home-%.log",
	  "Folder": "~/go-home/log/",
	  "Rolling": true,
	  "MaximumFileSize": 10485760,
	  "MaximumLogs": 10
	}`
)

// WriteConfigsToRedis will set a couple of KEY,VALUE pairs that act as configuration
func WriteConfigsToRedis(URL string, password string, db int) {
	// Open REDIS and read all the configurations
	redisClient := redis.NewClient(&redis.Options{
		Addr:     URL,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	redisClient.Set(ghPresenceConfigKey, ghPresenceConfigValue, 0)
	redisClient.Set(ghTimeOfDayConfigKey, ghTimeOfDayConfigValue, 0)
}
