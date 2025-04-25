# IoT Hardware

## HomeKit / Airplay capable TV's

- LG TV; https://www.lg.com/cn/tvs-soundbars/lg-oled65c5pca

## Smart Lights

Wiz (Philips) smart lights are connecting to WiFi and can easily be controlled with HomeKit by exposing them as HomeKit accessories by using Golang.
We can group many bulbl under one `light` and expose it as one light to HomeKit, also these lights are connected through WiFi should they should react a lot quicker than Zigbee lights.

Aqara light bulbs (white) are Zigbee based and can be controlled using the Aqara Hub and can thus end up in HomeKit.

Philips Hue (color or white) can be used in the bedrooms to avoid the need of WiFi. They can be controlled using the Hue Bridge. The Hue Bridge can be connected to LAN and then it can be used to control the lights. Only Hue can be controlled programmatically using Golang. 
Overall I think we should avoid Hue lights since we need an extra Hub (per floor?).

## Smart Plugs

Can measure power consumption and control devices remotely. For example it can detect a wash machine, dryer, dish washer, coffee machine running and ending when the power consumption drops to a certain level.

## Room Presense Sensors

I have bought 3 `LinknLink eMotion Pro` sensors on `Amazon.nl`, they are Wifi IP based and should be connected to a MQTT broker. I will also order 2 presence sensors from `SmartHomeShop`, these also include CO2, VOC, PM, Lux, NOx, and are based on ESPHome, we can write a process that will read the data from the sensors and send it to the MQTT broker.
They can be exposed to HomeKit (if necessary) through the use of Golang, however we can also just expose switches that indicate presence in a room/area.

- Living Room (1st floor)
- Living Room (2nd floor)
- Bath Room (2nd floor)
- Study Room (2nd floor)

## Bed Presence Sensors (DYI, WIP)

- SEN-09674 FSR (can be 600mm long, with two of them you can detect presence on both sides of the bed)
  - 2x 10K Ohm resistors
  - 2 x 72.5 RMB     = 145 RMB
- ESP32              = 22 RMB
- USB-C power supply = 20 RMB
- Total              = 187 RMB 

- https://github.com/eoncire/HA_bed_presence
- https://www.homeautomationguy.io/blog/making-my-own-bed-sensor

Door Contact sensors can be repurposed as pressure sensors. They can then be used in a chair, sofa or a bed to detect if someone is sitting or laying down. This can be used to trigger automations like turning on some lights when getting out of bed in the middle of the night. Or for the living room sofa to pause/resume a movie, resume when someone is sitting down and pause when they get up.

## Air Quality

- ESP32 WROOM 32E (https://www.espressif.com/sites/default/files/documentation/esp32-wroom-32e_esp32-wroom-32ue_datasheet_en.pdf)

### Luminosity

Light sensors can be used to detect if it is dark outside and turn on the lights in the house. They can also be used to detect if it is bright outside and turn off the lights in the house.

- BH1750 (16 bit I2C light sensor, 1-65535 lux)

### Temperature, Pressure and Humidity

The BME280 is a humidity sensor measuring relative humidity, barometric pressure and ambient temperature.

- BME280 (Bosch, https://www.bosch-sensortec.com/products/environmental-sensors/humidity-sensors-bme280/)

### Carbon Dioxide (CO2)

- SENSIRION SDC41 CO2 Sensor (https://www.sensirion.com/en/environmental-sensors/air-quality/sdc41-co2-sensor/)

## SmartThings Station

Can serve as an iPhone wireless charger but also can detect if it is charging, so some automation can be triggered when this event happens. For example, when the iPhone is charging we identify it with going to sleep, so turn off the lights in the bedroom. When waking up in the middle of the night, we can turn on some night lights in the bedroom/bathroom.

## DIY Presence Detection

### ESP32 with DFRobot C4001

- Wi-Fi 2.4GHz, Bluetooth 4.2
- 5V1A, USB
- ESPHome
- DFRobot 24GHz, 100°, 25m

Cost: 

- DFRobot C4001 = 180 RMB
- ESP32 = 25 RMB
- USB cable = 10 RMB
- Power supply = 20 RMB
- Total = 235 RMB

### ESP32 with LD2410C

Can only track one person at a time, but it is cheaper than the DFRobot C4001. It can be used to detect if the room is occupied or not.

- Wi-Fi 2.4GHz, Bluetooth 4.2
- 5V1A, USB
- ESPHome
- LD2410B/C 24GHz, 60°, 6m

Cost:

- LD2410B/C = 50 RMB
- ESP32 = 25 RMB
- USB cable = 10 RMB
- Power supply = 20 RMB
- Total = 105 RMB