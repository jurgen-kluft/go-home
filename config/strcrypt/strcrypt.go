package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Encrypt or decrypt a string"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "encrypt, e",
			Usage: "The flag to indicate if we have to encrypt",
		},
		&cli.BoolFlag{
			Name:  "decrypt, d",
			Usage: "The flag to indicate if we have to decrypt",
		},
		&cli.StringFlag{
			Name:  "string, s",
			Usage: "The string to encrypt or decrypt",
		},
	}

	app.Action = func(c *cli.Context) error {
		str := ""
		var err error
		if c.Bool("encrypt") {
			str, err = config.Encrypt(c.String("string"))
		} else if c.Bool("decrypt") {
			str, err = config.Decrypt(c.String("string"))
		}
		if err == nil {
			fmt.Println(str)
		} else {
			fmt.Println(err.Error())
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
