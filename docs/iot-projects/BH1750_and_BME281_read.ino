#include <Wire.h>
#include <BH1750.h>
#include <BME280I2C.h>

#define SERIAL_BAUD 115200

BME280I2C bme;  // Default : forced mode, standby time = 1000 ms
                // Oversampling = pressure ×1, temperature ×1, humidity ×1, filter off,

BH1750 lightMeter;

void setup()
{
    Serial.begin(SERIAL_BAUD);

    // Initialize the I2C bus (BH1750 library doesn't do this automatically)
    // On esp8266 devices you can select SCL and SDA pins using Wire.begin(D4, D3);
    // Wire.begin(8,9);
    Wire.begin();

    while (!bme.begin())
    {
        Serial.println("Could not find BME280 sensor!");
        delay(1000);
    }

    switch (bme.chipModel())
    {
        case BME280::ChipModel_BME280: Serial.println("Found BME280 sensor! Success."); break;
        case BME280::ChipModel_BMP280: Serial.println("Found BMP280 sensor! No Humidity available."); break;
        default: Serial.println("Found UNKNOWN sensor! Error!");
    }

    if (lightMeter.begin(BH1750::CONTINUOUS_HIGH_RES_MODE, 0x23, &Wire))
    {
        Serial.println(F("BH1750 initialised"));
    }
    else
    {
        Serial.println(F("Error initialising BH1750"));
    }
}

void loop()
{
    printBME280Data(&Serial);
    printLightSensor(&Serial);
    delay(1000);
}

//////////////////////////////////////////////////////////////////
void printBME280Data(Stream *client)
{
    float temp(NAN), hum(NAN), pres(NAN);

    BME280::TempUnit tempUnit(BME280::TempUnit_Celsius);
    BME280::PresUnit presUnit(BME280::PresUnit_Pa);

    bme.read(pres, temp, hum, tempUnit, presUnit);

    client->print("Temp: ");
    client->print(temp);
    client->print("°" + String(tempUnit == BME280::TempUnit_Celsius ? 'C' : 'F'));
    client->print("\t\tHumidity: ");
    client->print(hum);
    client->print("% RH");
    client->print("\t\tPressure: ");
    client->print(pres);
    client->println("Pa");

    delay(1000);
}

//////////////////////////////////////////////////////////////////
void printLightSensor(Stream *client)
{
    float lux = lightMeter.readLightLevel();
    client->print("Light: ");
    client->print(lux);
    client->println(" lx");
}
