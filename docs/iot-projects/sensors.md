# Sensors

- Light (lux)
- Temperature (°C)
- Humidity (%)
- Pressure (hPa)
- Sound (dB)
- Vibration
- CO2 (ppm)
- PM2.5 (ppm)
- PM10 (ppm)
- VOC (ppb)
- Motion, PIR (detection)
- Presence, mmWave (detection)

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
- RD-03D mmWave sensor 
    - Voltage: 5V, and the power ripple is required to be controlled within 100mV
    - I2C address: 0x10
    - URL: https://github.com/Gjorgjevikj/Ai-Thinker-RD-03
