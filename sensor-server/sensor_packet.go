package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Packet structure
// {
//     u16                   length; // Number of bytes in the packet
//     u16                   sequence; // Sequence number of the packet
//     u8                    version; // Version of the packet structure
//     DeviceLocation::Value location;
//     DeviceLabel::Value    label;
//     u8                    count;  // Number of sensor values in the packet (max 16)

//     // sensor value 1
//     u8 type_and_channel;
//     u8 state_and_field_type;
//     union
//     {
//         s8  s8_value;
//         s16 s16_value;
//         s32 s32_value;
//         u8  u8_value;
//         u16 u16_value;
//         u32 u32_value;
//         f32 f32_value;
//     } value;

//     // sensor value 2

//     // ... (up to max 16 sensor values)

//     // terminator, 2 bytes
//     u16 terminator; // 0xFFFF
// };

type SensorPacket struct {
	length   uint16
	sequence uint16
	version  uint8
	location DeviceLocation
	label    uint8
	count    uint8
	values   []SensorValue
}

type SensorValue struct {
	sensorType SensorType
	channel    SensorChannel
	state      SensorState
	fieldType  FieldType

	valueInt8  int8
	valueInt16 int16
	valueInt32 int32

	valueUint8  uint8
	valueUint16 uint16
	valueUint32 uint32

	valueFloat32 float32
}

func DecodeNetworkPacket(data []byte) (SensorPacket, error) {
	if len(data) < 4 {
		return SensorPacket{}, fmt.Errorf("data too short")
	}

	pkt := SensorPacket{
		length:   uint16(data[0]) | uint16(data[1])<<8,
		sequence: uint16(data[2]) | uint16(data[3])<<8,
		version:  data[4],
		location: DeviceLocation(data[5]),
		label:    data[6],
		count:    data[7],
		values:   make([]SensorValue, 0, data[7]),
	}

	fmt.Printf("Number of values: %d\n", pkt.count)

	if len(data) < int(pkt.length) {
		return pkt, fmt.Errorf("data length mismatch, %d < %d", len(data), pkt.length)
	}

	offset := 8
	for i := uint8(0); i < pkt.count; i++ {
		value := SensorValue{
			sensorType: SensorType(data[offset] >> 4),
			channel:    SensorChannel(data[offset] & 0x0F),
			state:      SensorState(data[offset+1] >> 4),
			fieldType:  FieldType(data[offset+1] & 0x0F),
		}

		offset += 2
		// depending on fieldType, read the appropriate value.
		// the written values are in little-endian format
		switch value.fieldType {
		case TypeS8:
			value.valueInt8 = int8(data[offset])
			offset += 1
		case TypeS16:
			value.valueInt16 = int16(binary.LittleEndian.Uint16(data[offset : offset+2]))
			offset += 2
		case TypeS32:
			value.valueInt32 = int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
			offset += 4
		case TypeU8:
			value.valueUint8 = data[offset]
			offset += 1
		case TypeU16:
			value.valueUint16 = binary.LittleEndian.Uint16(data[offset : offset+2])
			offset += 2
		case TypeU32:
			value.valueUint32 = binary.LittleEndian.Uint32(data[offset : offset+4])
			offset += 4
		case TypeF32:
			value32 := binary.LittleEndian.Uint32(data[offset : offset+4])
			value.valueFloat32 = math.Float32frombits(value32)
			offset += 4
		}

		pkt.values = append(pkt.values, value)
	}

	// Check for terminator
	if offset+2 > len(data) {
		return pkt, fmt.Errorf("data too short for terminator")
	}

	// terminator is CA FE
	if data[offset] != 0xCA || data[offset+1] != 0xFE {
		return pkt, fmt.Errorf("terminator mismatch, expected 0xCAFE, got 0x%X%X", data[offset], data[offset+1])
	}

	return pkt, nil
}

type SensorState uint8

const (
	Off   SensorState = 0x1
	On    SensorState = 0x2
	Error SensorState = 0x3
)

type SensorChannel uint8

const (
	Channel0 SensorChannel = 0x0
	Channel1 SensorChannel = 0x1
	Channel2 SensorChannel = 0x2
	Channel3 SensorChannel = 0x3
	Channel4 SensorChannel = 0x4
	Channel5 SensorChannel = 0x5
	Channel6 SensorChannel = 0x6
	Channel7 SensorChannel = 0x7
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

type SensorType uint8

const (
	Temperature SensorType = 0x0 //
	Humidity    SensorType = 0x1 //
	Pressure    SensorType = 0x2 //
	Light       SensorType = 0x3 //
	CO2         SensorType = 0x4 //
	Presence    SensorType = 0x5 // (float, 0.0-1.0)
	Target      SensorType = 0x6 // (channel index indicates X, Y, Z axis)
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
	case Target:
		return "Target"
	default:
		return fmt.Sprintf("Unknown (%d)", t)
	}
}

type DeviceLocation uint8

const (
	Unknown     DeviceLocation = 0
	Bedroom1    DeviceLocation = 1
	Bedroom2    DeviceLocation = 2
	Bedroom3    DeviceLocation = 3
	Bedroom4    DeviceLocation = 4
	Livingroom1 DeviceLocation = 5
	Livingroom2 DeviceLocation = 6
	Livingroom3 DeviceLocation = 7
	Livingroom4 DeviceLocation = 8
	Kitchen1    DeviceLocation = 9
	Kitchen2    DeviceLocation = 10
	Kitchen3    DeviceLocation = 11
	Kitchen4    DeviceLocation = 12
	Bathroom1   DeviceLocation = 13
	Bathroom2   DeviceLocation = 14
	Bathroom3   DeviceLocation = 15
	Bathroom4   DeviceLocation = 16
	Hallway     DeviceLocation = 17
	ChildARoom  DeviceLocation = 18
	ChildBRoom  DeviceLocation = 19
	ChildCRoom  DeviceLocation = 20
	ChildDRoom  DeviceLocation = 21
	Guest1Room  DeviceLocation = 22
	Guest2Room  DeviceLocation = 23
	Study1Room  DeviceLocation = 24
	Study2Room  DeviceLocation = 25
	Balcony1    DeviceLocation = 26
	Balcony2    DeviceLocation = 27
	Balcony3    DeviceLocation = 28
	Balcony4    DeviceLocation = 29
)

func (d DeviceLocation) String() string {
	switch d {
	case Unknown:
		return "Unknown"
	case Bedroom1:
		return "Bedroom 1"
	case Bedroom2:
		return "Bedroom 2"
	case Bedroom3:
		return "Bedroom 3"
	case Bedroom4:
		return "Bedroom 4"
	case Livingroom1:
		return "Living Room 1"
	case Livingroom2:
		return "Living Room 2"
	case Livingroom3:
		return "Living Room 3"
	case Livingroom4:
		return "Living Room 4"
	case Kitchen1:
		return "Kitchen 1"
	case Kitchen2:
		return "Kitchen 2"
	case Kitchen3:
		return "Kitchen 3"
	case Kitchen4:
		return "Kitchen 4"
	case Bathroom1:
		return "Bathroom 1"
	case Bathroom2:
		return "Bathroom 2"
	case Bathroom3:
		return "Bathroom 3"
	case Bathroom4:
		return "Bathroom 4"
	case Hallway:
		return "Hallway"
	case ChildARoom:
		return "Child A Room"
	case ChildBRoom:
		return "Child B Room"
	case ChildCRoom:
		return "Child C Room"
	case ChildDRoom:
		return "Child D Room"
	case Guest1Room:
		return "Guest 1 Room"
	case Guest2Room:
		return "Guest 2 Room"
	case Study1Room:
		return "Study 1 Room"
	case Study2Room:
		return "Study 2 Room"
	case Balcony1:
		return "Balcony 1"
	case Balcony2:
		return "Balcony 2"
	case Balcony3:
		return "Balcony 3"
	case Balcony4:
		return "Balcony 4"
	}
	return ""
}
