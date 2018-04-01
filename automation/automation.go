package automation

// All automation logic is in this package
// Here we react to:
// - presence (people arriving/leaving)
// - switches (pressed)
// - events (timeofday, calendar)
// -

import (
	"github.com/jurgen-kluft/go-home/pubsub"
	"time"
)

func main() {

	for {
		connected := true
		for connected {
			client := pubsub.New()
			err := client.Connect("automation")
			if err == nil {

				client.Subscribe("state/+")

				for connected {
					select {
					case msg := <-client.InMsgs:
						topic := msg.Topic()
						if topic == "automation/config" {
						} else if topic == "client/disconnected" {
							connected = false
						}
						break
					case <-time.After(time.Second * 1):
						break
					}
				}
			} else {
				panic("Error on Client.Connect(): " + err.Error())
			}
		}

		// Wait for 10 seconds before retrying
		time.Sleep(10 * time.Second)
	}
}
