package main

import (
	"context"
	"fmt"
	"github.com/brutella/dnssd"
	slog "log"
	"os"
	"os/signal"
	"time"
)

var timeFormat = "15:04:05.000"

func launchConfigService() {
	cfg := dnssd.Config{
		Name:   "Go Home",
		Type:   "_http._tcp",
		Domain: "local",
		Port:   12345,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if resp, err := dnssd.NewResponder(); err != nil {
		fmt.Println(err)
	} else {
		srv, err := dnssd.NewService(cfg)
		if err != nil {
			slog.Fatal(err)
		}

		go func() {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)

			select {
			case <-stop:
				cancel()
			}
		}()

		go func() {
			time.Sleep(100 * time.Millisecond)
			handle, err := resp.Add(srv)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%s	Got a reply for service %s: Name now registered and active\n", time.Now().Format(timeFormat), handle.Service().ServiceInstanceName())
			}
		}()
		err = resp.Respond(ctx)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {

	launchConfigService()
}
