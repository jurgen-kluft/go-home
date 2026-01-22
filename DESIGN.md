# Unix Socket

## UNIX Socket Messaging

Listen to config/request messages, when a client connects we can send him the configuration.
Also when we detect that the configuration on disk has changed, we can hot-load it and send
it to the associated channel. (This is part of the config service and is working)

