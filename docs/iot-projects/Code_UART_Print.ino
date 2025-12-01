

// UART Reader

// Reads data byte-by-byte and prints 16 bytes per line in HEX

#define RX_PIN    20      // Adjust to your RX pin
#define TX_PIN    21      // Adjust to your TX pin
#define BAUD_RATE 115200  // Sensor baud rate

HardwareSerial SensorSerial(1);  // Use UART1

void setup()
{
    Serial.begin(115200);
    while (!Serial)
    {
        ;
    }  // Wait for Serial Monitor
    SensorSerial.begin(BAUD_RATE, SERIAL_8N1, RX_PIN, TX_PIN);
    Serial.println("UART HEX Dump Started");
}

void loop()
{
    static int  count = 0;
    static char text[20];

    if (SensorSerial.available())
    {
        uint8_t byteData = SensorSerial.read();

        // Print byte in HEX (two digits)
        if (byteData < 0x10)
            Serial.print("0");  // Leading zero for single digit
        Serial.print(byteData, HEX);
        Serial.print(" ");

        text[count] = (char)byteData;
        count++;
        if (count >= 16)
        {
            Serial.println();  // New line after 16 bytes
            text[count] = 0;
            Serial.println(text);
            count = 0;
        }
    }
}
