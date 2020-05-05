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
	du     = flag.Duration("du", 3000*time.Second, "scanning duration")
)

type beacon struct {
	Address  string
	Name     string
	RSSI     int
	Services map[string]ble.UUID
}

var knownBeacons = map[string]struct{ Name string }{
	"cc:98:8b:d1:4a:0f": {Name: "Sony Headphone   "},
	"d8:0f:99:88:c3:7a": {Name: "Fitbit Charge 3  "},
	"60:03:08:ac:bb:0d": {Name: "Mini iPad        "},
	"6e:f0:14:6d:38:a9": {Name: "??               "},
}

func getNameForBeacon(address string) string {
	name, exists := knownBeacons[address]
	if !exists {
		return address
	}
	return name.Name
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
			b.Services = make(map[string]ble.UUID)
			b.Name = getNameForBeacon(b.Address)
			beacons[a.Address().String()] = b
			return true
		}

		// Append any new services
		for _, srv := range a.Services() {
			_, e := b.Services[srv.String()]
			if !e {
				b.Services[srv.String()] = srv
			}
		}

		if a.Connectable() {
			rssi := a.RSSI()
			if rssi > 0 {
				rssi = -rssi
			}

			if intAbs(b.RSSI-rssi) > 30 {
				b.RSSI = (b.RSSI + rssi) / 2 // Take average
				return true
			}
		}

		return false
	}
	advHandler := func(a ble.Advertisement) {

		b, _ := beacons[a.Address().String()]

		if a.Connectable() {
			fmt.Printf("[%s] C %3d:", b.Name, intAbs(b.RSSI))
		} else {
			fmt.Printf("[%s] N %3d:", b.Name, intAbs(b.RSSI))
		}
		comma := ""
		if len(a.LocalName()) > 0 {
			fmt.Printf(" Name: %s", a.LocalName())
			comma = ","
		}
		if len(b.Services) > 0 {
			comma = ""
			for _, srv := range b.Services {
				if comma == "" {
					fmt.Printf("Services: ")
				}
				if ble.IsServiceKnown(srv) {
					fmt.Printf("%s %v", comma, ble.KnownServiceName(srv))
				}
				comma = ","
			}
			comma = ","
		}
		fmt.Printf("\n")
	}

	// Scan for specified durantion, or until interrupted by user.
	fmt.Printf("Scanning for %s...\n", *du)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.Scan(ctx, true, advHandler, advFilter))
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
