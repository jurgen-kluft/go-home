# Smart Home

- A HUB (on every floor) that can display information from all devices
  - Should have touch display so that we can control it easily.
- An ESPNOW Gateway that can receive data from all the sensors over ESPNOW and send it to a server
  - Sensors that can measure air quality, presence, bed presence, door/window status, etc.
  - Gateway will also send a time-sync packet every hour.

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

## Presence

- mmWave Radar Sensor

The mmWave radar is always `on` when it detects nothing, but at the moment it detects something, it will slowly `back-off` at a certain rate, e.g. 10 seconds, 20 seconds, 1 minute, 5 minutes, etc.

## Bed Presence

- ESP32 C3 Mini
- mmWave Radar Sensor
  - We will use two setups; ESP8266 + Relay + Light sensor setup, to turn on/off another full setup, ESP32 and RD03D mmWave

In the evening between a certain time window (configurable, e.g. 8pm until detection, time-out at 2am) the mmWave sensor is used to determine if the bed is occupied or not. We need to run a prototype with a RD03D mmWave sensor to see if we are able to detect 2 persons in the bed.

When the light sensor detects that the lights are turned off, we do the detection for 'bed presence' during a certain time duration (e.g. 10 minutes). After collecting info during this time duration, we determine if the bed is occupied or not, and how many persons are in the bed. Then during the whole night (e.g. until 7am) we just assume the bed is occupied.

## Battery Powered Prototypes

These are only for low frequency measurements, e.g. once every 5 minutes or once every 15 minutes:

# Low Frequency measuments, e.g. Air Quality

We can make a battery version of this, but we need to use a 'latch circuit' to power the sensors only when needed. The ESP32 can control a 'latch' to turn on the power to the sensors, and after reading the sensors, it can turn off the power again. This way, we can save a lot of power. 

- Temperature
- Humidity
- Pressure
- CO2
- Luminosity

- Firebeetle ESP32
  - Deep Sleep or even Hibernation Mode
- Latch Circuit
  - To power the sensors only when needed
- 1600 mAh LiPo Battery

1600 mAh / 100 uA = 16000 hours = 666 days = 1.8 years (theoretical)

# Presence Sensor

If we use a mmWave (low power) sensor that has 'trigger' output, we can use this to wake up the ESP32. The ESP32 can then send a message to the server that presence is detected, and then go back to deep sleep. If no presence is detected for a certain time (e.g. 5 minutes), the ESP32 can go back to deep sleep.

## Battery Powered Presence Sensor

- ESP32 C3 Xiao
  - 1600 mAh LiPo Battery
  - Deep Sleep Mode, estimated 10 uA in deep sleep
  - HiLink LD2410S mmWave Sensor
    - 0.1 mA average current consumption
    - Use the 'trigger' output to wake up the ESP32

1600 mAh / 110 uA = 14545 hours = 606 days = 1.6 years (theoretical)

# Door / Window / Mailbox Sensor

LoRa based: https://github.com/foorschtbar/LoRaProMini

URL: https://github.com/gadjet/Window-Door-sensor-Version-5/tree/main
     https://gadjetsblog.blogspot.com/2022/03/the-many-versions-of-wireles-door.html
Ordered at JLC PCB: https://www.jlcpcb.com

- ESP8266 12F
  - Power Latch, estimated (closed/open) 5.4 uA/3.9 uA
- Magnet Switch or Button
- 400 mAh 3.7 V LiPo Battery?

400 mAh / 10 uA = 40000 hours = 1666 days = 4.5 years (theoretical)
Battery Dimensions: 3cm x 2cm x 1cm = 6 cm³
