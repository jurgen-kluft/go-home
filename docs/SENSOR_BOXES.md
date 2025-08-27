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

## Procedure

Read the sensor(s) at a certain interval (e.g. 10 times per second), and send a network message at another interval (e.g. once every  second). A message can contain multiple blocks of sensor data, each block with a timestamp.

