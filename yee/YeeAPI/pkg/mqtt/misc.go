package mqtt

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/thomasf/lg"
)

var validModels = map[string]bool{
	"color":  true,
	"mono":   true,
	"stripe": true,
}

func parseMessage(msg MQTT.Message) (Command, error) {

	topics := strings.SplitN(msg.Topic(), "/", 4)

	var m Command

	if len(topics) < 3 {
		return m, fmt.Errorf("too short topic path: %s", msg.Topic)
	}

	if topics[1] != "command" {
		return m, fmt.Errorf("not a command: %s", msg.Topic)
	}

	model := strings.TrimPrefix(topics[0], "yl")
	if !validModels[model] {
		return m, fmt.Errorf("invalid model: %s", topics[0])
	}
	m.DeviceModel = model

	m.DeviceName = topics[2]

	if len(topics) > 3 {
		m.Command = topics[3]
	}
	m.Value = string(msg.Payload())
	if lg.V(10) {
		lg.Infof("parsed message:\n%s\n%s", spew.Sdump(msg), spew.Sdump(m))
	}

	return m, nil
}

func publish(t *Device, updates []PropUpdate, s MQTT.Client) error {
	d := t.Device
	id := d.ID
	retain := false
	if d.Name != "" {
		id = d.Name
		retain = true
	}
	for _, v := range updates {
		topic := fmt.Sprintf("yl%s/state/%s/%s", t.Device.Model, id, v.Prop)
		s.Publish(topic, 0, retain, []byte(v.Value))

	}
	return nil
}

func publishUpdates(t *Device, mqttService MQTT.Client) bool {
	updates := t.Updates()
	n := len(updates)
	if n > 0 {
		err := publish(t, updates, mqttService)
		if err != nil {
			lg.Errorln(err)
		}
		return true
	}
	return false
}
