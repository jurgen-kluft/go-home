# go-home

Automated home using Go (krazygo.org, Go only OS for Raspberry PI)

Devices/Hardware:

- Netgear R6900
- Philips HUE
- Xiaomi Yee
- Xiaomi Aqara
  - Wall Switch
  - Button Switch
  - Power Switch
  - Sensors (Motion, Magnet)
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
  - Hue Emulator        (Philips HUE lighting emulator, turn on/off)
  - Yee                 (Xiaomi Yee lighting, turn on/off, change CT and BRI)
  - Xiaomi aqara        (Xiaomi Gateway connection, getting information from motion sensors and controlling switches and plugs)
  - Sony Bravia Remote  (Turn on/off Sony Bravia TV(s))
  - Samsung TV Remote   (Turn on/off Samsung TV(s))
  - Automation; reacting to all events and executing automation rules

Status:

I am currently not yet running this, however i have tested all sub-processes and they all function. The only one not
tested yet is 'automation'.
