package main

import (
	"fmt"
	"log"

	"github.com/brutella/hc"
	"github.com/jurgen-kluft/go-home/config"
)

func main() {
	var ahkConfig config.AhkConfig

	acsrs := &accessories{}
	accs := acsrs.initializeFromConfig(ahkConfig)

	// configure the ip transport
	config := hc.Config{Pin: ahkConfig.Pin}
	fmt.Println("bridge: " + acsrs.Bridge.Info.Name.GetValue())
	for _, acc := range accs {
		fmt.Println("   accessory: " + acc.Info.Name.GetValue())
	}
	t, err := hc.NewIPTransport(config, acsrs.Bridge.Accessory, accs...)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()

	// Receive complete state of lights, sensors and TV's

}
