package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

// Bluetooth instance, tracking of Beacon instances
type Bluetooth struct {
	beacons map[string]*Beacon
}

// Create returns an instance of Bluetooth that is responsible for holding on the Beacon instances
func Create() *Bluetooth {
	ble := &Bluetooth{beacons: make(map[string]*Beacon)}
	return ble
}

// Beacon is an instance that holds information of a BLE device
type Beacon struct {
	UUID  string
	Major uint16
	Minor uint16
	RSSI  int
}

// IsValidBeacon will return true when the packet is large enough and contains the correct ID
func (bt *Bluetooth) IsValidBeacon(data []byte) bool {
	if len(data) < 25 || binary.BigEndian.Uint32(data) != 0x4c000215 {
		return false
	}
	return true
}

// HaveSeenBeacon returns true if bluetooth has seen this beacon before and thus has a valid beacon instance
func (bt *Bluetooth) HaveSeenBeacon(data []byte) bool {
	if len(data) < 25 || binary.BigEndian.Uint32(data) != 0x4c000215 {
		return false
	}

	uuid := strings.ToUpper(hex.EncodeToString(data[4:8]) + "-" + hex.EncodeToString(data[8:10]) + "-" + hex.EncodeToString(data[10:12]) + "-" + hex.EncodeToString(data[12:14]) + "-" + hex.EncodeToString(data[14:20]))

	_, exists := bt.beacons[uuid]
	return exists
}

// NewBeacon returns an instance of Beacon when incoming data is recognized as a beacon, otherwise a nil is returned
func (bt *Bluetooth) NewBeacon(data []byte) (*Beacon, error) {
	if bt.IsValidBeacon(data) {
		uuid := strings.ToUpper(hex.EncodeToString(data[4:8]) + "-" + hex.EncodeToString(data[8:10]) + "-" + hex.EncodeToString(data[10:12]) + "-" + hex.EncodeToString(data[12:14]) + "-" + hex.EncodeToString(data[14:20]))
		beacon, exists := bt.beacons[uuid]
		if !exists {
			beacon = new(Beacon)
			beacon.UUID = uuid
			beacon.Major = binary.BigEndian.Uint16(data[20:22])
			beacon.Minor = binary.BigEndian.Uint16(data[22:24])
			beacon.RSSI = 0
			bt.beacons[uuid] = beacon
		}
		return beacon, nil
	}
	return nil, errors.New("Not an iBeacon")
}

// GetExistingBeacon returns a beacon instance if it already has been created before, otherwise nill
func (bt *Bluetooth) GetExistingBeacon(data []byte) *Beacon {
	if bt.IsValidBeacon(data) {
		uuid := strings.ToUpper(hex.EncodeToString(data[4:8]) + "-" + hex.EncodeToString(data[8:10]) + "-" + hex.EncodeToString(data[10:12]) + "-" + hex.EncodeToString(data[12:14]) + "-" + hex.EncodeToString(data[14:20]))
		beacon, exists := bt.beacons[uuid]
		if exists {
			return beacon
		}
	}
	return nil
}
