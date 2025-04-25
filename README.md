# go-home

Automated home using Go.

## Status

Total overhaul of the project, all code is being scratched and redone in Go.

## App Structure

Note:
  InfluxDB for tracking metrics and usage of some processes.
  Since all of this is written in Go it should be able to run anywhere, from a Raspberry PI, Arduino to a Windows/Mac machine.

- Since every process is just running it's own logic want we need is a pub/sub server where every process
  can register itself to specific events that it is interested in.
  -> MQTT Server (Pub/Sub server where you can subscribe to channel(s))
     This is also useful for any third-party device to integrate into the whole system.

- Following sub-processes:
  - Apple HomeKit       WIP, (Apple Home Kit accessory emulator, our **UI** solution)
  - Wemo                Ok, (Wemo WiFi powerplug)
  - Config              Ok, (A service that is the provider of configurations for all other services)
  - Presence            Ok, (Connects to Router to obtain list of devices present on the network)
  - Flux                Ok, (Calculates Color-Temperature and Brightness during the day)
  - AQI                 Ok, (Air Quality Index)
  - Suncalc             Ok, (Computes sun-rise, sun-set etc..)
  - Season              Ok, (Computes the season (Spring, Summer, Autumn, Winter) based on the date)
  - Weather             Ok, (Darksky)
  - Calendar            Ok, (Reads calendars from icloud and determines active events and makes sensors out of them)
  - Shout               Ok, (Push notifications to HomePods?)
  - Yee                 Ok, (Xiaomi Yee lighting, turn on/off, change CT and BRI)
  - Wiz                 Ok, (Xiaomi Yee lighting, turn on/off, change CT and BRI)
  - Sony Bravia TV      Ok, (Turn on/off, HDMI Input, Volume for Sony Bravia TV(s))
  - Samsung TV Remote   Ok, (Turn on/off Samsung TV(s))

TODO:
  - ESPHome; For custom mmWave/Air-Quality sensors (https://github.com/pteich/esphomekit)
  - DD-WRT; For some form of presence (https://github.com/awilliams/wifi-presence)
  - MQTT; https://github.com/antlinker/libmqtt
  
## Automation Logic
  
Automation, reacting to all events and executing automation rules, all written in Go.

## Devices / Hardware

- Home 1:
  - Aqara M2 Hub (Zigbee Sensors, Lights, Switches, Plugs)
  - Aqara E1 Hub (Zigbee Repeater)
  - Apple Mac Mini M4, running Go
  - Apple TV 4K (this is the HomeKit controller since it can be Ethernet connected)
  - Apple HomePod
  - 2 TP-Link Deco routers
  - Hue, Yee, Ikea, Wiz and Aqara lighting
  - Aqara Switches
  - Sony Bravia TV
  - ESPHome devices (bed presence, air quality)

- Home 2:
  - Apple Mac Mini M4, running Go
  - Apple TV 4K (this is the HomeKit controller since it can be Ethernet connected)
  - Aqara M3 Hub (Zigbee Sensors, Lights, Switches, Plugs)
  - 2 x Aqara M2 Hub (Zigbee Sensors, Lights, Switches, Plugs)
  - Apple HomePod
  - GL.inet routers
  - Wiz and Aqara lighting
  - Aqara Switches
  - LG TV's (HomeKit compatible)
  - ESPHome devices (room presence, bed presence, air quality)

## Apple HomePod

Acts as a HomeKit Hub (Server) and also is able to serve `announcements` to all HomePods in the house. This is useful for automations that need to announce something, like when the door bell is pressed or when the front door is opened or when the wash machine or dryer is done.

## Mac Mini As Home-Server

Doing it with a Mac Mini allows us to play audio on a HomePod by switching the `audio output` on the Mac Mini, we can do this from the terminal. This means we can write a script that can convert text to an mp4 and then play it on one or more HomePods. This is useful for automations that need to announce something, like when the door bell is pressed or when the front door is opened or when the wash machine or dryer is done.
We can do this by writing to an announcement.queue file, and we have a process that tails the file and plays the announcements one at a time on the designated HomePod(s).

Example announcement on HomePod1 and HomePod2:
```
say "Hello, this is a test announcement" --output announcement.mp4
select audio output as HomePod1
afplay announcement.mp4
select audio output as HomePod2
afplay announcement.mp4
```

## HomeKit Accessory Server

For simulating HK accessories, so that in your automations you can use those `virtual` accessories. Our HomeKit automation can easily incorporate these virtual devices, and we can use them to trigger automations but also trigger them.

We can use Golang and C#, our favorite languages to implement this. The HomeKit Accessory Server is a simple server that can be run on any device that supports Golang or C#. It can be run on a PC, a Mac, or any other device that supports Golang or C#.

### NFC Tags

To keep this simple for each Tag we will have a virtual accessory that is triggered when the tag is scanned. This accessory can be used in automations to trigger other accessories or automations.

So: Scan NFC Tag -> Set the state of a specific `Virtual Switch` to ON
