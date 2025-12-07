# Bedroom motion sensor

In the bedroom, we want to detect motion during the day but not during the night to avoid RF polution disturbing sleep.
Also it is not needed once the sensor knows the state of the bedroom.

- Uses ESP8266
- Uses AI Thinker RD03 mmWave Sensor
- Used Relay to turn-off ESP8266 which has the Motion Sensor to remove interference during the night

## 3D Case

- USB-C socket
- Main ESP8266
  - APDS9960 sensor (color, gesture, proximity)
- Relay board
  - Motion ESP8266
  - RD03D Sensor (sensor stick)
  - OLED display



