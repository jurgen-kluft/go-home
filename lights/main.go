package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrobinson/deconz-go"
)

func main() {
	var (
		// ./lights -apiKey 0A498B9909 -host 10.0.0.18 -ct 500 -bri 150 -setState -id 18 (-isOn -1)
		host   = flag.String("host", "", "The IP or hostname of the gateway")
		port   = flag.Int("port", 80, "The port of the gateway")
		apiKey = flag.String("apiKey", "", "The API key of the gateway")

		resourceID = flag.Int("id", 0, "The ID of the light to control")
		isOn       = flag.Int("isOn", 0, "Whether to turn the light on or off (1 is on, -1 is off)")
		bri        = flag.Int("bri", 0, "The brightness to set")
		ct         = flag.Int("ct", 0, "The color temperature level to set")
		hue        = flag.Int("hue", 0, "The hue to set")
		sat        = flag.Int("sat", 0, "The saturation level to set")
		name       = flag.String("name", "", "The name of the device to set")
		lightIDs   = flag.String("lightIDs", "", "The comma separated list of light IDs to add to a group")

		create    = flag.Bool("create", false, "Whether to create the specified group")
		setState  = flag.Bool("setState", false, "Whether to set the specified state fields")
		setConfig = flag.Bool("setConfig", false, "Whether to set the specified config fields")
		delete    = flag.Bool("delete", false, "Whether to delete the specified group")
	)
	flag.Parse()

	c := deconz.NewClient(&http.Client{}, *host, *port, *apiKey)

	if *create {
		req := &deconz.CreateGroupRequest{
			Name: *name,
		}

		newID, err := c.CreateGroup(context.Background(), req)
		if err != nil {
			fmt.Printf("error creating group: %s\n", err.Error())
			return
		}
		fmt.Printf("created group %d\n", newID)
		*resourceID = newID
	}

	// No resource ID, get all groups
	if *resourceID < 1 {
		resp, err := c.GetGroups(context.Background())
		if err != nil {
			fmt.Printf("error getting groups: %s\n", err.Error())
			return
		}
		spew.Dump(resp)
		return
	}

	if *setState {
		req := &deconz.SetGroupStateRequest{SetLightStateRequest: deconz.SetLightStateRequest{}}

		if *isOn == -1 {
			off := new(bool)
			*off = false
			req = &deconz.SetGroupStateRequest{
				SetLightStateRequest: deconz.SetLightStateRequest{
					On: off,
				},
			}
		} else if *isOn == 1 {
			off := new(bool)
			*off = true
			req = &deconz.SetGroupStateRequest{
				SetLightStateRequest: deconz.SetLightStateRequest{
					On: off,
				},
			}
		}

		req.Brightness = *bri
		req.CT = *ct

		if *hue > 0 {
			req.Hue = hue
		}
		if *sat > 0 {
			req.Saturation = sat
		}

		err := c.SetGroupState(context.Background(), *resourceID, req)
		if err != nil {
			fmt.Printf("error setting group %d: %s\n", *resourceID, err.Error())
			return
		}

		fmt.Printf("set complete\n")
	}
	if *setConfig {
		req := &deconz.SetGroupConfigRequest{}

		if len(*name) > 0 {
			req.Name = *name
		}
		if len(*lightIDs) > 0 {
			lights := strings.Split(*lightIDs, ",")
			req.LightIDs = lights
		}

		err := c.SetGroupConfig(context.Background(), *resourceID, req)
		if err != nil {
			fmt.Printf("error setting group %d config: %s\n", *resourceID, err.Error())
			return
		}

		fmt.Printf("set complete\n")
	}
	if *delete {
		err := c.DeleteGroup(context.Background(), *resourceID)
		if err != nil {
			fmt.Printf("error deleting group %d: %s\n", *resourceID, err.Error())
			return
		}

		fmt.Printf("delete complete\n")
	}

	resp, err := c.GetGroup(context.Background(), *resourceID)
	if err != nil {
		fmt.Printf("error getting group %d: %s\n", *resourceID, err.Error())
		return
	}
	spew.Dump(resp)
}
