# Services

Each service is responsible for a specific domain, e.g. weather, calendar, presence, etc. Each service has its own configuration and can communicate with other services through the micro-service framework. Also each service watches its own configuration file for changes and can restart itself and load the new configuration when it detects a change.

## UNIX Socket Messaging

The micro-service framework uses UNIX sockets for inter-process communication. Each service has its own socket file, and the framework manages the connections between services based on their dependencies. When a service starts, it creates its socket file and listens for incoming connections. When a service needs to communicate with another service, it connects to the respective socket file of that service.

## Configuration

Unix Socket configuration for multiple processes, where the outgoing sockets are determined from the incoming sockets of the respective services.

## 