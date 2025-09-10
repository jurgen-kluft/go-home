package main

import (
	"fmt"
)

type SensorLocation uint8

const (
	Unknown    SensorLocation = 0
	Location1  SensorLocation = 0x01
	Location2  SensorLocation = 0x02
	Location3  SensorLocation = 0x03
	Location4  SensorLocation = 0x04
	Location5  SensorLocation = 0x05
	Location6  SensorLocation = 0x06
	Location7  SensorLocation = 0x07
	Location8  SensorLocation = 0x08
	Location15 SensorLocation = 0x0F
	Bedroom    SensorLocation = 0x10
	Livingroom SensorLocation = 0x20
	Kitchen    SensorLocation = 0x30
	Bathroom   SensorLocation = 0x40
	Hallway    SensorLocation = 0x50
	Balcony    SensorLocation = 0x60
	Study      SensorLocation = 0x70
	Pantry     SensorLocation = 0x80
)

type SensorModel uint8

const (
	GPIO   SensorModel = 0x00
	BH1750 SensorModel = 0x10
	BME280 SensorModel = 0x20
	SCD4X  SensorModel = 0x30
)

type SensorType uint8

const (
	Temperature SensorType = 0x00 // (s8, °C)
	Humidity    SensorType = 0x01 // (u8, %)
	Pressure    SensorType = 0x02 // (u16, hPa)
	Light       SensorType = 0x03 // (u16, lux)
	CO2         SensorType = 0x04 // (u16, ppm)
	VOC         SensorType = 0x05 // (u16, ppm)
	PM1_0       SensorType = 0x06 // (u16, µg/m3)
	PM2_5       SensorType = 0x07 // (u16, µg/m3)
	PM10        SensorType = 0x08 // (u16, µg/m3)
	Noise       SensorType = 0x09 // (u16, dB)
	Presence    SensorType = 0x0A // (u8, 0-1)
	Distance    SensorType = 0x0B // (u16, cm)
	UV          SensorType = 0x0C // (u16, index)
	CO          SensorType = 0x0D // (u16, ppm)
	Vibration   SensorType = 0x0E // (u8, 0=none, 1=low, 2=medium, 3=high)
	State       SensorType = 0x0F // (u16 (u8[2]), sensor model, sensor state)
)

type SensorState uint8

const (
	Off   SensorState = 0x10
	On    SensorState = 0x20
	Error SensorState = 0x30
)

type SensorFieldType uint8

const (
	TypeNone SensorFieldType = 0x00
	TypeS8   SensorFieldType = 0x01
	TypeS16  SensorFieldType = 0x02
	TypeS32  SensorFieldType = 0x03
	TypeU8   SensorFieldType = 0x04
	TypeU16  SensorFieldType = 0x05
	TypeU32  SensorFieldType = 0x06
)

func ToSensorFieldType(st SensorType) SensorFieldType {
	switch st {
	case Temperature:
		return TypeS8
	case Humidity:
		return TypeU8
	case Pressure:
		return TypeU16
	case Light:
		return TypeU16
	case CO2:
		return TypeU16
	case VOC:
		return TypeU16
	case PM1_0:
		return TypeU16
	case PM2_5:
		return TypeU16
	case PM10:
		return TypeU16
	case Noise:
		return TypeU8
	case Presence:
		return TypeU8
	case Distance:
		return TypeU16
	case UV:
		return TypeU8
	case CO:
		return TypeU8
	case Vibration:
		return TypeU8
	case State:
		return TypeU16
	}
	return TypeNone
}

// String returns the string representation of the SensorType.
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
	case VOC:
		return "VOC"
	case PM1_0:
		return "PM1.0"
	case PM2_5:
		return "PM2.5"
	case PM10:
		return "PM10"
	case Noise:
		return "Noise"
	case Presence:
		return "Presence"
	case Distance:
		return "Distance"
	case UV:
		return "UV"
	case CO:
		return "CO"
	case Vibration:
		return "Vibration"
	case State:
		return "State"
	default:
		return "Unknown SensorType"
	}
}

// String returns the string representation of the SensorLocation.
func (d SensorLocation) String() string {
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
