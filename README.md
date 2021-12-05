# go-home

Automated home using Go.

## Raspberry PI

``gokrazy.org``, a Go only OS for Raspberry PI.

## Status

I am currently running this in the weekend, still testing some sub-processes and now focussing on Apple Home-Kit interaction so that we have an actual UI.

## App Structure

Note:
  InfluxDB for tracking metrics and usage of some processes.
  Since all of this is written in Go it should be able to run anywhere, from a Raspberry PI, Arduino to a Windows/Mac machine.

- Since every process is just running it's own logic want we need is a pub/sub server where every process
  can register itself to specific events that it is interested in.
  -> NATS Server (Pub/Sub server where you can subscribe to channel(s))
  -> Researching the possibility of using zeroconf 'Service Discovery - mDNS'

- Following sub-processes:
  - Apple HomeKit       WIP, (Apple Home Kit accessory emulator, our **UI** solution)
  - Conbee II DECONZ    Ok, (Philips HUE / IKEA / Xiaomi Aqara; lights, switches, sensors)
  - Wemo                Ok, (Wemo wifi powerplug)
  - Config              Ok, (A service that is the provider of configurations for all other services)
  - Presence            Ok, (Connects to Netgear Router to obtain list of devices present on the network)
  - Flux                Ok, (Calculates Color-Temperature and Brightness per day for Hue and Yee lights)
  - AQI                 Ok, (Air Quality Index)
  - Suncalc             Ok, (Computes sun-rise, sun-set etc..)
  - Weather             Ok, (Darksky)
  - Calendar            Ok, (Reads calendars from icloud and determines active events and makes sensors out of them)
  - Shout               Ok, (Has Slack as the back-end to send messages, but requires a VPN :-(, looking for something else)
  - Yee                 Ok, (Xiaomi Yee lighting, turn on/off, change CT and BRI)
  - Sony Bravia Remote  Ok, (Turn on/off, HDMI Input, Volume for Sony Bravia TV(s))
  - Samsung TV Remote   Ok, (Turn on/off Samsung TV(s))
  
## Automation Logic
  
Automation, reacting to all events and executing automation rules, all written in Go.

## Devices / Hardware

- Netgear R6900
- Conbee II DECONZ
- Philips HUE lights
- Xiaomi Yee lights
- Xiaomi Aqara
  - Button Switch
  - Power Switch
  - Sensors (Motion, Magnet)
- IKEA lights, switch and sensor
- Wemo switch
- Amazon Alexa
- Samsung TV
- Sony Bravia TV
