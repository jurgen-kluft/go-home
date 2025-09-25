# Zigbee, WiFi and BLE

- ESP32 WROOM 32E
- Temperature
- Humidity
- Pressure
- CO2
- Luminosity

## Floor

Every floor has a Zigbee2MQTT coordinator that directly connects to the MQTT broker, either over ethernet or WiFi.

## Sensors

- Door/Window open/close sensors, these are Zigbee based and connect to the Zigbee2MQTT coordinator of the floor they are on.
  - For example, the Aqara door/window sensor.
- Motion sensors, these are Zigbee based and connect to the Zigbee2MQTT coordinator of the floor they are on.
  - For example, the Aqara motion sensor.
- Temperature, Humidity, Pressure, CO2 and Luminosity sensors, these are ESP32 boards with the respective sensors connected to them. 
  They connect to the WiFi network and send their data to sensor server over TCP. 

## Lights

- Zigbee based lights, these connect to the Zigbee2MQTT coordinator of the floor they are on.

## Wall Panels / Switches

- Zigbee based wall panels / switches, these connect to the Zigbee2MQTT coordinator of the floor they are on.
  For example, the Aqara wall panel.


## MQTT Broker


## MQTT Client for the Sensor Server

We have a custom MQTT client that connects to the MQTT broker and listens for messages from particular sensors.
This client is written in Golang and will send the data to the sensor server over TCP.



## Automation

If we are able to collect all the data from the sensors, we can use it to automate certain actions in the house.
Since we can also setup a DB of all the lights in the house using Zigbee2MQTT, we can use the data from the sensors 
to automate the lights as well as other devices in the house, for example TV's.

Since we mostly work in Golang, we can write our automations in Golang.

