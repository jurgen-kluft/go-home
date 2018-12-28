# Emitter.IO

Been looking into emitter.io, a high-performance pub/sub server that seems a lot more suitable
to our situation. It also uses MQTT as the message protocol which makes it interesting from
a IoT point of view.

An emitter client can subscribe to channels, so generally for all of our processes they should
subscribe to their config channel, state-request channel.

## Config Emitter Client

Listen to presence messages, when a subscriber registers we can send him the configuration.
Also when we detect that the configuration on disk has changed, we can hot-load it and send
it to the associated channel.

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
