# ESP32, YD Wroom 32E

- GPIO 21 (SDA) and GPIO 22 (SCL)
- Has a Led-Strip
  - west build -b yd_esp32/esp32/procpu samples/drivers/led/led_strip

- 5V pin is directly connected to the USB-C connector, so it can be used to power other devices.

## Details

I2C ports: 

- SDA = 21
- SCL = 22

Working sensors:

- BME280
- SDC41
- BH1750


