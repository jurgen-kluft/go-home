#include <Arduino.h>
#include <SensirionI2cScd4x.h>
#include <Wire.h>

// From: https://github.com/Sensirion/arduino-i2c-scd4x/blob/master/examples/exampleUsage/exampleUsage.ino

// macro definitions
// make sure that we use the proper definition of NO_ERROR
#ifdef NO_ERROR
#    undef NO_ERROR
#endif
#define NO_ERROR 0

SensirionI2cScd4x sensorCO2;
bool              sensorCO2Initialized = false;

void PrintUint64(uint64_t &value)
{
    Serial.print("0x");
    Serial.print((uint32_t)(value >> 32), HEX);
    Serial.print((uint32_t)(value & 0xFFFFFFFF), HEX);
}

static char    errorMessage[64];
static int16_t error;

bool printErrorIfError(int16_t error, const char *msg, Stream *output)
{
    if (error != NO_ERROR)
    {
        output->print(msg);
        errorToString(error, errorMessage, sizeof errorMessage);
        output->println(errorMessage);
        return true;
    }
    return false;
}

bool initCO2Sensor()
{
    sensorCO2.begin(Wire, SCD41_I2C_ADDR_62);

    uint64_t serialNumber = 0;
    delay(30);

    // Ensure sensorCO2 is in clean state
    error = sensorCO2.wakeUp();
    if (printErrorIfError(error, "Error trying to execute wakeUp(): ", &Serial))
        return false;

    error = sensorCO2.stopPeriodicMeasurement();
    if (printErrorIfError(error, "Error trying to execute stopPeriodicMeasurement(): ", &Serial))
        return false;

    error = sensorCO2.reinit();
    if (printErrorIfError(error, "Error trying to execute reinit(): ", &Serial))
        return false;

    // Read out information about the sensorCO2
    error = sensorCO2.getSerialNumber(serialNumber);
    if (printErrorIfError(error, "Error trying to execute getSerialNumber(): ", &Serial))
        return false;

    Serial.print("serial number: ");
    PrintUint64(serialNumber);
    Serial.println();

    //
    // If temperature offset and/or sensorCO2 altitude compensation
    // is required, you should call the respective functions here.
    // Check out the header file for the function definitions.
    // Start periodic measurements (5sec interval)
    error = sensorCO2.startPeriodicMeasurement();
    if (printErrorIfError(error, "Error trying to execute startPeriodicMeasurement(): ", &Serial))
        return false;

    return true;
}

//
// Sampling should be at around 1 time per 5 seconds (0.2Hz)
//
void readAndPrintCO2Measurement(Stream *output)
{
    bool     dataReady        = false;
    uint16_t co2Concentration = 0;
    float    temperature      = 0.0;
    float    relativeHumidity = 0.0;

    if (!sensorCO2Initialized)
    {
        sensorCO2Initialized = initCO2Sensor();
    }

    error = sensorCO2.getDataReadyStatus(dataReady);
    if (printErrorIfError(error, "Error trying to execute getDataReadyStatus(): ", output))
        return;

    if (!dataReady)
    {
        output->println("CO2 data not ready");
        return;
    }

    //
    // If ambient pressure compenstation during measurement
    // is required, you should call the respective functions here.
    // Check out the header file for the function definition.
    error = sensorCO2.readMeasurement(co2Concentration, temperature, relativeHumidity);
    if (printErrorIfError(error, "Error trying to execute readMeasurement(): ", output))
        return;

    //
    // Print results in physical units.
    //
    output->print("CO2 concentration [ppm]: ");
    output->print(co2Concentration);
    output->println();
    output->print("Temperature [°C]: ");
    output->print(temperature);
    output->println();
    output->print("Relative Humidity [RH]: ");
    output->print(relativeHumidity);
    output->println();
}

void setup()
{
    Serial.begin(115200);
    while (!Serial)
    {
        delay(100);
    }
    Wire.begin(21, 22);

    initCO2Sensor();

    //
    // If low-power mode is required, switch to the low power
    // measurement function instead of the standard measurement
    // function above. Check out the header file for the definition.
    // For SCD41, you can also check out the single shot measurement example.
    //
}

void loop()
{
    readAndPrintCO2Measurement(&Serial);
    delay(5000);  // Wait for 5 seconds before the next measurement
}
