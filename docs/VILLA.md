# Villa

- Number of WiFi APs (every floor): 3
  - 3 x GL-inet Flint
- Zigbee Coordinator per floor
  - CC2531 USB Dongle
  - For:
    - Hue Lights
    - Any other possible Zigbee devices in future

- Number of presence sensors: 18
  - USB-C powered
  - WiFi connected
- Number of environment sensor boxes: 11
  - USB-C powered
  - WiFi connected
- Number of open/close: 31
  - Battery powered
  - WiFi connected
- Number of switch wall panels: TODO, we need to identify all of them
  - Note: We have Neutral in all switch boxes

- Number of sensor data streams (roughly): = 229
  - Presence = (state + RSSI + presence) * 18 = 54
  - Environment = (4 + state + RSSI) * 11 = 66
  - Open/Close = (1 + state + RSSI) * 31 = 93
  - Smart Sockets = (voltage + current + power + energy) * 4 = 16

## Mailbox

- (open/close) door8

## Basement

- Summary with Basement:
  - presence sensors: 2
  - environment sensors: 1
  - (open/close): 4

- Stairwell
- Storage Room
  - (open/close) door
  - (open/close) window
  - environment: temperature + humidity + CO2
  - presence
  - smart socket: 2 (freezer and dehumidifier)
- Main Room
  - (open/close) door
  - (open/close) window
  - presence

## 1st Floor

- Summary:
  - presence sensors: 6
  - environment sensors: 4
  - (open/close) sensors: 10

- Entrance
  - (open/close) door, front door
- Stairwell
  - presence
- Bathroom
  - (open/close) door
  - (open/close) window
  - environment: temperature + humidity
  - presence
- Kitchen 
  - (open/close) door
  - (open/close) window, left
  - (open/close) window, right
  - presence
  - environment: temperature + humidity
  - smart socket: 3 (fridge, dishwasher and water heater)
- Dining Area 
  - (open/close) window
  - presence
- Living Area
  - (open/close) door (sliding)
  - presence
  - environment: temperature + humidity + luminance
- Bedroom Area
  - (open/close) door
  - (open/close) window
  - presence
  - environment: temperature + humidity + luminance

## 2nd Floor

- Summary:
  - presence sensors: 6
  - environment sensors: 3
  - (open/close) sensors: 9

- Stairwell 
  - presence
- Study Room 
  - (open/close) door
  - (open/close) window
  - presence
  - environment: temperature + humidity + luminance + CO2
- Bathroom 
  - (open/close) door
  - (open/close) window
  - presence
  - environment: temperature + humidity + luminance
- Living Room 
  - (open/close) door, to stairwell
  - (open/close) door, to balcony
  - (open/close) door, between living room and washing area
  - (open/close) window, left
  - (open/close) window, right
  - presence
  - environment: temperature + humidity + luminance
- Washing Area 
  - presence
  - luminance
  - smart socket: 2 (washing machine and dryer)
- Balcony 
  - presence

## 3rd Floor

- Summary:
  - presence sensors: 4
  - environment sensors: 3
  - (open/close) sensors: 8

- Stairwell 
  - presence
- Bathroom 
  - (open/close) door
  - (open/close) window
  - presence
  - environment: temperature + humidity + luminance
- Bed Room 
  - (open/close) door
  - (open/close) window
  - presence
  - environment: temperature + humidity + luminance + CO2
- Main Bed Room
  - (open/close) door
  - (open/close) door, north balcony
  - (open/close) door, south balcony
  - (open/close) window
  - presence
  - environment: temperature + humidity + luminance + CO2
- South Balcony 
- North Balcony (small)
