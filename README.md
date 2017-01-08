# go-home

Automated home (currently only supporting WEMO and HUE)



# Services

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
