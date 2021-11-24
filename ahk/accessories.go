package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
	"github.com/jurgen-kluft/go-home/config"
)

type television struct {
	*accessory.Accessory
	Tv *service.Television
}

func (t *television) PowerSelect(state int) {

}

type coloredLightbulb struct {
	*accessory.Accessory
	Light *service.ColoredLightbulb
	ID    int // deconz light ID
}

var (
	setLightStateURL = "http://%s/api/%s/lights/%d/state"
	setGroupStateURL = "http://%s/api/%s/groups/%d/action"
)

func turnOnOffLightGroup(groupID int, on bool) {
	url := fmt.Sprintf(setGroupStateURL, "10.0.0.18", "0A498B9909", groupID)
	stateJSON := "{ \"on\": false }"
	if on {
		stateJSON = "{ \"on\": true }"
	}
	body := strings.NewReader(stateJSON)
	request, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	response.Body.Close()
}

func (c *coloredLightbulb) Callback(onoff bool) {
	turnOnOffLightGroup(int(c.Accessory.ID), onoff)
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

func newTelevision(info accessory.Info) *television {
	acc := &television{}
	acc.Accessory = accessory.New(info, accessory.TypeTelevision)
	acc.Tv = service.NewTelevision()
	acc.AddService(acc.Tv.Service)
	return acc
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
	Televisions       []*television
}

func (a *accessories) initializeFromConfig(config *config.AhkConfig) []*accessory.Accessory {

	bridgeInfo := accessory.Info{Name: "Bridge", ID: 1}
	bridgeInfo.FirmwareRevision = "1.0"
	bridgeInfo.Manufacturer = "go-home"
	bridgeInfo.Model = "micro"
	bridgeInfo.Name = "home"
	bridgeInfo.SerialNumber = "090A1-93EAM0"
	a.Bridge = accessory.NewBridge(bridgeInfo)

	a.ColoredLights = make([]*coloredLightbulb, 0, 10)
	a.WhiteLights = make([]*lightbulb, 0, 10)
	a.MotionSensors = make([]*motionSensor, 0, 10)
	a.LightSensors = make([]*lightSensor, 0, 10)
	a.OccupancySensors = make([]*occupancySensor, 0, 10)
	a.ContactSensors = make([]*contactSensor, 0, 10)
	a.AirQualitySensors = make([]*airQualitySensor, 0, 10)
	a.Switches = make([]*button, 0, 10)
	a.Televisions = make([]*television, 0, 10)

	for _, lght := range config.Lights {
		if lght.Type == "colored" {
			lightbulb := newColoredLightbulb(accessory.Info{Name: lght.Name, ID: lght.ID, Manufacturer: lght.Manufacturer})
			lightbulb.Light.On.OnValueRemoteUpdate(lightbulb.Callback)
			a.ColoredLights = append(a.ColoredLights, lightbulb)
		} else if lght.Type == "white" {
			lightbulb := newLightbulb(accessory.Info{Name: lght.Name, ID: lght.ID})
			lightbulb.Light.On.OnValueRemoteUpdate(lightbulb.Callback)
			a.WhiteLights = append(a.WhiteLights, lightbulb)
		}
	}

	for _, ms := range config.Sensors {
		if ms.Type == "motion" {
			sensor := newMotionSensor(accessory.Info{Name: ms.Name, ID: ms.ID, Manufacturer: ms.Manufacturer})
			a.MotionSensors = append(a.MotionSensors, sensor)
		} else if ms.Type == "contact" {
			sensor := newContactSensor(accessory.Info{Name: ms.Name, ID: ms.ID})
			a.ContactSensors = append(a.ContactSensors, sensor)
		} else if ms.Type == "air-quality" {
			sensor := newAirQualitySensor(accessory.Info{Name: ms.Name, ID: ms.ID})
			a.AirQualitySensors = append(a.AirQualitySensors, sensor)
		} else if ms.Type == "occupancy" {
			sensor := newOccupancySensor(accessory.Info{Name: ms.Name, ID: ms.ID})
			a.OccupancySensors = append(a.OccupancySensors, sensor)
		}
	}

	for _, swtch := range config.Switches {
		sw := newButton(accessory.Info{Name: swtch.Name, ID: swtch.ID, Manufacturer: swtch.Manufacturer})
		sw.Button.On.OnValueRemoteUpdate(sw.Callback)
		a.Switches = append(a.Switches, sw)
	}

	for _, tv := range config.Televisions {
		t := newTelevision(accessory.Info{Name: tv.Name, ID: tv.ID, Manufacturer: tv.Manufacturer})
		t.Tv.PowerModeSelection.OnValueRemoteUpdate(t.PowerSelect)
		a.Televisions = append(a.Televisions, t)
	}

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
	for _, acc := range a.Televisions {
		accs = append(accs, acc.Accessory)
	}

	return accs
}

func accessoryTypeToString(atype accessory.AccessoryType) string {
	switch atype {
	case accessory.TypeBridge:
		return "Bridge"
	case accessory.TypeFan:
		return "Fan"
	case accessory.TypeGarageDoorOpener:
		return "Garage Door Opener"
	case accessory.TypeLightbulb:
		return "Lightbulb"
	case accessory.TypeDoorLock:
		return "Door Lock"
	case accessory.TypeOutlet:
		return "Outlet"
	case accessory.TypeSwitch:
		return "Switch"
	case accessory.TypeThermostat:
		return "Thermostat"
	case accessory.TypeSensor:
		return "Sensor"
	case accessory.TypeSecuritySystem:
		return "Security System"
	case accessory.TypeDoor:
		return "Door"
	case accessory.TypeWindow:
		return "Window"
	case accessory.TypeWindowCovering:
		return "Window Covering"
	case accessory.TypeProgrammableSwitch:
		return "Programmable Switch"
	case accessory.TypeIPCamera:
		return "IP Camera"
	case accessory.TypeVideoDoorbell:
		return "Video Doorbell"
	case accessory.TypeAirPurifier:
		return "Air Purifier"
	case accessory.TypeHeater:
		return "Heater"
	case accessory.TypeAirConditioner:
		return "Air Conditioner"
	case accessory.TypeHumidifier:
		return "Humidifier"
	case accessory.TypeDehumidifier:
		return "Dehumidifier"
	case accessory.TypeSprinklers:
		return "Sprinklers"
	case accessory.TypeFaucets:
		return "Faucets"
	case accessory.TypeShowerSystems:
		return "Shower System"
	case accessory.TypeTelevision:
		return "Television"
	case accessory.TypeRemoteControl:
		return "Remote Control"
	}
	return ""
}
