- Number of WiFi APs: 3
  - GL-inet Berry, Slate and Flint

- Number of presence sensors: 9
  - USB-C powered
  - WiFi connected
- Number of environment sensor boxes: 7
  - USB-C powered
  - WiFi connected
- Number of open/close sensors: 8
  - Battery powered
  - WiFi connected

# Data Streams

Every value that we track is a file data stream on our server, so how many data streams do we have?

- Total number of data streams: 74

  - Presence sensors: 9 data streams (1 per sensor)
  - Open/close sensors: 8 data streams (1 per sensor)
  - Environment sensors: 28 data streams (4 per sensor: temperature, humidity, luminance, CO2)
  - Switch wall panels: 25 data streams (varies depending on number of buttons per panel)
  - Smart Sockets (voltage/current/power/energy): 1 * 4 data streams

# Appartment

According to the below setup, the appartment contains:
  
  -  Number of '(open/close) sensor' = 8
  -  Number of 'Presence sensor' = 9
  -  Number of 'Environment sensor' = 7
  -  Number of 'Switch wall panel (2/3 buttons)' = 8

There are some 'hotel' switches, they will have to be disabled, since some of them do not contain a neutral wire and we will
use some of the 'switching' wires to trace neutral to the switch box.

## Entrance
  
  - (open/close) sensor
  - Hallway light
  - Switch wall panel (WiFi, NoNeutral, 2 buttons) x 1
    - This switch box doesn't have a neutral wire, so we will use a NoNeutral Zigbee switch with 2 buttons.

## Kitchen
  
  - Presence sensor
  - Environment: Temperature + humidity + luminance + CO2
  - Switch wall panel (WiFi, Neutral, 3 buttons) x 1
    - To be verified: This switch box needs to pull a neutral wire from the ceiling light to the switch box
  - Kitchen counter lights
  - Kitchen table lights
  - Kitchen entrance light

## Bathroom

  - (open/close) sensor  
  - Presence sensor
  - Environment: Temperature + humidity + luminance
  - Switch wall panel (WiFi, Neutral, 2 buttons) x 1
    - Tricky; There is a neutral wire passing through, so we will have to hijack it for use in the switch box
  - Bathroom lights

## Living Room
  
  - Presence sensor
  - Presence sensor
  - (open/close) sensor, balcony  
  - Environment: Temperature + humidity + luminance + CO2
  - Switch wall panel (WiFi, Neutral, 2 buttons)
  - Switch wall panel (WiFi, Neutral, 2 buttons)
  - Switch wall panel (WiFi, Neutral, 3 buttons)
  - Switch wall panel (WiFi, Neutral, 3 buttons)
  - Living room main light
  - Living room stand lights
  - Living room chandelier

## Hall

  - Presence sensor
  - Switch wall panel (WiFi, Neutral, 2 buttons)

## Bedroom Sophia
  
  - (open/close) sensor  
  - (open/close) sensor, balcony  
  - Presence sensor
  - Environment: Temperature + humidity + luminance + CO2
  - Bedroom main light
  - Switch wall panel (WiFi, Neutral, 1 button)

## Bedroom Jennifer

  - (open/close) sensor  
  - Presence sensor
  - Environment: Temperature + humidity + luminance + CO2
  - Bedroom main light
  - Switch wall panel (WiFi, Neutral, 1 button)

## Bedroom Main

  - (open/close) sensor  
  - Presence sensor
  - Environment: Temperature + humidity + luminance + CO2
  - Bedroom main light
  - Bedroom stand lights
  - Switch wall panel (WiFi, Neutral, 2 buttons)
  - Smart Socket: For AirConditioner

##  Bathroom Main

  - (open/close) sensor  
  - Presence sensor
  - Environment: Temperature + humidity + luminance
  - Switch wall panel (WiFi, Neutral, 2 buttons)
  - Bathroom lights

## North Balcony

  - Light

##  South Balcony

  - Light
