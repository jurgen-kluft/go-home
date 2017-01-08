package main

/*
Groups of Lights:

* Bedroom
* Livingroom
* Kitchen
* SophiaRoom
* JenniferRoom

Scenes:
We pre-program scenes and we reference them when we modify the group.

Season   : TimeOfDay   : Group     = Scene
"WINTER" : "BREAKFAST" : "Kitchen" = "MorningWinter"
"SPRING" : "BREAKFAST" : "Kitchen" = "MorningSpring"
"WINTER" : "LUNCH"     : "Kitchen" = "NoonWinter"
"WINTER" : "DINNER"    : "Kitchen" = "SunSetWinter"
"WINTER" : "EVENING"   : "Kitchen" = "EveningWinter"
"WINTER" : "EVENING"   : "Bedroom" = "LateEveningWinter"

All these configurations are stored in REDIS, the key being
the TimeOfDay. When receiving the value we need to get the
matching Season/Group to find out the scene to apply.

We can also scope multiple fields into the key:
     KEY                     |      VALUE
Kitchen:Breakfast:Winter     |  MorningWinter
Kitchen:Breakfast:Summer     |  MorningSummer
Bedroom:Evening:Winter       |  LateEveningWinter

JSON
// KEY = TimeOfDay-BREAKFAST
{
    "scenes" : [
        { "season" : "Winter", "group" : "Kitchen", "scene" : "MorningWinter" },
        { "season" : "Spring", "group" : "Kitchen", "scene" : "MorningSpring" },
    ]
}

Then we need to create the HUE light configurations for the above scenes:
- MorningWinter
- MorningSpring
- NoonWinter
- SunSetWinter
- EveningWinter
- LateEveningWinter
- MorningSummer
- NoonSummer
-

Just need to find a way to create and easily modify these configurations.

*/

import (
	"flag"
	"fmt"
	"github.com/jurgen-kluft/go-home/devices/hue/groups"
	"github.com/jurgen-kluft/go-home/devices/hue/lights"
	"github.com/jurgen-kluft/go-home/devices/hue/portal"
	"os"
	"time"
)

var (
	apiKey     string = ""
	blinkState lights.State
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: get-light-state -key=[string]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	blinkState = lights.State{On: true, Alert: "lselect"}
	flag.StringVar(&apiKey, "key", os.Getenv("HUE_USERNAME"), "hue light api key")
	flag.Parse()
	flag.Usage = usage
}

func main() {
	if apiKey != "" {
		pp, err := portal.GetPortal()
		if err != nil {
			fmt.Println("portal.GetPortal() ERROR: ", err)
			os.Exit(1)
		}
		ll := lights.New(pp[0].InternalIPAddress, apiKey)
		allLights, err := ll.GetAllLights()
		if err != nil {
			fmt.Println("lights.GetAllLights() ERROR: ", err)
			os.Exit(1)
		}
		fmt.Println()
		fmt.Println("Lights")
		fmt.Println("------")
		for _, l := range allLights {
			fmt.Printf("ID: %d Name: %s\n", l.ID, l.Name)
		}
		gg := groups.New(pp[0].InternalIPAddress, apiKey)
		allGroups, err := gg.GetAllGroups()
		if err != nil {
			fmt.Println("groups.GetAllGroups() ERROR: ", err)
			os.Exit(1)
		}
		fmt.Println()
		fmt.Println("Groups")
		fmt.Println("------")
		for _, g := range allGroups {
			fmt.Printf("ID: %d Name: %s\n", g.ID, g.Name)
			for _, lll := range g.Lights {
				fmt.Println("\t", lll)
			}
			previousState := g.Action
			_, err := gg.SetGroupState(g.ID, blinkState)
			if err != nil {
				fmt.Println("groups.SetGroupState() ERROR: ", err)
				os.Exit(1)
			}
			time.Sleep(time.Second * time.Duration(10))
			_, err = gg.SetGroupState(g.ID, previousState)
			if err != nil {
				fmt.Println("groups.SetGroupState() ERROR: ", err)
				os.Exit(1)
			}
		}
	} else {
		usage()
	}
}
