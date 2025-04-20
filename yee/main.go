package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"log"
)

func main() {
	yeeBridge := accessory.NewBridge(accessory.Info{Name: "YeeBridge", ID: 1})
	yeeLight1 := NewYeelight(accessory.Info{Name: "Frontdoor Entrance", ID: 2, Manufacturer: "YeeLight"}, "10.0.0.57", "55443")
	yeeLight2 := NewYeelight(accessory.Info{Name: "Piano Light", ID: 3, Manufacturer: "YeeLight"}, "10.0.0.54", "55443")

	t, err := hc.NewIPTransport(hc.Config{Pin: "12341234"}, yeeBridge.Accessory, yeeLight1.Accessory, yeeLight2.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {
		t.Stop()
	})

	t.Start()
}
