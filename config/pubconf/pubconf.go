package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jurgen-kluft/go-home/micro-service"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Publish a config to emitter broker"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Value: "flux.config.json",
			Usage: "The JSON configuration file to read and publish",
		},
		cli.StringFlag{
			Name:  "channel",
			Value: "config/flux/",
			Usage: "The channel to publish to",
		},
	}

	app.Action = func(c *cli.Context) error {

		filename := c.String("file")

		filedata, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		jsonstr := string(filedata)

		channel := c.String("channel")
		register := []string{channel}
		subscribe := []string{}

		m := microservice.New("pubconf")
		m.RegisterAndSubscribe(register, subscribe)

		m.RegisterHandler("*", func(m *microservice.Service, topic string, msg []byte) bool {
			fmt.Printf("message received, topic:'%s', msg:'%s'\n", topic, string(msg))
			return true
		})

		m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
			err := m.Pubsub.PublishTTL(channel, jsonstr, 5*60)
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
