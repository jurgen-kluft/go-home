# Sensors

Here a, incomplete, list of sensors that are very practical:

- Light (lux)
  - BH1750
- Temperature (°C)
  - BME280
  - BMP280
  - SCD41
- Humidity (%)
  - BME280
  - SCD41
- Pressure (hPa)
  - BME280
  - BMP280
- Magnetic field (uT)
  - HMC5883L
- Acceleration (g)
- Gyroscope (°/s)
- Vibration
  - 801s
- Sound (dB)
- CO2 (ppm)
  - SCD41
- PM2.5 (ppm)
- PM10 (ppm)
- VOC (ppb)
- Motion, PIR (detection)
  - HC-SR501
- Presence, mmWave (detection)
  - RD-03D

## Modules

- BH1750 light sensor 
  - Voltage: 3.3V - 5V
  - I2C address: 0x23
  - URL: https://github.com/claws/BH1750
- BME/BMP280 temperature, humidity, and pressure sensor
  - Voltage: 3.3V
  - I2C address: 0x76
  - URL: https://github.com/finitespace/BME280
- SCD41 CO2 sensor
    - Voltage: 3.3V - 5V, most important is a stable power supply
    - I2C address: 0x68
    - URL: https://github.com/Sensirion/arduino-i2c-scd4x
- HC-SR501 PIR motion sensor
    - Voltage: 5V
    - URL: https://github.com/helenhoffman/ESP32_MotionSensor
- RD-03D mmWave sensor 
    - Voltage: 5V, and the power ripple is required to be controlled within 100mV
    - I2C address: 0x10
    - URL: https://github.com/Gjorgjevikj/Ai-Thinker-RD-03
