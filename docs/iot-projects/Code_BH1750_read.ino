#include <Wire.h>
#include <BH1750.h>

BH1750 lightMeter;

void setup()
{
    Serial.begin(9600);

    //Wire.begin(8,9);
    Wire.begin(21,22);
    
    lightMeter.configure(BH1750::CONTINUOUS_HIGH_RES_MODE);

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
