package main

import (
	"fmt"
	"log"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type coloredLightbulb struct {
	*accessory.Accessory
	Light *service.ColoredLightbulb
}

func (c *coloredLightbulb) Callback(onoff bool) {

}

type lightbulb struct {
	*accessory.Accessory
	Light *service.Lightbulb
}

func (c *lightbulb) Callback(onoff bool) {

}

type button struct {
	*accessory.Accessory
	Button *service.Switch
}

func (c *button) Callback(onoff bool) {

}

type motionSensor struct {
	*accessory.Accessory
	Sensor *service.MotionSensor
}

type lightSensor struct {
	*accessory.Accessory
	Sensor *service.LightSensor
}

type occupancySensor struct {
	*accessory.Accessory
	Sensor *service.OccupancySensor
}

type contactSensor struct {
	*accessory.Accessory
	Sensor *service.ContactSensor
}

type airQualitySensor struct {
	*accessory.Accessory
	Sensor *service.AirQualitySensor
}

func newColoredLightbulb(info accessory.Info) *coloredLightbulb {
	acc := &coloredLightbulb{}
	acc.Accessory = accessory.New(info, accessory.TypeLightbulb)
	acc.Light = service.NewColoredLightbulb()
	acc.Light.Brightness.SetValue(100)
	acc.AddService(acc.Light.Service)
	return acc
}
func newLightbulb(info accessory.Info) *lightbulb {
	acc := &lightbulb{}
	acc.Accessory = accessory.New(info, accessory.TypeLightbulb)
	acc.Light = service.NewLightbulb()
	acc.AddService(acc.Light.Service)
	return acc
}
func newButton(info accessory.Info) *button {
	acc := &button{}
	acc.Accessory = accessory.New(info, accessory.TypeSwitch)
	acc.Button = service.NewSwitch()
	acc.AddService(acc.Button.Service)
	return acc
}

func newMotionSensor(info accessory.Info) *motionSensor {
	acc := &motionSensor{}
	acc.Accessory = accessory.New(info, accessory.TypeSensor)
	acc.Sensor = service.NewMotionSensor()
	acc.AddService(acc.Sensor.Service)
	return acc
}

func newLightSensor(info accessory.Info) *lightSensor {
	acc := &lightSensor{}
	acc.Accessory = accessory.New(info, accessory.TypeSensor)
	acc.Sensor = service.NewLightSensor()
	acc.AddService(acc.Sensor.Service)
	return acc
}

func newOccupancySensor(info accessory.Info) *occupancySensor {
	acc := &occupancySensor{}
	acc.Accessory = accessory.New(info, accessory.TypeSensor)
	acc.Sensor = service.NewOccupancySensor()
	acc.AddService(acc.Sensor.Service)
	return acc
}

func newContactSensor(info accessory.Info) *contactSensor {
	acc := &contactSensor{}
	acc.Accessory = accessory.New(info, accessory.TypeSensor)
	acc.Sensor = service.NewContactSensor()
	acc.AddService(acc.Sensor.Service)
	return acc
}

func newAirQualitySensor(info accessory.Info) *airQualitySensor {
	acc := &airQualitySensor{}
	acc.Accessory = accessory.New(info, accessory.TypeSensor)
	acc.Sensor = service.NewAirQualitySensor()
	acc.AddService(acc.Sensor.Service)
	return acc
}

type accessories struct {
	Bridge            *accessory.Bridge
	ColoredLights     []*coloredLightbulb
	WhiteLights       []*lightbulb
	MotionSensors     []*motionSensor
	LightSensors      []*lightSensor
	OccupancySensors  []*occupancySensor
	ContactSensors    []*contactSensor
	AirQualitySensors []*airQualitySensor
	Switches          []*button
}

func (a *accessories) registerAll() []*accessory.Accessory {
	a.Bridge = accessory.NewBridge(accessory.Info{Name: "Bridge", ID: 1})

	lightbulb := newColoredLightbulb(accessory.Info{Name: "Sophia Stand", ID: 3})
	lightbulb.Light.On.OnValueRemoteUpdate(lightbulb.Callback)

	a.ColoredLights = make([]*coloredLightbulb, 0, 10)
	a.ColoredLights = append(a.ColoredLights, lightbulb)

	button1 := newButton(accessory.Info{Name: "Meeting", ID: 2})
	button1.Button.On.OnValueRemoteUpdate(button1.Callback)
	a.Switches = make([]*button, 0, 10)
	a.Switches = append(a.Switches, button1)

	// Colored bulbs (HUE)
	// Light bulbs (IKEA)
	// Switches (Aqara)

	// Air quality sensor (PM 2.5, PM 10, Ozone, Nitrogen Dioxide, CO, CO2)
	// Contact Sensor (Front door)
	// Light Sensor
	// Motion Sensor
	// Occupancy Sensor

	accs := make([]*accessory.Accessory, 0, 10)
	for _, acc := range a.ColoredLights {
		accs = append(accs, acc.Accessory)
	}
	for _, acc := range a.WhiteLights {
		accs = append(accs, acc.Accessory)
	}
	for _, acc := range a.MotionSensors {
		accs = append(accs, acc.Accessory)
	}
	for _, acc := range a.LightSensors {
		accs = append(accs, acc.Accessory)
	}
	for _, acc := range a.OccupancySensors {
		accs = append(accs, acc.Accessory)
	}
	for _, acc := range a.ContactSensors {
		accs = append(accs, acc.Accessory)
	}
	for _, acc := range a.AirQualitySensors {
		accs = append(accs, acc.Accessory)
	}
	for _, acc := range a.Switches {
		accs = append(accs, acc.Accessory)
	}

	return accs
}

func main() {
	pin := "53664337"
	acsrs := &accessories{}
	accs := acsrs.registerAll()

	// configure the ip transport
	config := hc.Config{Pin: pin}
	fmt.Println("bridge: " + acsrs.Bridge.Info.Name.GetValue())
	for _, acc := range accs {
		fmt.Println("   accessory: " + acc.Info.Name.GetValue())
	}
	t, err := hc.NewIPTransport(config, acsrs.Bridge.Accessory, accs...)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
