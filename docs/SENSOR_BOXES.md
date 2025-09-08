# Air Quality

- ESP32 WROOM 32E
- Temperature
- Humidity
- Pressure
- CO2
- Luminosity

## Presence

- ESP32 C3 Mini
- mmWave Radar Sensor

The mmWave radar is always `on` when it detects nothing, but at the moment it detects something, it will slowly `back-off` at a certain rate, e.g. 10 seconds, 20 seconds, 1 minute, 5 minutes, etc.

## Procedure

Read the sensor(s) at a certain interval (e.g. 10 times per second), and send a network message at another interval (e.g. once every  second). A message can contain multiple blocks of sensor data, each block with a timestamp.

