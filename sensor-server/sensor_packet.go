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
	sensorType  SensorType
	sensorModel SensorModel
	state       SensorState
	fieldType   FieldType
	value       float64
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
			sensorType:  SensorType(data[offset] >> 4),
			sensorModel: SensorModel(data[offset] & 0x0F),
			state:       SensorState(data[offset+1] >> 4),
			fieldType:   FieldType(data[offset+1] & 0x0F),
		}

		offset += 2
		// depending on fieldType, read the appropriate value.
		// the written values are in little-endian format
		switch value.fieldType {
		case TypeS8:
			value.value = float64(data[offset])
			offset += 1
		case TypeS16:
			value.value = float64(binary.LittleEndian.Uint16(data[offset : offset+2]))
			offset += 2
		case TypeS32:
			value.value = float64(int32(binary.LittleEndian.Uint32(data[offset : offset+4])))
			offset += 4
		case TypeU8:
			value.value = float64(data[offset])
			offset += 1
		case TypeU16:
			value.value = float64(binary.LittleEndian.Uint16(data[offset : offset+2]))
			offset += 2
		case TypeU32:
			value.value = float64(binary.LittleEndian.Uint32(data[offset : offset+4]))
			offset += 4
		case TypeF32:
			value32 := binary.LittleEndian.Uint32(data[offset : offset+4])
			value.value = float64(math.Float32frombits(value32))
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
