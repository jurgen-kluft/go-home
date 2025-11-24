# Connectivity

- WiFi
- Ethernet
- EspNow
- Bluetooth (BLE, 5 and 6, see Nordic nRF54L15)

## Central Gateway

The Central Gateway is a dedicated ESP32 device using Ethernet that act as a bridge between 
ESP-NOW devices and the LAN network. It will receive data from multiple ESP-NOW devices 
and forward this data over ethernet to a central server.
When the central-server on the LAN is not reachable, the Central Gateway will buffer the data
locally on a TF card until the server is reachable again.

The Central Gateway also has a procedure to 'search' for the Central Server on the LAN network.
It will broadcast a UDP message on the LAN network requesting the server to identify itself.
When the server receives this message it will respond with its IP address and port number.
The Central Gateway will then store this information in its non-volatile memory for future use.

TODO:

- Design 3D case for Central Gateway
- Integrate the LilyGo T-Ethernet board
- Implement the Central Gateway firmware


## General Boot Procedure

- On power-up, the ESP32 initializes ESP-NOW
  
- Sends a broadcast message with information that it requests credentials
- A Central Gateway (dedicated ESP32) listens for these broadcast messages and responds with:
  - WiFi SSID+Password (optional)
  - Mac Address of N Central Gateways for future ESP-NOW communication
- Sends a unicast message to each of the Central Gateways to confirm receipt of credentials.
  This will take some time and it will determine which Central Gateways are reachable.
- It will then connect to each of the Central Gateways via ESP-NOW to confirm itself as a peer.
- It can now start to send sensor data to the Central Gateways via ESP-NOW.

TODO:

- Implement the General Boot Procedure for ESP32 and ESP8266 devices.

## ESP32 using ESP-NOW

These devices are mains powered and will use ESP-NOW to send data to one or more Central Gateways.
When the device is not configured it will follow the above boot procedure to get the necessary
connectivity information.

## ESP8266 using ESP-NOW

These devices are battery powered and will use ESP-NOW to send data to a pre-programmed Central Gateway.
When the device is not configured it will follow the above boot procedure to get the necessary
connectivity information.

