#ifndef __DOOR_MONITOR_SENSORS_H__
#define __DOOR_MONITOR_SENSORS_H__

#ifdef USE_PRAGMA_ONCE 
#pragma once 
#endif

struct Sensors;

struct SensorData
{
    float Temperature;
    float Humidity;
    float Pressure;
    int   Magnetic[3];
    int   Accelerator[3];
    int   Gyroscope[3];

    bool ReadAll(Sensors* sensors);

    bool ReadTemperature(Sensors* sensors);
    bool ReadHumidity(Sensors* sensors);
    bool ReadPressure(Sensors* sensors);
    bool ReadAcceleration(Sensors* sensors);
    bool ReadGyroscope(Sensors* sensors);
    bool ReadMagnetic(Sensors* sensors);
};

// Forward declares
//struct DevI2C;
//struct LIS2MDLSensor;
//struct HTS221Sensor;
//struct LPS22HBSensor;
//struct LSM6DSLSensor;

struct Sensors
{
    DevI2C*        m_i2c;
    LIS2MDLSensor* m_lis2mdl;
    HTS221Sensor*  m_hts221sensor;
    LPS22HBSensor* m_lps22hbsensor;
    LSM6DSLSensor* m_lsm6dslsensor;

    void Init(OLEDDisplay* screen);
};


#endif