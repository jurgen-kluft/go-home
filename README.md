# go-home

<<<<<<< HEAD
Automated home (currently only supporting WEMO and HUE)


=======
Automated home using Go (krazygo.org, Go only OS for Raspberry PI)
>>>>>>> bd2165cd9aca0b4f932c61aae51b4b8774e15eed

Devices:
- Philips HUE
- Yeelight
- Xiaomi Aqara

<<<<<<< HEAD
Presence
- This service is ticking itself and will push messages

TimeOfDay
- This service is ticking itself and will push messages

Devices (HUE and WEMO)
- This service is ticking itself and will push messages


Messaging (HipChat)
- This service is receiving messages

Logging
- This service is receiving messages




# Communication BUS

NATS: https://github.com/nats-io/gnatsd, https://nats.io/
I still think we should go this way, NATS is a simple installation.


# Configuration

Consul: https://github.com/hashicorp/consul
Consul can help to keep the configuration dynamic, also it is just one file so managing this is easy.
It can also act as a KV store and that is all we need for go-home.
=======
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
>>>>>>> bd2165cd9aca0b4f932c61aae51b4b8774e15eed
