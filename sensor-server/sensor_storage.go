package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

// Storage for sensor data
// One storage stream for [DeviceLocation | SensorType]

type SensorDataFile struct {
	Filename   string
	FileHandle *os.File
}

// DataBlock, we can write byte, uint16, uint32, float32, etc. to this block
type SensorDataBlock struct {
	Info   *SensorDataFile
	Year   uint16
	Month  uint16
	Day    uint16
	Hour   uint16
	Buffer *bytes.Buffer
	Writer *io.Writer
}

type SensorStorage struct {
	blocks            map[*SensorDataFile]*SensorDataBlock // Free data blocks
	streams           map[uint32]*SensorDataBlock          //
	writeChannel      chan *SensorDataBlock                // Channel for writing data blocks
	blockHeaderBuffer bytes.Buffer
	blockHeaderWriter io.Writer
}

func (s *SensorStorage) StoreSensor(location SensorLocation, sensorType SensorType, sensorValue SensorValue) error {
	// Create or get the data block for the given location and sensor type
	dataFile := &SensorDataFile{Filename: string(location) + "_" + sensorType.String() + ".dat"}
	block, exists := s.blocks[dataFile]
	if !exists {
		block = &SensorDataBlock{
			Info:   dataFile,
			Buffer: new(bytes.Buffer),
		}
		s.blocks[dataFile] = block
	}

	// Write the sensor value to the block's buffer
	if err := binary.Write(block.Buffer, binary.LittleEndian, sensorValue); err != nil {
		return err
	}

	// Send the block to the write channel for processing
	s.writeChannel <- block

	return nil
}

// The go-routine for writing data blocks to the file system
func (s *SensorStorage) WriteDataBlock(block *SensorDataBlock) error {
	// Check if the block is valid and has data
	if block == nil || block.Buffer.Len() == 0 {
		return nil // Nothing to write
	}

	// Open the file for writing, append-only
	file, err := os.OpenFile(block.Info.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	binary.Write(&s.blockHeaderBuffer, binary.LittleEndian, block.Year)
	binary.Write(&s.blockHeaderBuffer, binary.LittleEndian, block.Month)
	binary.Write(&s.blockHeaderBuffer, binary.LittleEndian, block.Day)
	binary.Write(&s.blockHeaderBuffer, binary.LittleEndian, block.Hour)

	blockData := block.Buffer.Bytes()
	binary.Write(&s.blockHeaderBuffer, binary.LittleEndian, len(blockData))

	// Write the header and data to the file
	_, err = file.Write(s.blockHeaderBuffer.Bytes())
	_, err = file.Write(blockData)

	// Reset the block after writing
	s.blockHeaderBuffer.Reset()
	block.Buffer.Reset()

	return err
}
