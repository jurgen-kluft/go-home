#include <Wire.h>
#include <BH1750.h>

BH1750 lightMeter;

void setup()
{
    Serial.begin(9600);

    // Initialize the I2C bus (BH1750 library doesn't do this automatically)
    // On esp8266 devices you can select SCL and SDA pins using Wire.begin(D4, D3);
    // Wire.begin(8,9);
    Wire.begin();

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
    float lux = lightMeter.readLightLevel();
    Serial.print("Light: ");
    Serial.print(lux);
    Serial.println(" lx");
    delay(1000);
}
