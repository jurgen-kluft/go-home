# NATS

## NATS Messaging

Listen to config/request/ messages, when a subscriber registers we can send him the configuration.
Also when we detect that the configuration on disk has changed, we can hot-load it and send
it to the associated channel. (This part of the config service and is working)

## Azure IoT Devkit - MXCHIP

Record and transmit, high frequency (1000 Hz?)

* accelerometer
* magnetometer
* gyroscope
* audio (maybe pre-processed?)

Record and transmit, low frequency (10 Hz)

* temperature
* humidity
* pressure
