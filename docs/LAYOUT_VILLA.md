# Zigbee, WiFi and BLE

- ESP32 WROOM 32E
- Temperature
- Humidity
- Pressure
- CO2
- Luminosity

## Floor

Every floor has WiFi cover, so we are targeting to mostly use WiFi and we can experiment with EspNOW.

## Sensors

- Door/Window open/close sensors, these are UDP/EspNOW based and connect to the WiFi of the floor they are on.
- Motion sensors, these are ESP32 boards with mmWave sensors and connect to WiFi of the floor they are on.
- Temperature, Humidity, Pressure, CO2 and Luminosity sensors, these are ESP32 boards with the respective sensors connected to them. 
  They connect to the WiFi network and send their data to sensor server over TCP. 

## Lights

- HUE lights, these connect to the HUE hub of the floor they are on.

## Wall Panels / Switches

- ESP32 based wall panels / switches, these connect to the Sensor Server.

## Automation

If we are able to collect all the data from the sensors, we can use it to automate certain actions in the house. Since we can also setup a DB of all the lights in the house using HUE, we can use the data from the sensors to automate the lights as well as other devices in the house, for example TV's.

Since we mostly work in Golang, we can write our automations in Golang.

