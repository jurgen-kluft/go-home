// just some code used when testing stuff out, might not work as expected and has alot of hard coded settings
package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/google/subcommands"
	"github.com/thomasf/yeelight/pkg/cli"
	"github.com/thomasf/yeelight/pkg/ssdp"
	"github.com/thomasf/yeelight/pkg/yeel"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cli.ListCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))

	// os.Exit(0)
	// main2()

	// d := yeel.Device{
	// 	Location: "yeelight://192.168.0.54:55443",
	// 	ID:       "poo",
	// 	Support:  map[string]bool{"set_name": true},
	// }
	// conn := &yeel.Conn{
	// 	Device: d,
	// }

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	err := conn.Open()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()
	// time.Sleep(time.Second)
	// log.Println("execmc")
	// res, err := conn.ExecCommand(yeel.SetNameCommand{Name: "tf-desk"})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(res)
	// log.Println("dn")
	// wg.Wait()

}
func main2() {
	rand.Seed(time.Now().UTC().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	deviceCh := make(chan yeel.Device, 100)

	var wg sync.WaitGroup

	conn := ssdp.Conn{}
	err := conn.Start("ens135")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := conn.Listen(deviceCh)
		if err != nil {
			log.Println(err)
		}

	}()
	time.Sleep(time.Second)
	go func() {
		for {
			err := conn.Search()
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Minute)
		}
	}()

	// trackedDevice .
	type trackedDevice struct {
		Device             yeel.Device
		LastUpdated        time.Time
		ConnectionAttempts int
	}

	conns := make(map[string]*yeel.Conn, 0)
	connEndedC := make(chan string, 0)
	tickr := time.NewTicker(time.Minute)
	devices := make(map[string]trackedDevice, 0)

	handleConn := func(conn *yeel.Conn) {
		err := conn.Open()
		if err != nil {
			log.Println(err)
		}
		connEndedC <- conn.Device.ID
	}

	for {
		select {
		case <-tickr.C:
			for _, v := range conns {
				go func(c *yeel.Conn) {
					res, err := c.ExecCommand(
						yeel.GetPropCommand{
							Properties: []string{"power"},
						},
					)
					if err != nil {
						log.Println(err)
					}
					log.Println(c.Device.ID, res)
				}(v)
			}

		case id := <-connEndedC:
			delete(conns, id)
			if d, ok := devices[id]; ok {
				if d.ConnectionAttempts < 5 {
					conn := &yeel.Conn{Device: d.Device}
					conns[d.Device.ID] = conn
					d.ConnectionAttempts = d.ConnectionAttempts + 1
					go handleConn(conn)

				}
			}

		case d := <-deviceCh:
			devices[d.ID] = trackedDevice{
				Device:      d,
				LastUpdated: time.Now(),
			}
			if c, ok := conns[d.ID]; !ok {
				conn := &yeel.Conn{Device: d}
				conns[d.ID] = conn
				go handleConn(conn)
			} else {
				log.Println("already connnnn", c)
			}
		}
	}
	wg.Add(1)
	wg.Wait()

	// // d3 := time.Second
	// for d := range deviceCh {

	// 	err := conn.Open()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// {
	// 	commands := []yeel.Commander{
	// 		// yeel.SetNameCommand{
	// 		// 	Name: "23ddd",
	// 		// },
	// 		yeel.SetBrightnessNormCommand{
	// 			Brightness: 0.1,
	// 			Duration:   500 * time.Millisecond,
	// 		},
	// 		yeel.SetBrightnessNormCommand{
	// 			Brightness: 1,
	// 			Duration:   500 * time.Millisecond,
	// 		},
	// 		yeel.SetRGBNormCommand{
	// 			R: .5, G: .5, B: .5,
	// 			Duration: time.Second,
	// 		},
	// 		yeel.SetRGBNormCommand{
	// 			R: .8, G: .5, B: .3,
	// 			Duration: time.Second,
	// 		},
	// 		yeel.SetRGBNormCommand{
	// 			R: .1, G: .5, B: .9,
	// 			Duration: time.Second,
	// 		},
	// 		yeel.SetBrightnessNormCommand{
	// 			Brightness: 1,
	// 			Duration:   500 * time.Millisecond,
	// 		},
	// 		yeel.SetBrightnessNormCommand{
	// 			Brightness: 0.5,
	// 			Duration:   500 * time.Millisecond,
	// 		},

	// 		yeel.SetHSVNormCommand{
	// 			Hue:        0,
	// 			Saturation: 1,
	// 			Duration:   time.Second,
	// 		},
	// 		yeel.SetHSVNormCommand{
	// 			Hue:        0.5,
	// 			Saturation: 0.5,
	// 			Duration:   time.Second,
	// 		},
	// 	}

	// 	for _, cmd := range commands {
	// 		res, err := conn.ExecCommand(cmd)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		log.Println(res)

	// 		time.Sleep(1100 * time.Millisecond)

	// 	}
	// }

	// {
	// 	// d1 := 150 * time.Millisecond
	// 	// d1 := 1500 * time.Millisecond
	// 	d1 := 1000 * time.Millisecond
	// 	// d2 := d1
	// 	// d2 := 400 * time.Millisecond
	// 	d2 := 1000 * time.Millisecond

	// 	// dur := time.Duration(rand.Intn(500)+3000) * time.Millisecond

	// 	// H := rand.Intn(259)
	// 	// S := rand.Intn(100)
	// 	// res, err := conn.SetHSV(H,S, dur)

	// 	expressions := yeel.Animation{

	// 		// yeel.SleepKeyframe{Duration: d2},
	// 		// yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(1, 0, 0)},
	// 		// yeel.SleepKeyframe{Duration: d2},
	// 		// yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(0, 0, 1)},
	// 		// yeel.SleepKeyframe{Duration: d2},
	// 		// yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(1, 0, 0)},
	// 		// yeel.SleepKeyframe{Duration: d2},
	// 		// yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(0, 0, 1)},

	// 		// // yeel.SleepKeyframe{Duration: d2},
	// 		// yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(0, 0, 0)},
	// 		// // yeel.SleepKeyframe{Duration: d2},
	// 		// yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(0.2, 0, 0.4)},
	// 		// // yeel.SleepKeyframe{Duration: d2},
	// 		// yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(0.8, 0.4, 0.7)},
	// 		// // yeel.SleepKeyframe{Duration: d2},
	// 		// yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(1, 0.5, 1)},
	// 		// // yeel.SleepKeyframe{Duration: d2},
	// 		// yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(1, 1, 1)},
	// 		yeel.SleepKeyframe{Duration: d2},

	// 		// yeel.SleepKeyframe{time.Second},
	// 		// yeel.ColorTempratureKeyframe{800 * time.Millisecond, 50, 1700},
	// 		yeel.ColorTempratureKeyframe{2000 * time.Millisecond, 70, 4000},
	// 		yeel.ColorTempratureKeyframe{2000 * time.Millisecond, 100, 6400},
	// 	}
	// 	expressions = yeel.Animation{
	// 		yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(1, 0, 0)},
	// 		yeel.SleepKeyframe{Duration: d1},
	// 		yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(0, 1, 0)},
	// 		yeel.SleepKeyframe{Duration: d1},
	// 		yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(1, 0, 1)},
	// 		yeel.SleepKeyframe{Duration: d1},
	// 		yeel.RGBKeyframe{Duration: d1, Brightness: 100, RGB: yeel.NewRGBNorm(0, 1, 0)},
	// 		yeel.SleepKeyframe{Duration: d1},
	// 		yeel.RGBKeyframe{Duration: d2, Brightness: 100, RGB: yeel.NewRGBNorm(1, 1, 0)},
	// 		yeel.RGBKeyframe{Duration: d2, Brightness: 100, RGB: yeel.NewRGBNorm(0, 1, 0)},
	// 		yeel.RGBKeyframe{Duration: d2, Brightness: 100, RGB: yeel.NewRGBNorm(1, 0, 1)},
	// 		yeel.RGBKeyframe{Duration: d2, Brightness: 100, RGB: yeel.NewRGBNorm(0, 1, 1)},
	// 	}
	// 	cmd := yeel.StartColorFlowCommand{
	// 		Count:     10000,
	// 		Action:    0,
	// 		Animation: expressions,
	// 	}
	// 	res, err := conn.ExecCommand(cmd)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println("RESSSIL", res)
	// 	time.Sleep(time.Minute)
	// 	os.Exit(0)
	// }

	// }

	// wg.Wait()
	// setname()

}
