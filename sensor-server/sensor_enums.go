package main

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
	Humidity    SensorType = 0x01 // (s8, %)
	Pressure    SensorType = 0x02 // (s16, hPa)
	Light       SensorType = 0x03 // (s16, lux)
	CO2         SensorType = 0x04 // (s16, ppm)
	VOC         SensorType = 0x05 // (s16, ppm)
	PM1_0       SensorType = 0x06 // (s16, µg/m3)
	PM2_5       SensorType = 0x07 // (s16, µg/m3)
	PM10        SensorType = 0x08 // (s16, µg/m3)
	Noise       SensorType = 0x09 // (s16, dB)
	Presence    SensorType = 0x0A // (s8, 0-1)
	Distance    SensorType = 0x0B // (s16, cm)
	UV          SensorType = 0x0C // (s16, index)
	CO          SensorType = 0x0D // (s16, ppm)
	Vibration   SensorType = 0x0E // (s8,  <=16=none, <=64=low, <=128=medium, <=192=high, <=255=extreme)
	State       SensorType = 0xFF // (s32 (u8[4]), sensor model, sensor state)
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
	TypeBit  SensorFieldType = 0x01
	TypeS8   SensorFieldType = 0x08
	TypeS16  SensorFieldType = 0x10
	TypeS32  SensorFieldType = 0x20
)

func (t SensorFieldType) SizeInBits() int {
	return int(t)
}

func ToSensorFieldType(st SensorType) SensorFieldType {
	switch st {
	case Temperature:
		return TypeS8
	case Humidity:
		return TypeS8
	case Pressure:
		return TypeS16
	case Light:
		return TypeS16
	case CO2:
		return TypeS16
	case VOC:
		return TypeS16
	case PM1_0:
		return TypeS16
	case PM2_5:
		return TypeS16
	case PM10:
		return TypeS16
	case Noise:
		return TypeS8
	case Presence:
		return TypeS8
	case Distance:
		return TypeS16
	case UV:
		return TypeS8
	case CO:
		return TypeS8
	case Vibration:
		return TypeS8
	case State:
		return TypeS16
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
