package com

import (
	"fmt"
	"gopkg.in/redis.v5"
)

type comRedis struct {
	redis *redis.Client
}

func (com *comRedis) Open() error {
	com.redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return fmt.Errorf("Not implemented")
}

func (com *comRedis) Close() error {
	return fmt.Errorf("Not implemented")
}

func (com *comRedis) Subscribe(topic string) (Channel, error) {
	return -1, fmt.Errorf("Not implemented")
}

func (com *comRedis) SubSend(channel Channel, message string) error {
	return fmt.Errorf("Not implemented")
}

func (com *comRedis) SubRecv(channel Channel) (message string, result error) {
	return "", fmt.Errorf("Not implemented")
}

func (com *comRedis) SetKV(key string, value string) error {
	return fmt.Errorf("Not implemented")
}

func (com *comRedis) GetKV(key string) (string, error) {
	return "", fmt.Errorf("Not implemented")
}
