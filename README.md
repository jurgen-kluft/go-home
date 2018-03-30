# go-home

Automated home using Go (krazygo.org, Go only OS for Raspberry PI)

Devices:
- Philips HUE
- Yeelight
- Xiaomi Aqara

Modules:
- Netgear for attached devices
- Flux type lighting computation
- Bayesian probability on presence (convert from C# to Go)
- Weather attributes from Darksky
- AQI
- Suncalc
- Calendar (reading Apple icloud public calendar(s))
- Slack
- Wemo

Note:
  There is a HUE emulator (Java, NodeJS), this could be used to have Alexa control virtual devices like 
  the Xiaomi Gateway light, our DualWiredWallSwitch, Power Plug etc..

Note:
  Promotheus for tracking metrics and state of all processes.

Note:
  Log (info, warning, error) module that sends messages to a Log-Actor which in turn can save it on disk, write to console or whatnot.

App Structure:
- Since every process is just running it's own logic want we need is kindof pub/sub server where every process
  can register itself to specific events that it is interested in.
  - Mist Server (Pub/Sub server where you can subscribe to 'tagged' messages)

- Following sub-processes:
  - Log
  - Presence
  - Flux
  - AQI
  - Suncalc
  - TimeOfDay
  - Weather (Darksky)
  - Calendar
  - Slack
  - Wemo
  - Lighting (Hue and Yee)
  - Xiaomi aqara
  - Sony Bravia Remote
  - Samsung TV Remote
  - Automation; reacting to all events and executing automation rules
