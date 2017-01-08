package main

import (
	"bufio"
	"fmt"
	"github.com/jurgen-kluft/go-home/devices/hue/configuration"
	"github.com/jurgen-kluft/go-home/devices/hue/portal"
	"os"
)

func main() {
	//hubHostname := ssdpDiscover()
	pp, err := portal.GetPortal()
	if err != nil || len(pp) == 0 {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	hubHostname := pp[0].InternalIPAddress
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please press the link button on your hub, then press [enter] to continue.")
	reader.ReadLine()
	fmt.Println("Please enter your application name:")
	data, _, _ := reader.ReadLine()
	applicationName := string(data)
	fmt.Println("Please enter your device type:")
	data1, _, _ := reader.ReadLine()
	deviceType := string(data1)
	c := configuration.New(hubHostname)
	response, err := c.CreateUser(applicationName, deviceType)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	username := response[0].Success["username"].(string)
	fmt.Printf("Your username is %s\n", username)
}
