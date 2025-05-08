I would like the ESP32 devices to simply send small binary UDP messages to a custom server.
The server, written in Golang, should be able to receive, process, and redirect the messages 
to a MQTT broker and even Grafana/InfluxDB.
For presence it will also be responsible to compute the presence of areas in a room. These
areas are configured in the server and are not part of the message.

We just need to define a basic message format:

- `size` (size of the message in bytes (including this field))
- `device_location` 
- `device_label` 
- `sensor_count` = [`type`, {optional: `channel index`}, {optional: `state`}, `value`]

Field Type (only 4 types):
- int8 = 1 byte
- int16 = 2 bytes
- int32 = 4 bytes
- float32 = 4 bytes

Device Locations:

- 0 = Unknown
- 1 = Bedroom
- 2 = Living Room
- 3 = Kitchen
- 4 = Bathroom
- 5 = Hallway
- 6 = Sophia Room
- 7 = Jennifer Room

Device Labels:

- 0 = Unknown
- 1 = Bed Presence
- 2 = Air Quality
- 3 = Air Quality + Room Presence
- 4 = Bed Presence + Air Quality

Sensor State (optional):

- default = 1
- 0 = Off
- 1 = On
- -1 = Error

Sensor Channel Index (optional):

- default = 0
- 0 - 255

Sensor Types:

- 0 = None
- 1 = Temperature
- 2 = Humidity
- 3 = Pressure
- 4 = Light
- 5 = CO2
- 6 = Presence
- 7 = Motion
- 8 = Target (channel index indicates X, Y, Z axis)
