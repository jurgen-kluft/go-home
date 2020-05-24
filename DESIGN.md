# NATS

## NATS Messaging

Listen to config/request/ messages, when a subscriber registers we can send him the configuration.
Also when we detect that the configuration on disk has changed, we can hot-load it and send
it to the associated channel. (This part of the config service and is working)

## ZeroConf and MsgBus/

We could also make the whole service infrastructure zero-conf. Using Bonjour (mDNS / DNS-SD Service Discovery)
to find running services. Then when we find services that we need to subscribe/publish on we connect to them
using MsgBus. In this way we do not need NATS and we can also run a service anywhere on the LAN without having
to configure anything. We can also launch a service that can announce the IP:Port of InfluxDB. Makes it a lot
easier to locally (iMac) test a service.

A micro-service thus has the following:

- SD Announce
  - Keep announcing continuesly until we have the expected subscribers and then switch to a
    low frequency announcement (10 seconds every 5 minutes)
- SD Resolve; until we have found all services that we are interested in
- PubSub server (For subscribers to connect to)
- PubSub client (We connecting to the PubSub server of the services)

Also Bonjour can find apple devices like iPhone, Apple TV etc.. so this can also serve as a presence detection.

## Azure IoT Devkit - MXCHIP

Record and transmit, high frequency (1000 Hz?)

- accelerometer
- magnetometer
- gyroscope
- audio (maybe pre-processed?)

Record and transmit, low frequency (10 Hz)

- temperature
- humidity
- pressure
