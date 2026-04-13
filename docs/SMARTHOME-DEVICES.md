# Smart Home

- A HUB (on every floor) that can display information from all devices
  - Should have touch display so that we can control it easily.
- A local server that can receive data from all the sensors over WiFi
  - Sensors that can measure air quality, presence, bed presence, door/window status, etc.
  - Server will also send a UDP time-sync packet N seconds.

# Environment (Air Quality, Temperature, Humidity, Pressure, CO2, Luminosity)

These are 220V/USB-C powered.

- Temperature
- Humidity
- Pressure
- CO2
- Luminosity

## Air Quality Elite

This is for an area where we want to also measure PM, VOC, CO and NOx.

- SEN66, Sensirion Sensor (500 RMB, 70 USD)
  - Measures: Temperature, Humidity, PM (0.5, 1.0, 2.5, 10.0), CO2, VOC, NOx
  - I2C Interface
  - URL: https://sensirion.com/products/catalog/SEN66

## Air Quality Basic

- Temperature
- Humidity
- Pressure
- CO2
- Luminosity
- Presence

Would like to have a self designed 3D printed enclosure that can hold all the sensors as well as a screen to display the values. The enclosure should have good airflow to the sensors. Also when presence detects someone is nearby, the screen should light up to show the values, otherwise it should either be off or show the date, time and weather.

## Rain

- https://smartsolutions4home.com/ss4h-rg-rain-gauge/

This requires 3D printing of the design, seems very doable, very nice idea.

## Color Light Strips, Color Sync

We can make a setup that can sync the color of the light strips to the color of the ambient light or light coming from the TV/Monitor.

- One or more TCS34725; RGB Color Sensor with IR filter and White LED
  - I2C Interface
  - URL: https://www.adafruit.com/product/1334
- WS2812B or SK6812 RGB LED Strips
  - Addressable RGB LED Strips
  - URL: https://www.adafruit.com/product/1138

# Presence

Using mmWave Radar Sensors

The mmWave radar is always `on` when it detects nothing, but at the moment it detects something, it will slowly `back-off` at a certain rate, e.g. 10 seconds, 20 seconds, 1 minute, 5 minutes, etc.

## Bed Presence

- ESP32 C3 Mini
- light sensor (BH1750)
- mmWave sensor (RD03D 24GHz, 60°, 8m)
  - We will use two setups; ESP8266 + Relay + Light sensor setup, to turn on/off another full setup, ESP32 and RD03D mmWave

In the evening within a certain time window (configurable, e.g. 8pm until detection, time-out at 2am) the mmWave sensor is used to determine if the bed is occupied or not. We need to run a prototype with a RD03D mmWave sensor to see if we are able to detect 2 persons in the bed.

When the light sensor detects that the lights are turned off, we do the detection for 'bed presence' during a certain time duration (e.g. 10 minutes). After collecting info during this time duration, we determine if the bed is occupied or not, and how many persons are in the bed. Then during the whole night (e.g. until 7am) we just assume the bed is occupied.

## Room Presence

If we use a mmWave (low power) sensor that has 'trigger' output, we can use this to wake up the ESP32. The ESP32 can then send a message to the server that presence is detected, and then go back to deep sleep. If no presence is detected for a certain time (e.g. 5 minutes), the ESP32 can go back to deep sleep.

# Smart 

## Gas Meter

Using a camera and ESP32-S3 to periodically take a picture of the gas meter and send it over the network to the server. The server can then use OCR to read the gas meter and determine the gas usage.

## Water Meter

Using a camera and ESP32-S3 to periodically take a picture of the water meter and send it over the network to the server. The server can then use OCR to read the water meter and determine the water usage.

## Electricity Meter

We need to find the main electricity line and wrap it with a current sensor, e.g. SCT-013-000, to measure the electricity usage. The ESP32 can then send the data over the network to the server.

## Magic Cube

This is a cube that can be used to control automations in the house. 

## Button or Door/Window/Mailbox Sensor

URL: https://github.com/gadjet/Window-Door-sensor-Version-5/tree/main
     https://gadjetsblog.blogspot.com/2022/03/the-many-versions-of-wireles-door.html
Ordered at JLC PCB: https://www.jlcpcb.com

- ESP8266 12F
  - Power Latch, estimated (closed/open) 5.4 uA/3.9 uA
- Magnet Switch or Button
- 900 mAh 3.7 V LiPo Battery?

900 mAh / 10 uA = 90000 hours = 3750 days = 10.3 years (theoretical)

Battery Dimensions: 3cm x 2cm x 1cm = 6 cm³
