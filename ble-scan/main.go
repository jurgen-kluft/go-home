package main

import (
	"fmt"
	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
	"log"
)

func onStateChanged(device gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning for iBeacon Broadcasts...")
		device.Scan([]gatt.UUID{}, true)
		return
	default:
		device.StopScanning()
	}
}

func main() {
	bluetooth := Create()
	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	device.Handle(gatt.PeripheralDiscovered(func(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
		if !bluetooth.HaveSeenBeacon(a.ManufacturerData) {
			b, err := bluetooth.NewBeacon(a.ManufacturerData)
			if err == nil {
				b.RSSI = ((rssi / 10) * 10)
				fmt.Println("UUID: ", b.UUID)
				fmt.Println("Major: ", b.Major)
				fmt.Println("Minor: ", b.Minor)
				fmt.Println("RSSI: ", rssi)
			}
		} else {
			b := bluetooth.GetExistingBeacon(a.ManufacturerData)
			if b.RSSI != ((rssi / 10) * 10) {
				b.RSSI = ((rssi / 10) * 10)
				fmt.Println("UUID: ", b.UUID)
				fmt.Println("Major: ", b.Major)
				fmt.Println("Minor: ", b.Minor)
				fmt.Println("RSSI: ", rssi)
			}
		}
	}))

	device.Init(onStateChanged)
	select {}
}
