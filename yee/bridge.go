package main

import (
	"github.com/brutella/hc/accessory"
	"github.com/oleggator/goyeelight"
	"log"
	"strconv"
)

// Accessory factory
func NewYeelight(info accessory.Info, ip, port string) *accessory.ColoredLightbulb {
	yeelight := goyeelight.New(ip, port)

	acc := accessory.NewColoredLightbulb(info)

	r, err := yeelight.GetProp("sat", "hue", "power", "bright")
	if err != nil {
		log.Println(err)
	}
	var (
		saturation = r["sat"]
		hue        = r["hue"]
		power      = r["power"]
		bright     = r["bright"]
	)

	err = updateStatus(acc, saturation, hue, power, bright)
	if err != nil {
		log.Println(err)
	}

	acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		var result string
		var err error
		if on == true {
			result, err = yeelight.SetPower("on", "smooth", "100")
		} else {
			result, err = yeelight.SetPower("off", "smooth", "100")
		}

		if err != nil {
			log.Println(err)
		}

		log.Println("Power:", result)
	})

	acc.Lightbulb.Brightness.OnValueRemoteUpdate(func(brightness int) {
		result, err := yeelight.SetBright(strconv.Itoa(brightness), "smooth", "500")
		if err != nil {
			log.Println(err)
		}

		log.Println("Brightness:", result)
	})

	acc.Lightbulb.Hue.OnValueRemoteUpdate(func(hueInput float64) {
		hue = strconv.Itoa(int(hueInput))

		result, err := yeelight.SetHSV(hue, saturation, "smooth", "500")
		if err != nil {
			log.Println(err)
		}

		log.Println("Color:", result)
	})

	acc.Lightbulb.Saturation.OnValueRemoteUpdate(func(saturationInput float64) {
		saturation = strconv.Itoa(int(saturationInput))

		result, err := yeelight.SetHSV(hue, saturation, "smooth", "500")
		if err != nil {
			log.Println(err)
		}
		log.Println("Saturation:", result)
	})

	return acc
}

// Update HomeKit accessory status
func updateStatus(lightbulb *accessory.ColoredLightbulb, saturation, hue, power, bright string) error {
	lightbulb.Lightbulb.On.SetValue(power == "on")

	brightInt, err := strconv.Atoi(bright)
	if err != nil {
		return err
	}
	lightbulb.Lightbulb.Brightness.SetValue(brightInt)

	hueInt, err := strconv.Atoi(hue)
	if err != nil {
		return err
	}
	lightbulb.Lightbulb.Hue.SetValue(float64(hueInt))

	saturationInt, err := strconv.Atoi(saturation)
	if err != nil {
		return err
	}
	lightbulb.Lightbulb.Saturation.SetValue(float64(saturationInt))

	return nil
}
