package mqtt

import (
	"time"

	"github.com/thomasf/yeelight/pkg/yeel"
)

var (
	animations map[string]yeel.Animation
)

func init() {
	animations = make(map[string]yeel.Animation)

	// strobe
	{

		animD := 50 * time.Millisecond
		// sleepD := 50*  time.Millisecond

		bri_max := yeel.NewBrightnessNorm(1).Int()
		bri_min := yeel.NewBrightnessNorm(0).Int()
		black := yeel.NewRGBNorm(0, 0, 0)
		white := yeel.NewRGBNorm(1, 1, 1)
		// color := yeel.NewRGBNorm(0, 1, 1)
		//color := yeel.NewRGBNorm(0, 0, 0)

		anim := yeel.Animation{
			yeel.RGBKeyframe{Duration: animD, Brightness: bri_min, RGB: black},
			// yeel.SleepKeyframe{Duration: sleepD},
			yeel.RGBKeyframe{Duration: animD, Brightness: bri_max, RGB: white},
			// yeel.SleepKeyframe{Duration: sleepD},
		}

		animations["strobe"] = anim

	}
}
