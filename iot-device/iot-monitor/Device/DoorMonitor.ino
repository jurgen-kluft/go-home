#include "AZ3166WiFi.h"
#include "AZ3166WiFiUdp.h"
#include "HTS221Sensor.h"
#include "LIS2MDLSensor.h"
#include "LPS22HBSensor.h"
#include "LSM6DSLSensor.h"
#include "OledDisplay.h"

#include "Sensors.h"

#define LOOP_DELAY 10

#define TITLE "Mon - v1.1.13"

static Sensors gSensors;

static bool gHasWifi     = false; // Indicate whether WiFi is ready

// A UDP instance to let us send and receive packets over UDP
static WiFiUDP Udp;
#define LOCAL_PORT 2390
#define REMOTE_IP "10.0.0.22"
#define REMOTE_PORT 7331

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Utilities
static bool InitWiFi()
{
    Screen.print(2, "Connecting...");

    if (WiFi.begin() == WL_CONNECTED)
    {
        IPAddress ip = WiFi.localIP();
        Screen.print(1, ip.get_address());
        Udp.begin(LOCAL_PORT);
        Udp.beginPacket(REMOTE_IP, REMOTE_PORT);
        Screen.print(2, "Running... \r\n");
        return true;
    }
    else
    {
        Screen.print(1, "No Wi-Fi\r\n ");
        return false;
    }
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Arduino sketch
void setup()
{
    Screen.init();
    Screen.print(0, TITLE);

    Screen.print(2, "Initializing...");
    Screen.print(3, " > Serial");
    Serial.begin(115200);

    // Initialize the WiFi module
    Screen.print(2, "Initializing...");
    Screen.print(3, " > WiFi");
    gHasWifi = InitWiFi();
    if (!gHasWifi)
    {
        return;
    }

    //OLEDDisplay* screen = &Screen;
    gSensors.Init(&Screen);

    Screen.print(2, "Monitoring...");
    Screen.print(3, "");
}

static void sendUdpPacket(void* data, int size)
{
    size_t written = 0;
    while (written < size)
    {
        written += Udp.write((const unsigned char*)data + written, size - written);
    }
}

static int         gSensorIndex = 0;
static SensorData  gSensorData[16];
static int         gIteration = 0;
static const char* gProgress  = "..........";

void loop()
{
    if (gHasWifi)
    {
        gSensorData[gSensorIndex].ReadAll(&gSensors);
        gSensorIndex += 1;

        if (gSensorIndex == 16)
        {
            // 16 * 48 bytes = 768 bytes
            sendUdpPacket(&gSensorData[0], 16 * sizeof(SensorData));
            gSensorIndex = 0;
        }

        gIteration = (gIteration + 1) % 128;
        Screen.print(3, gProgress + (gIteration/16));
    }

    delay(LOOP_DELAY);
}
