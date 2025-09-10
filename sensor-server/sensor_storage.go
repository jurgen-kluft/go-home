package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

// Storage for sensor data
// One storage stream for [DeviceLocation | SensorType]
// DeviceLocation is associated with

type SensorStore struct {
	Name  string // 'LivingRoomA_Temperature'
	Id    int    // Unique ID for this sensor
	Year  uint16 // Current year
	Month uint8  // Current month
	Day   uint8  // Current day of month
}

// DataBlock format:
//   - Year (int)
//   - Month (int)
//   - Day (int)
//   - Hour (int)
//   - SampleType (int) - bit, int8, int16, int32
//   - SampleFreq (int) - samples per hour
//   - DataLength (int) - length of data in bytes
//   - Data (variable length, depends on SampleType and SampleFreq)
type SensorDataBlock struct {
	Info         *SensorStore
	Time         time.Time       // Time (Year, Month, Day, at zero hour)
	SampleType   SensorFieldType // Type of samples in this block
	SampleFreq   int32           // Samples per hour
	SamplePeriod int32           // Milliseconds between samples
	SampleCount  int32           // Number of samples in this block
	LastSample   SensorValue     // Last sample value
	Buffer       []byte          // Buffer to hold the samples
}

const (
	SensorDataBlockHeaderSize = 64 // Size of the header in bytes
)

func NewSensorDataBlock(info *SensorStore, sampleType SensorFieldType, sampleFreq int32) *SensorDataBlock {
	// Construct the correct time for the start of this block
	t := time.Now()
	blockTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	// TODO This datablock might already exist on disk, load it instead of creating a new one.
	//      This can happen when the server crashes and/or restarts during the day.

	blockSizeInBytes := (int(sampleFreq*24)*int(sampleType.SizeInBits()) + 7) / 8
	blockSizeInBytes = blockSizeInBytes + SensorDataBlockHeaderSize
	block := &SensorDataBlock{
		Info:       info,
		Time:       blockTime,
		SampleType: sampleType,
		SampleFreq: sampleFreq,
		Buffer:     make([]byte, blockSizeInBytes),
	}

	binary.LittleEndian.PutUint32(block.Buffer, uint32(blockTime.Year()))
	binary.LittleEndian.PutUint32(block.Buffer[4:], uint32(blockTime.Month()))
	binary.LittleEndian.PutUint32(block.Buffer[8:], uint32(blockTime.Day()))
	binary.LittleEndian.PutUint32(block.Buffer[12:], uint32(block.SampleType))
	binary.LittleEndian.PutUint32(block.Buffer[16:], uint32(block.SampleFreq))
	binary.LittleEndian.PutUint32(block.Buffer[20:], uint32(len(block.Buffer)-SensorDataBlockHeaderSize))

	return block
}

func (s *SensorDataBlock) WriteSensorValue(sensorValue SensorValue) {
	s.LastSample = sensorValue
	if sensorValue.IsZero() {
		return // Default value
	}

	sampleTime := time.Now()
	sampleIndex := sampleTime.Sub(s.Time).Milliseconds() / (60 * 60 * 1000 / int64(s.SampleFreq))
	if sampleIndex < 0 || sampleIndex >= int64(s.SampleFreq*24) {
		// Sample is out of range for this block, ignore it
		return
	}

	s.SampleCount = int32(sampleIndex) + 1

	// Write the sensor value to the buffer based on its type
	switch s.SampleType {
	case TypeBit:
		byteIndex := SensorDataBlockHeaderSize + (sampleIndex / 8)
		bitIndex := sampleIndex % 8
		s.Buffer[byteIndex] |= (1 << bitIndex) // Set the bit
	case TypeS8:
		byteIndex := SensorDataBlockHeaderSize + sampleIndex
		s.Buffer[byteIndex] = byte(sensorValue.value)
	case TypeS16:
		byteIndex := SensorDataBlockHeaderSize + sampleIndex*2
		binary.LittleEndian.PutUint16(s.Buffer[byteIndex:byteIndex+2], uint16(sensorValue.value))
	case TypeS32:
		byteIndex := SensorDataBlockHeaderSize + sampleIndex*4
		binary.LittleEndian.PutUint32(s.Buffer[byteIndex:byteIndex+4], uint32(sensorValue.value))
	case TypeU8:
		byteIndex := SensorDataBlockHeaderSize + sampleIndex
		s.Buffer[byteIndex] = byte(sensorValue.value)
	case TypeU16:
		byteIndex := SensorDataBlockHeaderSize + sampleIndex*2
		binary.LittleEndian.PutUint16(s.Buffer[byteIndex:byteIndex+2], uint16(sensorValue.value))
	case TypeU32:
		byteIndex := SensorDataBlockHeaderSize + sampleIndex*4
		binary.LittleEndian.PutUint32(s.Buffer[byteIndex:byteIndex+4], uint32(sensorValue.value))
	}
}

func (s *SensorDataBlock) IsDone() bool {
	// Determine if the block is full based on SampleFreq and time
	return s.SampleCount >= s.SampleFreq
}

type SensorStorage struct {
	dirPath           string                // Directory path for storing sensor data files
	sensorNameList    []string              // List of sensor identifiers
	sensorTypeList    []SensorType          // List of sensor types
	sensorFreqList    []int32               // Sample frequency per sensor type
	sensorNameToId    map[string]int        // Map from sensor name to ID
	streams           []*SensorDataBlock    // Active data streams
	writeChannel      chan *SensorDataBlock // Channel for writing data blocks
	blockHeaderBuffer bytes.Buffer
	blockHeaderWriter io.Writer
}

func NewSensorStorage(dirPath string, writeChannel chan *SensorDataBlock) *SensorStorage {
	storage := &SensorStorage{
		dirPath:        dirPath,
		sensorNameList: make([]string, 0),
		sensorTypeList: make([]SensorType, 0),
		sensorFreqList: make([]int32, 0),
		sensorNameToId: make(map[string]int),
		streams:        make([]*SensorDataBlock, 0),
		writeChannel:   writeChannel,
	}
	storage.blockHeaderWriter = &storage.blockHeaderBuffer

	// The go-routine that writes data blocks to the file system
	go func() {
		for block := range writeChannel {
			err := storage.WriteDataBlock(block)
			if err != nil {
				fmt.Printf("Error writing data block: %v\n", err)
			}
		}
	}()

	return storage
}

// RegisterSensor registers a new sensor and returns its ID.
// Note: 'sensorName' should be unique.
func (s *SensorStorage) RegisterSensor(sensorName string, sensorType SensorType) int {
	if id, exists := s.sensorNameToId[sensorName]; exists {
		return id // Sensor already registered
	}

	id := len(s.sensorNameList)
	s.sensorNameList = append(s.sensorNameList, sensorName)

	// For simplicity, assume sensor type and frequency are derived from the name
	var sensorFreq int32

	s.sensorTypeList = append(s.sensorTypeList, sensorType)
	s.sensorFreqList = append(s.sensorFreqList, sensorFreq)
	s.sensorNameToId[sensorName] = id

	// Ensure streams slice is large enough
	for len(s.streams) <= id {
		s.streams = append(s.streams, nil)
	}

	return id
}
func (s *SensorStorage) WriteSensorValue(id int, sensorValue SensorValue) error {
	// Create or get the data block for the given location and sensor type
	block := s.streams[id]
	if block == nil {
		sensorStore := &SensorStore{Name: s.sensorNameList[id], Id: id}
		block = NewSensorDataBlock(sensorStore, SensorFieldType(s.sensorTypeList[id]), s.sensorFreqList[id])
		s.streams[id] = block
	}

	// Write the sensor value to the block's buffer
	block.WriteSensorValue(sensorValue)
	if block.IsDone() {
		// If the block is full, send it to the write channel and create a new block
		s.writeChannel <- block
		s.streams[id] = nil // Clear the current block to create a new one next time
	}

	return nil
}

// The go-routine for writing data blocks to the file system
func (s *SensorStorage) WriteDataBlock(block *SensorDataBlock) error {
	// Check if the block is valid and has data
	if block == nil || len(block.Buffer) == 0 {
		return nil // Nothing to write
	}

	// Open the file for writing, append-only
	sensorStorePath := path.Join(s.dirPath, block.Info.Name, fmt.Sprintf("sensor_data_%04d_%02d_%02d.dat", block.Time.Year(), block.Time.Month(), block.Time.Day()))
	file, err := os.OpenFile(sensorStorePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the header and data to the file
	_, err = file.Write(block.Buffer)

	return err
}
