# go-home

Automated home using Go (gokrazy.org, Go only OS for Raspberry PI)

Devices/Hardware:

- Netgear R6900
- Conbee II DECONZ
- Philips HUE lights
- Xiaomi Yee lights
- Xiaomi Aqara
  - Wall Switch
  - Button Switch
  - Power Switch
  - Sensors (Motion, Magnet)
- IKEA lights, switch and sensor
- Wemo switch
- Amazon Alexa
- Samsung TV
- Sony Bravia TV

Note:
  InfluxDB for tracking metrics and usage of all processes.
  Since all of this is written in Go it should be able to run anywhere, from a Raspberry PI to a Windows/Mac machine.

App Structure:

- Since every process is just running it's own logic want we need is kindof pub/sub server where every process
  can register itself to specific events that it is interested in.
  - NATS Server (Pub/Sub server where you can subscribe to channel(s))

- Following sub-processes:
  - Config              Ok, (A service that is the provider of configurations for all other services)
  - Presence            Ok, (Connects to Netgear Router to obtain list of devices present on the network)
  - Flux                Ok, (Calculates Color-Temperature and Brightness per day for Hue and Yee lights)
  - AQI                 Ok, (Air Quality Index)
  - Suncalc             Ok, (Computes sun-rise, sun-set etc..)
  - Weather             Ok, (Darksky)
  - Calendar            Ok, (Reads calendars from icloud and determines active events and makes sensors out of them)
  - Shout               Ok, (Has Slack as the back-end to send messages)
  - Wemo                Untested, (Wemo devices)
  - Conbee II DECONZ    WIP, (Philips HUE / IKEA / Xiaomi Aqara; lights, switches, sensors)
  - Apple HomeKit       WIP, (Apple Home Kit accessory emulator)
  - Yee                 Ok, (Xiaomi Yee lighting, turn on/off, change CT and BRI)
  - Sony Bravia Remote  Ok, (Turn on/off, HDMI Input, Volume for Sony Bravia TV(s))
  - Samsung TV Remote   Ok, (Turn on/off Samsung TV(s))
  
  - Automation; reacting to all events and executing automation rules (all written in Go)

Status:

I am currently not yet running this, however i have tested (nearly) all sub-processes and they all function. 
The only one not fully tested yet is 'automation'.
