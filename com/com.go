package com

// Com is our generic Communication interface
type Com interface {
	Open() error
	Close() error

	SubscribeToTopic(topic string) (Topic, error)
	SendToTopic(t Topic, message string) error
	RecvFromTopic(t Topic) (message string, result error)

	// Some helper
	SetKV(key, value string) error
	GetKV(key string) (string, error)
}

// Topic is our communication group
type Topic int32

// New will return an instance of Com
func New() Com {
	com := &comRedis{}
	return com
}
