# Time Sync

On our network, we are running a local NTP server, so all devices can sync their time to this server. This is important for logging and scheduling tasks. We can simply use UDP to send NTP requests to the server and get the current time.

Upon setup we have to synchronize the tick, and we have to have a global reference time (e.g. epoch time) to calculate the current time.
64 bit, 1us accuracy, results in a rollover every 584 years.
To avoid time drift, we can sync the time every day or so.

# Air Quality

- Temperature
- Humidity
- Pressure
- CO2
- Luminosity

## Presence

- mmWave Radar Sensor

The mmWave radar is always `on` when it detects nothing, but at the moment it detects something, it will slowly `back-off` at a certain rate, e.g. 10 seconds, 20 seconds, 1 minute, 5 minutes, etc.

## Magnetic Sensor

A3144E Hall Effect Sensor, the Vcc pin should be connected to 3.3V, the GND pin to GND and the OUT pin to a GPIO pin on the ESP32.
The signal pin should be pulled up to 3.3V using a 10K resistor. When the magnet is close to the sensor, the output will be LOW, otherwise it will be HIGH.

Pinout: Vcc | GND | Signal

## Bed Presence

- ESP32 C3 Mini
- mmWave Radar Sensor
- Light Sensor

In the evening between a certain time window (e.g. 8pm to 11pm) the mmWave sensor is used to determine if the bed is occupied or not. We need to run the prototype with a RD03D mmWave sensor to see if we are able to detect 2 persons in the bed.
When the light sensor detects that the light are turned off, we do the detection for 'bed presence' during a certain time duration (e.g. 10 minutes). After collecting info during this time duration, we determine if the bed is occupied or not, and how many persons are in the bed. Then during the whole night (e.g. 11pm to 7am) we just assume the bed is occupied.

## Procedure

Read the sensor(s) at a certain interval (e.g. 10 times per second), and send a network message at another interval (e.g. once every  second). A message can contain multiple blocks of sensor data, each block with a timestamp.

## Battery Powered Prototypes

These are only for low frequency measurements, e.g. once every 5 minutes or once every 15 minutes:

- Temperature
- Humidity
- Pressure
- CO2
- Luminosity

- Firebeetle ESP32
  - Deep Sleep or even Hibernation Mode
- 1600 mAh LiPo Battery

1600 mAh / 100 uA = 16000 hours = 666 days = 1.8 years (theoretical)