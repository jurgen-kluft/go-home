// From: https://github.com/finitespace/BME280/blob/master/examples/BME_280_I2C_Test/BME_280_I2C_Test.ino

#include <BME280I2C.h>
#include <Wire.h>

#define SERIAL_BAUD 115200

BME280I2C sensorBme280;  // Default : forced mode, standby time = 1000 ms
                         // Oversampling = pressure ×1, temperature ×1, humidity ×1, filter off,
uint64_t sensorBme280LastTryInit     = 0;
uint64_t sensorBme280LastRead        = 0;
uint64_t sensorBme280TryInitInterval = 5 * 1000 * 1000;  // 5 seconds
uint64_t sensorBme280ReadInterval    = 5 * 1000 * 1000;  // 5 seconds
bool     sensorBme280Initialized     = false;

// ESP32 YD
// State: ?
const int sdaPin = 21;
const int sclPin = 22;

// ESP32S3 Dev Module
// State: Working
//const int sdaPin = 8;
//const int sclPin = 9;

void initializeSensorBme280(uint64_t now)
{
    if (now - sensorBme280LastTryInit < sensorBme280TryInitInterval)
    {
        return;
    }
    sensorBme280LastTryInit = now;

    if (!sensorBme280.begin())
    {
        Serial.println("Could not find BME280 sensor!");
        sensorBme280Initialized = false;
        return;
    }

    switch (sensorBme280.chipModel())
    {
        case BME280::ChipModel_BME280:
            Serial.println("Found BME280 sensor! Success.");
            sensorBme280Initialized = true;
            break;
        case BME280::ChipModel_BMP280:
            Serial.println("Found BMP280 sensor! No Humidity available.");
            sensorBme280Initialized = true;
            break;
        default: Serial.println("Found UNKNOWN sensor! Error!");
    }
}

//////////////////////////////////////////////////////////////////
void printBME280Data(uint64_t now, Stream *client)
{
    if (!sensorBme280Initialized)
    {
        initializeSensorBme280(now);
        return;
    }

    if (now - sensorBme280LastRead < sensorBme280ReadInterval)
    {
        return;
    }
    sensorBme280LastRead = now;

    float temp(NAN), hum(NAN), pres(NAN);

    BME280::TempUnit tempUnit(BME280::TempUnit_Celsius);
    BME280::PresUnit presUnit(BME280::PresUnit_Pa);

    sensorBme280.read(pres, temp, hum, tempUnit, presUnit);

    client->print("Temp: ");
    client->print(temp);
    client->print("°" + String(tempUnit == BME280::TempUnit_Celsius ? 'C' : 'F'));
    client->print("\t\tHumidity: ");
    client->print(hum);
    client->print("% RH");
    client->print("\t\tPressure: ");
    client->print(pres);
    client->println("Pa");
}

//////////////////////////////////////////////////////////////////
void setup()
{
    Serial.begin(SERIAL_BAUD);
    while (!Serial)
    {
        delay(100); // Wait for serial port to connect. Needed for native USB
    }  

    Wire.begin(sdaPin, sclPin);

    uint64_t now = micros();

    sensorBme280LastTryInit     = now;
    sensorBme280LastRead        = now;

    initializeSensorBme280(now);
}

//////////////////////////////////////////////////////////////////
void loop()
{
    uint64_t now = micros();
    printBME280Data(now, &Serial);

    delay(100);
}
