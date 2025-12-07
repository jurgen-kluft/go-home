package cadlib

// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// Global constants
const (
	Unit         = 1.0    // 1 unit = 1 mm
	WT           = 1.6    // Wall Thickness
	W1           = WT     // Thin Wall
	W2           = 2 * WT // Medium Wall
	Rounding     = 0.5    // Small Rounding
	MainRounding = 10     // Main Rounding
)

// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// Power Blocks (USB power supplies)

// Power block dimensions
const (
	PowerBlockWidth    = 36.0
	PowerBlockLength   = 46.52
	PowerBlockHeight   = 23.0 + 0.25
	PowerBlockRounding = 2.5
)

// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// Relay boards
const (
	RelayBoardW            = 17.0
	RelayBoardL            = 57.0
	RelayBoardH            = 17.0 - RelayBoardPCBThickness
	RelayBoardPCBThickness = 1.6
	RelayBoardBottomHeight = 3.0
)

// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// USB, Antenna

// External antenna dimensions
const (
	AntennaW = 20.0
	AntennaL = 40.0
	AntennaH = 1.5
)

// USB-A dimensions
const (
	USBAWidth  = 12.0
	USBALength = 16.0
	USBAHeight = 7.0
)

// USB-C hole dimensions (The USB-C plug connector)
const (
	UsbCHoleDiameter     = 11.6 // Radius of the USB-C hole, 1.16 cm
	UsbCHoleRadius       = UsbCHoleDiameter / 2.0
	UsbCHoleRingDiameter = 17.0
)

// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// Microcontroller boards

// Seeed studio XIAO ESP32-C3 dimensions
// USB-C connector
const (
	ESP32C3W = 18.0
	ESP32C3L = 21.0
	ESP32C3H = 5
	ESP32C3T = 1.25 // Thickness of the PCB
)

// LilyGo ETH PCB board measurements
const (
	LiliGo_PCB_W          = 28.0
	LiliGo_PCB_L          = 59.5
	LiliGo_PCB_MountingHR = 1.5  // Radius of the mounting holes
	LiliGo_PCB_H2HL       = 54.0 // Distance between the mounting holes along the length
	LiliGo_PCB_H2HW       = 23.0 // Distance between the mounting holes along the width
)

// Waveshare ETH PCB board measurements
const (
	WS_PCB_W          = 21.0
	WS_PCB_L          = 72.8
	WS_PCB_MountingHR = 1.6   // Radius of the mounting holes
	WS_PCB_H2HL       = 54.15 // Distance between the mounting holes along the length
	WS_PCB_H2HW       = 18.25 // Distance between the mounting holes along the width
)

// ESP8266 NodeMcu board dimensions
// TODO: Verify the dimensions
const (
	ESP8266W    = 30.4
	ESP8266L    = 57.0
	ESP8266H    = 14.0
	ESP8266T    = 1.6  // Thickness of the PCB
	ESP8266BT   = 3.0  // Bottom height (height of components on bottom side
	ESP8266H2HW = 25.4 // Widht, mounting hole to mounting hole
	ESP8266H2HL = 52.0 // Length, mounting hole to mounting hole
	ESP8266HR   = 3.0  // Radius of the mounting holes
)

// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// DISPLAYS

// Display, 128x128, SH1107 OLED
const (
	Sh1107ScreenW       = 37.3 // Actual display width
	Sh1107ScreenL       = 34.0 // Actual display length
	Sh1107ScreenR       = 0.5  // Actual display corner rounding
	Sh1107W             = 47.1 // Overall width including bezel
	Sh1107L             = 34.1 // Overall length including bezel
	Sh1107MountingW     = 42   // Mounting hole to hole width
	Sh1107MountingL     = 29   // Mounting hole to hole length
	Sh1107MountingHoleD = 2.2  // Mounting hole diameter
)

// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// SENSORS

// RD03D mmWave sensor dimensions
const (
	RD03DW      = 15.25
	RD03DL      = 44.5
	RD03DHeight = 5.0  // Actual measured height
	RD03DH      = 11.0 // Enclosed height
	RD03DT      = 3.2  // PCB and things thickness
	RD03DBT     = 0.7  // Bottom height (height of components on bottom side
	RD03DTT     = 1.0  // Bottom height (height of components on bottom side
)

// APDS9960 sensor
const (
	APDS9960W      = 4.0
	APDS9960L      = 3.0
	APDS9960Height = 1.0  // Actual measured height
	APDS9960H      = 11.0 // Enclosed height
	APDS9960T      = 1.6  // PCB and things thickness
	APDS9960BT     = 0.5  // Bottom height (height of components on bottom side
)

// Scd41 (CO2, temp, humidity)
const (
	Scd41SensorWidth          = 8.5                                              // Width of the sensor body
	Scd41SensorLength         = 8.5                                              // Length of the sensor body
	Scd41SensorHeight         = 6.5                                              // Height of the sensor body
	Scd41BottomToSensorTop    = 18.0                                             // Length until the top of the sensor body (measured from the bottom)
	Scd41BottomToSensorMiddle = Scd41BottomToSensorTop - (Scd41SensorHeight / 2) // Length until the middle of the sensor body (measured from the bottom)
	Scd41W                    = 13.2                                             // Width of the PCB
	Scd41L                    = 22.0                                             // Length of the PCB
	Scd41H                    = 8
	Scd41T                    = 1.6
)

// Bh1750 (light)
const (
	Bh1750Spacing              = 1.0
	Bh1750BottomToSensorTop    = 12.25                                              // Length until the top of the sensor body (measured from the bottom)
	Bh1750TopToSensorBottom    = Bh1750L - 8.25                                     // Length until the bottom of the sensor body (measured from the top)
	Bh1750BottomToSensorMiddle = Bh1750BottomToSensorTop - (Bh1750SensorHeight / 2) // Length until the middle of the sensor body (measured from the bottom)
	Bh1750SensorWidth          = 3.2                                                // Width of the sensor body
	Bh1750SensorHeight         = Bh1750BottomToSensorTop - Bh1750TopToSensorBottom  // Height of the sensor body

	Bh1750W      = 14.2 // Width of the PCB
	Bh1750L      = 18.5 // Length of the PCB
	Bh1750Height = 1.0  // Height of the PCB
	Bh1750H      = 2.0
	Bh1750T      = 1.6
)

// Bme280 (temperature, humidity, pressure)
const (
	Bme280W                    = 10.5
	Bme280L                    = 13.2
	Bme280H                    = 3.5
	Bme280T                    = 1.6
	Bme280SensorWidth          = 2.0
	Bme280SensorLength         = 3.0
	Bme280SensorHeight         = 1.0
	Bme280BottomToSensorMiddle = 11.0
)
