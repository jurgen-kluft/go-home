package com

// Com is our generic Com interface
type Com interface {
	Open() error
	Close() error

	Subscribe(topic string) (Channel, error)
	SubSend(channel Channel, message string) error
	SubRecv(channel Channel) (message string, result error)

	SetKV(key, value string) error
	GetKV(key string) (string, error)
}

// Channel is our generic communication channel
type Channel int32

// New will return an instance of Com
func New() Com {
	com := &comRedis{}
	return com
}
