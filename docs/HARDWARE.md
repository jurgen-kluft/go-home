# IoT Hardware

## HomeKit / Airplay capable TV's

- LG TV; https://www.lg.com/cn/tvs-soundbars/lg-oled65c5pca

## MAC Mini

Mac Mini M4 to run:

- Sensor Server (golang)
- HomeKit Bridge?
- MQTT Broker
- Image Recognition (for presence detection)
- Other services
  - Immich (photo server)

## NAS

- [X] Ugreen NAS DH4800 Plus (https://www.youtube.com/watch?v=hEu6LTKbqcA)

## ESP8266 boards

Very good for DIY trigger based sensors, e.g. door/window sensors, water leak sensors, vibration sensors, etc.

Also low frequency (wakeup) based sensors can be built using ESP8266 modules, like: 

- temperature
- humidity
- pressure
- luminosity

## ESP32 boards

Good board brands:

- DFRobot, e.g. FireBeetle
- LOLIN, e.g. D1 Mini
- Waveshare
- Seeed Studio, Xiao

## CYD (cheap yellow display)

These can work in combination with OpenHAB :)

- There are many different versions, I have bought the 2.4 inch version (39 RMB)
- There is also a 2.8 and 3.5 inch version (43, 55 RMB)

## Lab Power Supply

Checked out many, and found this one to be of good quality and price.

- UNI-T UDP3305S

## Smart Wall Panels

We would like to have smart wall panels in every room to control lights, plugs, scenes and other automations. The wall panels should be flush mounted in the wall and have a nice touch screen to control everything. The wall panels should be connected to the network through Zigbee.

- GeekOpen, these look like normal light switches, but inside they have a ESP8266 that sends out JSON packets to a configurable TCP server. You can send JSON
  commands to the switch to change the state of the switch. They are powered using 220V AC.

  - https://www.smart-bird.cn/smart.html

## Smart Plugs

Can measure power consumption and control devices remotely. For example it can detect a wash machine, dryer, dish washer, coffee machine running, mobile phone charging.

  - https://www.smart-bird.cn/smart.html

## Smart Lights

Wiz (Philips) smart lights are connecting to WiFi and can easily be controlled with HomeKit by exposing them as HomeKit accessories by using Golang.
We can group many bulbs under one `light` and expose it as one light to HomeKit, also these lights are connected through WiFi so they should react a lot quicker than Zigbee lights.

- https://github.com/squarejaw/wiz
- https://github.com/achetronic/wizgo

Aqara light bulbs (white) are Zigbee based and can be controlled using the Aqara Hub and can thus end up in HomeKit.

Philips Hue (color or white) can be used in the bedrooms to avoid the need of WiFi. They can be controlled using the Hue Bridge. The Hue Bridge can be connected to LAN and then it can be used to control the lights. Only Hue can be controlled programmatically using Golang. Overall I think we should avoid Hue lights since we need an extra Hub (per floor?).

## Room Presense Sensors

We will also build them ourselves using ESP32 devices on WiFi can easily be exposed to HomeKit (if necessary) through the use of Golang, however we can also just expose switches that indicate presence in a room/area using `https://github.com/brutella/hap`.

## Bed Presence Sensors (DYI, WIP)

One other option

- RD03D 24GHz, 60°, 8m
- ESP32              = 22 RMB
- USB-C power supply = 20 RMB
- Total              = 187 RMB 

- https://github.com/eoncire/HA_bed_presence
- https://www.homeautomationguy.io/blog/making-my-own-bed-sensor

### Luminosity

Light sensors can be used to detect if it is dark outside and turn on the lights in the house. 
They can also be used to detect if it is bright outside and turn off the lights in the house.

### Temperature, Pressure and Humidity

The TMP117 is a high accuracy temperature sensor with an accuracy of ±0.1°C over the -20°C to +50°C temperature range.
Price is a bit high, but it is very accurate, around 90 RMB on TaoBao.

- TMP117 (Texas Instruments, https://www.ti.com/product/TMP117)

The BME280 is a humidity sensor measuring relative humidity, barometric pressure and ambient temperature.

- BME280 (Bosch, https://www.bosch-sensortec.com/products/environmental-sensors/humidity-sensors-bme280/)

### Carbon Dioxide (CO2)

- SENSIRION SDC41 CO2 Sensor (https://www.sensirion.com/en/environmental-sensors/air-quality/sdc41-co2-sensor/)

## DIY Presence Detection

### ESP32 with DFRobot C4001

- Wi-Fi 2.4GHz, Bluetooth 4.2
- 5V1A, USB
- ESPHome or Custom
- DFRobot 24GHz, 100°, 25m

Cost: 

- DFRobot C4001 = 100 RMB
- ESP32 = 30 RMB
- USB cable = 10 RMB
- Power supply = 20 RMB
- Total = 160 RMB

### ESP32 with AiThinker RD03D

Can track multiple targets at a time, it can be used to detect if the room is occupied or not.

- Wi-Fi 2.4GHz, Bluetooth 4.2
- 5V1A, USB
- RD03D 24GHz, 60°, 11m

Cost:

- RD03D = 30 RMB
- ESP32 = 30 RMB
- USB cable = 10 RMB
- Power supply = 20 RMB
- Total = 90 RMB

