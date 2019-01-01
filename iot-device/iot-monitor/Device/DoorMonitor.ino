// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license.
// To get started please visit
// https://microsoft.github.io/azure-iot-developer-kit/docs/projects/door-monitor?utm_source=ArduinoExtension&utm_medium=ReleaseNote&utm_campaign=VSCode
#include "AZ3166WiFi.h"
#include "AZ3166WiFiUdp.h"
#include "HTS221Sensor.h"
#include "LIS2MDLSensor.h"
#include "LPS22HBSensor.h"
#include "LSM6DSLSensor.h"
#include "OledDisplay.h"

#define LOOP_DELAY 10

// The magnetometer sensor
static DevI2C *i2c;
static LIS2MDLSensor *lis2mdl;
static HTS221Sensor *hts221sensor;
static LPS22HBSensor *lps22hbsensor;
static LSM6DSLSensor *lsm6dslsensor;


// Indicate whether the magnetometer sensor has been initialized
static bool initialized = false;

// Indicate whether WiFi is ready
static bool hasWifi = false;

// A UDP instance to let us send and receive packets over UDP
static WiFiUDP Udp;
unsigned int localPort = 2390;                  // local port to listen for UDP packets

// 48 bytes
struct SensorData {
  float Temperature;
  float Humidity;
  float Pressure;
  int Magnetic[3];
  int Accelerator[3];
  int Gyroscope[3];
};

void sendSensorpacket(char* address, uint16_t port, void* data, int size);

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Audio, RMS and Gain
float calcGain16LE(char *buf, uint16_t len);
float calcRMS16LE(char *buf, uint16_t len);

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Utilities
static void InitWiFi() {
  Screen.print(2, "Connecting...");

  if (WiFi.begin() == WL_CONNECTED) {
    IPAddress ip = WiFi.localIP();
    Screen.print(1, ip.get_address());
    hasWifi = true;
    Udp.begin(localPort);
    Udp.beginPacket("10.0.0.22", 7331);
    Screen.print(2, "Running... \r\n");
  } else {
    hasWifi = false;
    Screen.print(1, "No Wi-Fi\r\n ");
  }
}

static void InitHTS221Sensor() {
  // init the sensor
  hts221sensor = new HTS221Sensor(*i2c);
  hts221sensor->init(NULL);
  hts221sensor->enable();
  hts221sensor->reset();
}

static bool readTemperature(SensorData* sd) {
    sd->Temperature = 0;
    int res = hts221sensor->getTemperature(&sd->Temperature);
    return res == 0;
}

static bool readHumidity(SensorData* sd) {
    sd->Humidity = 0;
    int res = hts221sensor->getHumidity(&sd->Humidity);
    return res == 0;
}

static void InitLPS22HBSensor() {
    lps22hbsensor = new LPS22HBSensor(*i2c);
    lps22hbsensor->init(NULL);    
}

static bool readPressure(SensorData* sd) {
    sd->Pressure = -1.0;
    int res = lps22hbsensor->getPressure(&sd->Pressure);
    return res == 0;
}

static void InitLSM6DSLSensor() {
    lsm6dslsensor = new LSM6DSLSensor(*i2c, D4, D5);
    lsm6dslsensor->init(NULL);
    lsm6dslsensor->enableAccelerator();
    lsm6dslsensor->enableGyroscope();
}

static bool readAcceleration(SensorData* sd) {
    int res = lsm6dslsensor->getXAxes(sd->Accelerator);  
    return res == 0;
}

static bool readGyroscope(SensorData* sd) {
    int res = lsm6dslsensor->getGAxes(sd->Gyroscope);
    return res == 0;
}

static void InitMagnetometer() {
  Screen.print(2, "Initializing...");
  i2c = new DevI2C(D14, D15);
  lis2mdl = new LIS2MDLSensor(*i2c);
  lis2mdl->init(NULL);
}

static bool readMagnetic(SensorData* sd) {
    int res = lis2mdl->getMAxes(sd->Magnetic);
    return res == 0;
}


static void InitSensors() {
  Screen.print(2, "Initializing...");

  Screen.print(3, " > Magnetometer");
  InitMagnetometer();

  Screen.print(3, " > HTS221Sensor");
  InitHTS221Sensor();

  Screen.print(3, " > LSM6DSLSensor");
  InitLSM6DSLSensor();

  Screen.print(3, " > LPS22HBSensor");
  InitLPS22HBSensor();
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Arduino sketch
void setup() {
  Screen.init();
  Screen.print(0, "Mon - v1.1.12");

  Screen.print(2, "Initializing...");
  Screen.print(3, " > Serial");
  Serial.begin(115200);

  // Initialize the WiFi module
  Screen.print(2, "Initializing...");
  Screen.print(3, " > WiFi");
  hasWifi = false;
  InitWiFi();
  if (!hasWifi) {
    return;
  }

  InitSensors();

  Screen.print(2, "Monitoring...");
  Screen.print(3, ">>");
}

// send an NTP request to the time server at the given address
void sendSensorpacket(void* data, int size) {
    size_t written = 0;
    while (written < size) {
      written += Udp.write((const unsigned char*)data + written, size - written);
    }
}

static SensorData sensors;
static int iteration = 0;
static const bool verbose = false;

void loop() {
  if (hasWifi) {    
    readTemperature(&sensors);
    readHumidity(&sensors);
    readPressure(&sensors);
    readMagnetic(&sensors);
    readAcceleration(&sensors);
    readGyroscope(&sensors);
    sendSensorpacket(&sensors, sizeof(SensorData));
  }
  iteration+=1;
  if (iteration == 1)
  {
    Screen.print(3, "> ");
  }
  else if (iteration == 25)
  {
    Screen.print(3, " > ");
  }
  else if (iteration == 50)
  {
    Screen.print(3, "  > ");
  }
  else if (iteration == 75)
  {
    Screen.print(3, "   > ");
  }
  else if (iteration == 100)
  {
    iteration = 0;
  }
  delay(LOOP_DELAY);
}
