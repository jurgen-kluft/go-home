# go-home

Automated home using NATS (messaging), Redis (database) and micro services.


# NATS

NATS (https://github.com/nats-io/gnatsd) is used as the PubSub server since they have a very solid go client.


# NATS golang client

Has a JSON encode/decode mechanic so you can push Go structures over the wire directly.
This makes it very convenient to communicate between the services.


