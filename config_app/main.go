package main

import (
	"github.com/jurgen-kluft/go-home/config"
)

func main() {
	WriteConfigsToRedis("localhost:6379", "", 0)
}
