package lightshue

import (
	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/pubsub"
	"github.com/stefanwichmann/go.hue"
)

type instance struct {
	config *config.HueConfig
}

func main() {
	hue := &instance{}
	for {
		pb := pubsub.New()
		err := pb.Connect("hue")
		if err == nil {

			client.Subscribe("sensor/light/+")

			for {
				select {
				case msg := <-pb.InMsgs:
					if msg.Topic() == "hue/config" {
						if hue.config == nil {
							hue.config, err = config.HueConfigFromJSON(string(msg.Payload()))
						}
					} else if msg.Topic() == "sensor/light/hue" {
						sensor, _ := config.SensorStateFromJSON(string(msg.Payload()))
					} else if msg.Topic() == "sensor/light/yee" {
						yeesensor, _ := NewSuncalc(msg.Data)
					}
					break
				case <-time.After(time.Second * 10):
					// do something if messages are taking too long
					// or if we haven't received enough state info.
					Process(flux, client)
					break
				}
			}
		}

		// Wait for 10 seconds before retrying
		time.Sleep(10 * time.Second)
	}
}
