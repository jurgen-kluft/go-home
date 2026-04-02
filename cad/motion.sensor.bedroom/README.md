# Bedroom motion sensor

In the bedroom, we want to detect motion during the day but not during the night to avoid RF polution disturbing sleep.
Also it is not needed once the sensor knows the state of the bedroom.

- ESP32
  - Uses Relay to turn-off:
    - ESP8266 (+ motion sensor)
    - AI Thinker RD03 mmWave Sensor

## 3D Case

- USB-C socket
- Main ESP32
  - APDS9960 sensor (color, gesture, proximity)
  - Relay board (turn On/Off power to the ESP8266 and RD03D sensor)
    - ESP8266
    - RD03D Sensor (sensor stick)

