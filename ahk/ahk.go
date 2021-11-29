package main

import (
	"fmt"
	"io/ioutil"

	"github.com/brutella/hc"
	"github.com/jurgen-kluft/go-home/config"
)

func main() {
	filename := "../config/ahk.config.json"
	filedata, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error)
		return
	}
	jsonbytes := filedata

	var ahkConfig *config.AhkConfig
	ahkConfig, err = config.AhkConfigFromJSON(jsonbytes)

	acsrs := &accessories{}
	accs := acsrs.initializeFromConfig(ahkConfig)

	// configure the ip transport
	config := hc.Config{Pin: ahkConfig.Pin}
	fmt.Println("bridge: " + acsrs.Bridge.Info.Name.GetValue())
	for _, acc := range accs {
		fmt.Println("   accessory: " + acc.Info.Name.GetValue() + " type:" + accessoryTypeToString(acc.Type))
	}
	t, err := hc.NewIPTransport(config, acsrs.Bridge.Accessory, accs...)
	if err != nil {
		fmt.Println(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()

	// Receive complete state of lights, sensors and TV's
	// Somewhere we should hold a single source of truth, the full state of
	// all lights, TV's, booleans and other variables.

	return
}
