# Automations

## Automation High-Level Rules

- People's Intent is more important than Automation
- Turning lights OFF automatically when no one is in the room is OK, however it should be marked that the lights where turned off by automation so that when people enter the room again, the lights are automatically turned ON.

So per room the state of the lights have an additional marker called `last change`, which can be either `automation` or `user`.

e.g. If someone goes to take a nap during the day in their bedroom, the lights should not be turned ON automatically when someone enters the room. First of all, very likely the room OCCUPANCY is 1 (the person taking a nap), and second of all, the lights where turned OFF by the user, so we should not turn them ON automatically.

## House, Per Floor, Per Room, State

For example, we want to know if the kitchen lights should be ON or OFF based on the time of day, season, and other sensors.

## Absolutely necessary

- Saving: Lights should not stay on when no one is in the room/area, especially at night.
- Auto: House should know the light intensity (per room) so as to start turning on/off lights.
- Moods: Lights should be dimmed when watching TV or a movie and brighten when paused.
- Door and Window sensors (open/close) should be used everywhere (50 RMB)
  - esp32-c3 super mini
  - US1881 hall effect sensor or reed switch
  - lipo battery (3.7V 1300mAh)
  - magnet

Example:

It is near bed time, automation detects there are 2 people in the bedroom, and it detects they are both in bed. Lights are being turned off, but on the 1st and 2nd floor some lights are still on. The automation should turn off those lights on the 1st and 2nd floor.

Furthermore, If possible, the front door should be locked and the alarm should be armed.

In the morning, when both people are out of bed and people are detected on the 1st floor, the alarm should be disarmed and the front door should be unlocked.

## Timing

A `virtual switch` that turns ON 30 minutes before sunset. This can be used to trigger other automations.
A `virtual switch` that turns ON 30 minutes before sunrise. This can be used to trigger other automations.

## HomeKit Scenes

Every scene in HomeKit has an associated `virtual switch`, this can be used to trigger the scene. 

Example 3:

Each motion sensor will have a `virtual sensor` that is set to ON when motion is detected and set to OFF when no motion is detected. 

We can then have an automation that uses these `virtual sensor`s to determine if the kitchen lights should be ON or OFF.
This automation can also use more information like the time of day, season, and other sensors to determine if the kitchen lights should be ON or OFF.

## Appartment

- Main Bedroom
  - Lights
    - Main Light
    - Stand Jurgen
    - Stand Faith
    - Shower Light
  - Temperature/Humidity/Pressure/CO2
  - Bed Presence
- Sophia Room
  - Lights
    - Main Light
  - Temperature/Humidity/Pressure/CO2
- Jennifer Room
  - Lights
    - Main Light
  - Temperature/Humidity/Pressure/CO2
- Kitchen
  - Lights
    - Diner
    - Counter
  - Air Quality
  - Presence
- Living Room
  - Lights
    - Main
    - Stand
    - Chandelier
  - Temperature/Humidity/Pressure/CO2
  - Presence
- Bathroom
  - Lights
    - Main Light
    - Shower
  - Presence
- Entrance
  - Lights
    - Frontdoor Light
  - Presence

## Virtual Switches

So we need to create a `virtual switch` for each `light` in the house.
In total we have 13 `virtual switch`es.

We can also make N `moods`

