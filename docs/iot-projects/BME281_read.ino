/*
BME280 I2C Test.ino

This code shows how to record data from the BME280 environmental sensor
using I2C interface. This file is an example file, part of the Arduino
BME280 library.

GNU General Public License

Written: Dec 30 2015.
Last Updated: Oct 07 2017.

Connecting the BME280 Sensor:
Sensor              ->  Board
-----------------------------
Vin (Voltage In)    ->  3.3V
Gnd (Ground)        ->  Gnd
SDA (Serial Data)   ->  A4 on Uno/Pro-Mini, 20 on Mega2560/Due, 2 Leonardo/Pro-Micro
SCK (Serial Clock)  ->  A5 on Uno/Pro-Mini, 21 on Mega2560/Due, 3 Leonardo/Pro-Micro

 */

// From: https://github.com/finitespace/BME280/blob/master/examples/BME_280_I2C_Test/BME_280_I2C_Test.ino

#include <BME280I2C.h>
#include <Wire.h>

#define SERIAL_BAUD 115200

BME280I2C sensorBme280; // Default : forced mode, standby time = 1000 ms
                        // Oversampling = pressure ×1, temperature ×1, humidity ×1, filter off,
static uint64_t sensorBme280LastTryInit = 0;
static uint64_t sensorBme280LastRead = 0;
static uint64_t sensorBme280TryInitInterval = 5000 * 1000; // 5 seconds
static uint64_t sensorBme280ReadInterval = 5000 * 1000;    // 5 seconds
static bool sensorBme280Initialized = false;

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
    return
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
  default:
    Serial.println("Found UNKNOWN sensor! Error!");
  }
}

//////////////////////////////////////////////////////////////////
void printBME280Data(uint64_t now, Stream *client)
{
  if (!sensorBme280Initialized)
  {
    initializeSensorBme280();
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
  } // Wait

  Wire.begin(21, 22);

  initializeSensorBme280();
}

//////////////////////////////////////////////////////////////////
void loop()
{
  printBME280Data(&Serial);

  delay(100);
}
