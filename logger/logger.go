package main

import (
	"fmt"
	"time"

	microservice "github.com/jurgen-kluft/go-home/micro-service"
)

type context struct {
	serviceName            string
	thisConfigFilepath     string
	servicesConfigFilepath string
	service                *microservice.Service
}

func new() *context {
	c := &context{}
	c.serviceName = "logger"
	c.thisConfigFilepath = "config/logger.config.json"
	c.servicesConfigFilepath = "config/services.config.json/"
	return c
}

func main() {
	c := new()

	var err error
	c.service, err = microservice.New(c.serviceName, c.thisConfigFilepath, c.servicesConfigFilepath, time.Second*30)
	if err != nil {
		fmt.Println("Error creating microservice:", err)
		return
	}
	c.service.Connect()

	tickCount := 0

	c.service.RegisterHandler("config", func(m *microservice.Service, msg *microservice.Message) bool {
		// handle new config
		return true
	})

	c.service.RegisterHandler("weather", func(m *microservice.Service, msg *microservice.Message) bool {
		// write weather message to log
		return true
	})

	c.service.RegisterHandler("sun", func(m *microservice.Service, msg *microservice.Message) bool {
		// write sun message to log
		return true
	})

	c.service.RegisterHandler("calendar", func(m *microservice.Service, msg *microservice.Message) bool {
		// write calendar message to log
		return true
	})

	c.service.RegisterHandler("tick", func(m *microservice.Service, msg *microservice.Message) bool {
		if (tickCount % 30) == 0 {
		}
		tickCount++
		return true
	})

	c.service.Loop()
}
