package main

import (
	"fmt"
)

type SensorModel uint8

const (
	BH1750 SensorModel = 0x0
	BME280 SensorModel = 0x1
	SCD4X  SensorModel = 0x2
)

type SensorType uint8

const (
	Temperature SensorType = 0x0 // (float, °C)
	Humidity    SensorType = 0x1 // (float, %)
	Pressure    SensorType = 0x2 // (float, hPa)
	Light       SensorType = 0x3 // (float, lux)
	CO2         SensorType = 0x4 // (float, ppm)
	VOC         SensorType = 0x5 // (float, ppm)
	PM1_0       SensorType = 0x6 // (float, µg/m3)
	PM2_5       SensorType = 0x7 // (float, µg/m3)
	PM10        SensorType = 0x8 // (float, µg/m3)
	Noise       SensorType = 0x9 // (float, dB)
	Presence    SensorType = 0xA // (float, 0.0-1.0)
)

type DeviceLocation uint16

const (
	Unknown    DeviceLocation = 0
	Location1  DeviceLocation = 0x01
	Location2  DeviceLocation = 0x02
	Location3  DeviceLocation = 0x03
	Location4  DeviceLocation = 0x04
	Location5  DeviceLocation = 0x05
	Location6  DeviceLocation = 0x06
	Location7  DeviceLocation = 0x07
	Location8  DeviceLocation = 0x08
	Area1      DeviceLocation = 0x10
	Area2      DeviceLocation = 0x20
	Area3      DeviceLocation = 0x30
	Area4      DeviceLocation = 0x40
	Area5      DeviceLocation = 0x50
	Area6      DeviceLocation = 0x60
	Area7      DeviceLocation = 0x70
	Area8      DeviceLocation = 0x80
	Bedroom    DeviceLocation = 0x100
	Livingroom DeviceLocation = 0x200
	Kitchen    DeviceLocation = 0x300
	Bathroom   DeviceLocation = 0x400
	Hallway    DeviceLocation = 0x500
	Balcony    DeviceLocation = 0x600
	Study      DeviceLocation = 0x700
	Pantry     DeviceLocation = 0x800
)

type SensorState uint8

const (
	Off   SensorState = 0x1
	On    SensorState = 0x2
	Error SensorState = 0x3
)

type FieldType uint8

const (
	TypeS8  FieldType = 0x0
	TypeS16 FieldType = 0x1
	TypeS32 FieldType = 0x2
	TypeU8  FieldType = 0x3
	TypeU16 FieldType = 0x4
	TypeU32 FieldType = 0x5
	TypeF32 FieldType = 0x6
)

func (t SensorType) String() string {
	switch t {
	case Temperature:
		return "Temperature"
	case Humidity:
		return "Humidity"
	case Pressure:
		return "Pressure"
	case Light:
		return "Light"
	case CO2:
		return "CO2"
	case Presence:
		return "Presence"
	default:
		return fmt.Sprintf("Unknown (%d)", t)
	}
}

func (d DeviceLocation) String() string {
	designator := d & 0xF0
	if designator == 0 {
		return fmt.Sprintf("Room %d", d&0x0F)
	}
	if designator == Bedroom {
		return fmt.Sprintf("Bedroom %d", d&0x0F)
	}
	if designator == Livingroom {
		return fmt.Sprintf("Livingroom %d", d&0x0F)
	}
	if designator == Kitchen {
		return fmt.Sprintf("Kitchen %d", d&0x0F)
	}
	if designator == Bathroom {
		return fmt.Sprintf("Bathroom %d", d&0x0F)
	}
	if designator == Hallway {
		return fmt.Sprintf("Hallway %d", d&0x0F)
	}
	if designator == Balcony {
		return fmt.Sprintf("Balcony %d", d&0x0F)
	}
	if designator == Study {
		return fmt.Sprintf("Studyroom %d", d&0x0F)
	}
	if designator == Pantry {
		return fmt.Sprintf("Pantry %d", d&0x0F)
	}
	return fmt.Sprintf("Unknown Location %d", d)
}
