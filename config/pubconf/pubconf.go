package main

import (
	"fmt"
	"log"
	"os"
	"time"

	microservice "github.com/jurgen-kluft/go-home/micro-service"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Publish a config to emitter broker"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "file",
			Value: "flux.config.json",
			Usage: "The JSON configuration file to read and publish",
		},
		&cli.StringFlag{
			Name:  "destination",
			Value: "config/flux",
			Usage: "The channel to write to",
		},
	}

	app.Action = func(c *cli.Context) error {

		filename := c.String("file")
		channel := c.String("destination")

		filedata, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		jsonbytes := filedata

		m, err := microservice.New("conf", "config/config", time.Second*15)
		if err != nil {
			return err
		}

		m.RegisterHandler("*", func(m *microservice.Service, msg *microservice.Message) bool {
			fmt.Printf("message received, from:'%d', msg:'%s'\n", msg.SrcID(), string(msg.Payload))
			return true
		})

		m.RegisterHandler("tick", func(m *microservice.Service, msg *microservice.Message) bool {
			err := m.SendJsonTo(channel, string(jsonbytes))
			if err != nil {
				fmt.Println(err)
			}

			// Only do one tick
			return false
		})

		m.Loop()

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
