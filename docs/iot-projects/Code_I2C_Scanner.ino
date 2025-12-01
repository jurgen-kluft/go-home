#include <Wire.h>

void I2C_ScannerWire(TwoWire* wire)
{
    byte error, address;
    int nDevices;

    Serial.println("Scanning...");

    nDevices = 0;
    for (address = 1; address < 127; address++)
    {
        wire->beginTransmission(address);
        error = wire->endTransmission();

        if (error == 0)
        {
            Serial.print("I2C device found at address 0x");
            if (address < 16)
                Serial.print("0");
            Serial.print(address, HEX);
            Serial.println("  !");

            nDevices++;
        }
        else if (error == 4)
        {
            Serial.print("Unknown error at address 0x");
            if (address < 16)
                Serial.print("0");
            Serial.println(address, HEX);
        }
    }
    if (nDevices == 0)
        Serial.println("No I2C devices found\n");
    else
        Serial.println("done\n");
}

void setup()
{
    Serial.begin(115200);

    // Wire.begin(int sda, int scl);
    Wire.begin();        // default: 21, 22
    Wire1.begin(16, 17); // I don't know the default pins !

    Serial.println("---------- Scanning Wire -------------");
    I2C_ScannerWire(&Wire);

    Serial.println("---------- Scanning Wire1 ------------");
    I2C_ScannerWire(&Wire1);
}

void loop() { delay(10); }
