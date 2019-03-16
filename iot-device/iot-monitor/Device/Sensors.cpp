#include "HTS221Sensor.h"
#include "LIS2MDLSensor.h"
#include "LPS22HBSensor.h"
#include "LSM6DSLSensor.h"
#include "OledDisplay.h"

#include "Sensors.h"

void Sensors::Init(OLEDDisplay* screen)
{
    screen->print(2, "Initializing...");
    m_i2c = new DevI2C(D14, D15);

    screen->print(3, " > Magnetometer");
    m_lis2mdl = new LIS2MDLSensor(*m_i2c);
    m_lis2mdl->init(NULL);

    screen->print(3, " > HTS221Sensor");
    m_hts221sensor = new HTS221Sensor(*m_i2c);
    m_hts221sensor->init(NULL);
    m_hts221sensor->enable();
    m_hts221sensor->reset();

    screen->print(3, " > LSM6DSLSensor");
    m_lsm6dslsensor = new LSM6DSLSensor(*m_i2c, D4, D5);
    m_lsm6dslsensor->init(NULL);
    m_lsm6dslsensor->enableAccelerator();
    m_lsm6dslsensor->enableGyroscope();

    screen->print(3, " > LPS22HBSensor");
    m_lps22hbsensor = new LPS22HBSensor(*m_i2c);
    m_lps22hbsensor->init(NULL);
}

void SensorData::Init()
{
    this->HDR = 0xF00D;
    Length = sizeof(int) + 3 * sizeof(int) + 3 * 3 * sizeof(int);
}

bool SensorData::ReadAll(Sensors* sensors)
{
    int res = 0;

    this->Temperature = 0;
    this->Humidity = 0;
    this->Pressure = -1.0;
    
    res += sensors->m_hts221sensor->getTemperature(&this->Temperature);
    res += sensors->m_hts221sensor->getHumidity(&this->Humidity);

    res += sensors->m_lps22hbsensor->getPressure(&this->Pressure);

    res += sensors->m_lsm6dslsensor->getXAxes(this->Accelerator);
    res += sensors->m_lsm6dslsensor->getGAxes(this->Gyroscope);

    res += sensors->m_lis2mdl->getMAxes(this->Magnetic);

    return res == 0;
}

bool SensorData::ReadTemperature(Sensors* sensors)
{
    this->Temperature = 0;
    int res           = sensors->m_hts221sensor->getTemperature(&this->Temperature);
    return res == 0;
}

bool SensorData::ReadHumidity(Sensors* sensors)
{
    this->Humidity = 0;
    int res        = sensors->m_hts221sensor->getHumidity(&this->Humidity);
    return res == 0;
}

bool SensorData::ReadPressure(Sensors* sensors)
{
    this->Pressure = -1.0;
    int res        = sensors->m_lps22hbsensor->getPressure(&this->Pressure);
    return res == 0;
}

bool SensorData::ReadAcceleration(Sensors* sensors)
{
    int res = sensors->m_lsm6dslsensor->getXAxes(this->Accelerator);
    return res == 0;
}

bool SensorData::ReadGyroscope(Sensors* sensors)
{
    int res = sensors->m_lsm6dslsensor->getGAxes(this->Gyroscope);
    return res == 0;
}

bool SensorData::ReadMagnetic(Sensors* sensors)
{
    int res = sensors->m_lis2mdl->getMAxes(this->Magnetic);
    return res == 0;
}
