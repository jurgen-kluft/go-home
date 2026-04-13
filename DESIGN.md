# Services

Each service is responsible for a specific domain, e.g. weather, calendar, presence, etc. Each service has its own configuration and can communicate with other services through UNIX sockets. Also each service watches its own configuration file for changes and can hot-load the new configuration when it detects a change.

## UNIX Socket Messaging

Listen to config/request messages, when a client connects we can send him the configuration.
Also when we detect that the configuration on disk has changed, we can hot-load it and send
it to the associated channel. (This is part of the config service and is working)

## Configuration

Unix Socket configuration for multiple processes, where the outgoing sockets are determined from the incoming sockets of the respective services:

```json
{
  "services": [
    {
      "id": "weather"
    },
    {
      "id": "aqi"
    },
    {
      "id": "sun"
    },
    {
      "id": "calendar"
    },
    {
      "id": "presence"
    },
    {
      "id": "flux",
      "depends_on": [
        "calendar",
        "weather",
        "sun"
      ]
    },
    {
      "id": "automation",
      "depends_on": [
        "calendar",
        "presence"
      ]
    }
  ]
}
```