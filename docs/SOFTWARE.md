# Software

- Go; https://golang.org/, 1.22.12

## ESP32, UDP messages

I would like the ESP32 devices to simply send small binary UDP messages to a custom (Golang?) server, and the server should be able to receive, process, and redirect the messages to the correct MQTT topic.

We just need to define a basic message format:
- `size` (int16, size of the message(including this field))
- `device_location` (int16, 0 = unknown, 1 = bedroom, 2 = living room, 3 = kitchen, 4 = bathroom, 5 = hallway, 6 = balcony)
- `device_id` (int16, )
- `sensor_count` (int16)
  - [
    - `sensor_type` (int16)
    - `sensor_state` (int8, 0 = off, 1 = on, -1 = error)
    - `sensor_value` (float32[]/int32[])
  - ]

Device Locations:

- 0 = Unknown
- 1 = Bedroom
- 2 = Living Room
- 3 = Kitchen
- 4 = Bathroom
- 5 = Hallway
- 6 = Sophia Room
- 7 = Jennifer Room

Device IDs:

- 0 = Unknown
- 1 = Bed Presence
- 2 = Air Quality
- 3 = Air Quality + Room Presence
- 4 = Bed Presence + Air Quality

Sensor Types:

- 0 = None
- 1 = Temperature
- 2 = Humidity
- 3 = Pressure
- 4 = Light
- 5 = CO2
- 6 = Presence
- 7 = Motion
- 8 = Target (X/Y/Z) (e.g., for the RD-03D sensor)

## Message Broker

Custom message broker to receive UDP messages from ESP32 devices and redirect them to the correct MQTT topic.

## MQTT Broker

- Mochi MQTT; https://github.com/mochi-mqtt/server (golang)
- MQTT Broker; Mosquitto

## ESPHome

From golang we can receive data from ESPHome devices, and we can send commands to them.

- https://github.com/mycontroller-org/esphome_api

### Custom Devices

- https://emanuelduss.ch/posts/co2-measurement/    

## Philips Wiz Lights

- https://github.com/squarejaw/wiz
- https://github.com/achetronic/wizgo

## Go

- Actor Model; https://github.com/vladopajic/go-actor

