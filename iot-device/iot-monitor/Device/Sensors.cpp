#include "HTS221Sensor.h"
#include "LIS2MDLSensor.h"
#include "LPS22HBSensor.h"
#include "LSM6DSLSensor.h"
#include "OledDisplay.h"
#include "Sensors.h"

void Sensors::Init()
{
    Screen.print(2, "Initializing...");
    m_i2c = new DevI2C(D14, D15);

    Screen.print(3, " > Magnetometer");
    m_lis2mdl = new LIS2MDLSensor(*m_i2c);
    m_lis2mdl->init(NULL);

    Screen.print(3, " > HTS221Sensor");
    m_hts221sensor = new HTS221Sensor(*m_i2c);
    m_hts221sensor->init(NULL);
    m_hts221sensor->enable();
    m_hts221sensor->reset();

    Screen.print(3, " > LSM6DSLSensor");
    m_lsm6dslsensor = new LSM6DSLSensor(*m_i2c, D4, D5);
    m_lsm6dslsensor->init(NULL);
    m_lsm6dslsensor->enableAccelerator();
    m_lsm6dslsensor->enableGyroscope();

    Screen.print(3, " > LPS22HBSensor");
    m_lps22hbsensor = new LPS22HBSensor(*m_i2c);
    m_lps22hbsensor->init(NULL);
}

bool SensorData::ReadAll()
{
    int res = 0;

    this->Temperature = 0;
    this->Humidity = 0;
    this->Pressure = -1.0;
    
    res += m_hts221sensor->getTemperature(&this->Temperature);
    res += m_hts221sensor->getHumidity(&this->Humidity);

    res += m_lps22hbsensor->getPressure(&this->Pressure);

    res += m_lsm6dslsensor->getXAxes(this->Accelerator);
    res += m_lsm6dslsensor->getGAxes(this->Gyroscope);

    res += m_lis2mdl->getMAxes(this->Magnetic);

    return res == 0;
}

bool SensorData::ReadTemperature(SensorData* sd)
{
    this->Temperature = 0;
    int res           = m_hts221sensor->getTemperature(&this->Temperature);
    return res == 0;
}

bool SensorData::ReadHumidity()
{
    this->Humidity = 0;
    int res        = m_hts221sensor->getHumidity(&this->Humidity);
    return res == 0;
}

bool SensorData::ReadPressure()
{
    this->Pressure = -1.0;
    int res        = m_lps22hbsensor->getPressure(&this->Pressure);
    return res == 0;
}

bool SensorData::ReadAcceleration()
{
    int res = m_lsm6dslsensor->getXAxes(this->Accelerator);
    return res == 0;
}

bool SensorData::ReadGyroscope()
{
    int res = m_lsm6dslsensor->getGAxes(this->Gyroscope);
    return res == 0;
}

bool SensorData::ReadMagnetic()
{
    int res = m_lis2mdl->getMAxes(this->Magnetic);
    return res == 0;
}
