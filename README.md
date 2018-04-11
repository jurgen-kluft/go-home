# go-home

Automated home using Go (krazygo.org, Go only OS for Raspberry PI)

Devices:
- Philips HUE
- Yeelight
- Xiaomi Aqara
- Wemo
- Samsung TV
- Sony Bravia TV

Modules:
- Netgear for attached devices
- Flux type lighting computation
- Bayesian probability on presence
- Weather attributes from Darksky
- AQI, air quality index
- Suncalc
- Calendar (reading Apple icloud public calendar(s))
- Shout messaging using Slack

Note:
  There is a HUE emulator (Java, NodeJS), this could be used to have Alexa control virtual devices like 
  the Xiaomi Gateway light, our DualWiredWallSwitch, Power Plug etc..

Note:
  InfluxDB for tracking metrics and usage of all processes.

App Structure:
- Since every process is just running it's own logic want we need is kindof pub/sub server where every process
  can register itself to specific events that it is interested in.
  - Emitter.io Server (Pub/Sub server where you can subscribe to channel(s))

- Following sub-processes:
  - Presence
  - Flux
  - AQI
  - Suncalc
  - Weather (Darksky)
  - Calendar
  - Shout
  - Wemo
  - Hue
  - Yee
  - Xiaomi aqara
  - Sony Bravia Remote
  - Samsung TV Remote
  - Automation; reacting to all events and executing automation rules
