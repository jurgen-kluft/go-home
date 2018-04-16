# go-home

Automated home using Go (krazygo.org, Go only OS for Raspberry PI)

Devices/Hardware:
- Netgear R6900
- Philips HUE
- Xiaomi Yee
- Xiaomi Aqara
- Wemo switch
- Samsung TV
- Sony Bravia TV

Note:
  There is a HUE emulator in Go, this could be used to have Alexa control virtual devices like 
  the Xiaomi Gateway light, our DualWiredWallSwitch, Power Plug etc..
  Github: https://github.com/pborges/huemulator

Note:
  InfluxDB for tracking metrics and usage of all processes.

App Structure:
- Since every process is just running it's own logic want we need is kindof pub/sub server where every process
  can register itself to specific events that it is interested in.
  - Emitter.io Server (Pub/Sub server where you can subscribe to channel(s))

- Following sub-processes:
  - Presence            (Connects to Netgear Router to obtain list of devices present on the network)
  - Flux                (Calculates Color-Temperature and Brightness per day for Hue and Yee lights)
  - AQI                 (Air Quality Index)
  - Suncalc             (Computes sun-rise, sun-set etc..)
  - Weather             (Darksky)
  - Calendar            ()
  - Shout               (Has Slack as the back-end to send messages)
  - Wemo                (Wemo devices)
  - Hue                 (Philips HUE lighting, turn on/off, change CT and BRI)
  - Yee                 (Xiaomi Yee lighting, turn on/off, change CT and BRI)
  - Xiaomi aqara        (Xiaomi Gateway connection, getting information from motion sensors and controlling switches and plugs)
  - Sony Bravia Remote  (Turn on/off Sony Bravia TV(s))
  - Samsung TV Remote   (Turn on/off Samsung TV(s))
  - Automation; reacting to all events and executing automation rules
