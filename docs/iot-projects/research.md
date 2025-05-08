# Sensors

- BH1750 light sensor 
  - I2C address: 0x23
  - URL: https://github.com/claws/BH1750
- BME280 temperature, humidity, and pressure sensor
  - I2C address: 0x76
  - URL: https://github.com/finitespace/BME280
- SDC41 CO2 sensor
    - I2C address: 0x68
    - URL: https://github.com/Sensirion/arduino-i2c-scd4x
- RD-03D mmWave sensor 
    - I2C address: 0x10
    - URL: 
      - https://github.com/javier-fg/arduino_rd-03d
      - https://github.com/MauricioOrtega10/Rd-03
      - https://github.com/bertrik/aithinker-rd03
      - https://github.com/Gjorgjevikj/Ai-Thinker-RD-03

# Breadboard prototype, soldering

A breadboard prototype board that can be soldered, should make it a lot easier to connect wiring together.

# 10cm x 10cm x 3.2cm case

Sufficient to build all the sensor units:

- Air quality
- Air quality + Room Presence
- Bed presence

# Soldering 4P (2.54 mm) Connector

Every sensor has a 4P connector, so we can easily connect/disconnect.

# ESP32 S3 N16R8

I2C ports: 

- SDA = 8
- SCL = 9

Working sensors:

- BME280 
- SDC41

# VCC-GND STUDIO, ESP32 WROOM 32E

I2C ports: 

- SDA = 21
- SCL = 22

Working sensors:

- BME280
- SDC41
- BH1750


