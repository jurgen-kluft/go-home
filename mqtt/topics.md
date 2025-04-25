# MQTT Topics

Naming convention:

- `device-type/location/room/zone/category/device-name`

e.g.

- `sensor/1stfloor/kitchen/ceiling/temperature/sensor-01`
- `sensor/1stfloor/kitchen/ceiling/humidity/sensor-01`
- `sensor/1stfloor/livingroom/main/motion/sensor-02`
- `tv/1stfloor/livingroom/main/power/lg-tv-01`

The payload is always JSON and should contain only fields that require updates.
- `sensor/1stfloor/kitchen/ceiling/temperature/sensor-01`
    - `{"temperature": 22.5}`
- `sensor/1stfloor/kitchen/ceiling/humidity/sensor-01`
    - `{"humidity": 45}`
- `sensor/1stfloor/livingroom/tv/motion/sensor-02`
    - `{"motion": true}`
- `tv/1stfloor/livingroom/main/power/lg-tv-01`
    - `{"power": "on"}`
