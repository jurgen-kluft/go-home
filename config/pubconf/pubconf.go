package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Publish a config to emitter broker"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: config.EmitterSecrets["host"],
			Usage: "The 'IP:Port' URI of the emitter broker",
		},
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

		host := c.String("host")
		filename := c.String("file")

		filedata, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		jsonstr := string(filedata)

		running := true
		for running {
			client := pubsub.New(host)
			channel := c.String("channel")
			register := []string{channel}
			subscribe := []string{}
			err := client.Connect("pubconf", register, subscribe)
			if err == nil {
				connected := true
				for connected {
					select {
					case msg := <-client.InMsgs:
						fmt.Printf("Emitter message received, topic:'%s', msg:'%s'\n", msg.Topic(), string(msg.Payload()))

					case <-time.After(time.Second * 1):
						err = client.PublishTTL(channel, jsonstr, 5*60)
						if err != nil {
							fmt.Println(err)
						}
						connected = false
						running = false
					}
				}
			}

			if err != nil {
				panic("Error: " + err.Error())
			}

			// Wait for 10 seconds before retrying
			time.Sleep(5 * time.Second)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
