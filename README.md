# go-home

HomeKit based home automation using Go, Arduino (C/C++, ESP32 and ESP8266).

## Status

Total overhaul of the project, all code is being scratched and redone in Go.

## App Structure

Note:
  Since all of this is written in Go it should be able to run anywhere, from a Raspberry PI, Arduino to any other Linux based machine.

- Since every process is just running it's own logic want we need is client/server mechanism where every process
  can connect itself to for specific events that it is interested in.

- Following sub-processes:
  - Apple HomeKit       WIP,(Apple Home Kit bridge, Hue Bridge and more are connected to this, and it exposes all devices)
  - GeekOpen            Ok, (GeekOpen WiFi wall-sockets/wall-switches/power-plugs)
  - Config              Ok, (A service that is the provider of configurations for all other services)
  - Presence            Ok, (Connects to Router to obtain list of devices present on the network)
  - Flux                Ok, (Calculates Color-Temperature and Brightness during the day and makes sensors out of them)
  - AQI                 Ok, (Air Quality Index, and makes sensor out of it)
  - Suncalc             Ok, (Computes sun-rise, sun-set etc.., and makes sensors out of it)
  - Season              Ok, (Computes the season (Spring, Summer, Autumn, Winter), and makes sensors out of it)
  - Weather             Ok, (Obtains weather data from OpenWeatherMap and makes sensors out of it)
  - Calendar            Ok, (Reads calendars from icloud and determines active events and makes sensors out of them)
  - Shout               Ok, (Push notifications to HomePods?)
  - Yee                 Ok, (Xiaomi Yee lighting, turn on/off, change CT and BRI)
  - Wiz                 Ok, (Xiaomi Yee lighting, turn on/off, change CT and BRI)
  - Sony Bravia TV      Ok, (Turn on/off, HDMI Input, Volume for Sony Bravia TV(s))
  - Samsung TV          Ok, (Turn on/off Samsung TV(s))

TODO:
  - DD-WRT; For some form of presence (https://github.com/awilliams/wifi-presence)
  
## Automation Logic
  
Automation, reacting to all events and executing automation rules, all written in Go.

## Devices / Hardware

- Home 1:
  - Apple Mac Mini M4, running Go
  - Apple TV 4K (this is the HomeKit controller since it can be Ethernet connected)
  - Apple HomePod
  - GLiNet routers
  - Hue, Yee, Ikea, Wiz lighting
  - Geek-Open, WiFi Switches/Sockets (can report energy usage)
    - Can detect if dish-washer, washing machine, air-conditioners or dryer is running based on energy usage
  - Sony Bravia TV
  - Custom (WiFi, UDP/TCP) sensors: 
    - door/window
    - motion/presence sensors
    - temperature/humidity sensors
    - light sensors (all based on ESP32/ESP8266 and Arduino)
  
- Home 2:
  - Apple Mac Mini M4, running Go
  - Apple TV 4K (this is the HomeKit controller since it can be Ethernet connected)
  - Apple HomePod
  - GL.inet routers
  - Hue, Ikea, Wiz lighting
  - Geek-Open, WiFi Switches/Sockets
    - Can detect if dish-washer, washing machine, air-conditioners or dryer is running based on energy usage
  - LG TV's (HomeKit compatible)
  - Custom (WiFi, UDP/TCP) sensors: 
    - door/window (magnetic switch)
    - motion/presence sensors
    - temperature/humidity sensors
    - light sensors (all based on ESP32/ESP8266 and Arduino)
  
## Apple HomePod

Acts as a HomeKit Hub (Server) and also is able to serve `announcements` to all HomePods in the house. This is useful for automations that need to announce something, like when the door bell is pressed or when the front door is opened or when the wash machine or dryer is done.

## HomeKit Accessory Server

For simulating HK accessories, so that in your automations you can use those `virtual` accessories. Our HomeKit automation can easily incorporate these virtual devices, and we can use them to trigger automations but also trigger them.

We can use Golang and C#, our favorite languages to implement this. The HomeKit Accessory Server is a simple server that can be run on any device that supports Golang or C#. It can be run on a PC, a Mac, or any other device that supports Golang or C#.

