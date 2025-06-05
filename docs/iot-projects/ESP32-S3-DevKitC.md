# ESP32, S3 Devkit C

- NOTE: 5V pin cannot be used to power other devices, it seems to be a IN only pin.
  Seems there is a IN-OUT solder pad, if you solder it up then you can use the 
  5V pin to power other devices.

- NOTE: SDA and SCL pins are not labeled on the board, but they are the same as the ESP32-S3-DevKitC-1:
  - SDA is GPIO 8
  - SCL is GPIO 9

## Details

I2C ports: 

- SDA = 8
- SCL = 9

Working sensors:

- BME280 
- SDC41
