package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
)

var (
	device = flag.String("device", "default", "implementation of ble")
	du     = flag.Duration("du", 15*time.Second, "scanning duration")
)

type beacon struct {
	Address string
	RSSI    int
}

func intAbs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func main() {
	flag.Parse()

	d, err := dev.NewDevice(*device)
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	beacons := make(map[string]*beacon)
	advFilter := func(a ble.Advertisement) bool {
		b, exists := beacons[a.Address().String()]
		if !exists {
			b := &beacon{Address: a.Address().String(), RSSI: a.RSSI()}
			beacons[a.Address().String()] = b
			return true
		}

		if intAbs(b.RSSI-a.RSSI()) > 20 {
			b.RSSI = (b.RSSI + a.RSSI()) / 2 // Take average
			return true
		}

		return false
	}

	// Scan for specified durantion, or until interrupted by user.
	fmt.Printf("Scanning for %s...\n", *du)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.Scan(ctx, true, advHandler, advFilter))
}

func advHandler(a ble.Advertisement) {
	if a.Connectable() {
		fmt.Printf("[%s] C %3d:", a.Address(), a.RSSI())
	} else {
		fmt.Printf("[%s] N %3d:", a.Address(), a.RSSI())
	}
	comma := ""
	if len(a.LocalName()) > 0 {
		fmt.Printf(" Name: %s", a.LocalName())
		comma = ","
	}
	if len(a.Services()) > 0 {
		fmt.Printf("%s Svcs: %v", comma, a.Services())
		comma = ","
	}
	if len(a.ManufacturerData()) > 0 {
		fmt.Printf("%s MD: %X", comma, a.ManufacturerData())
	}
	fmt.Printf("\n")
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}
