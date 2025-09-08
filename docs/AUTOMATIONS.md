# Automations

## House, Per Floor, Per Room, State

For example, we want to know if the kitchen lights should be ON or OFF based on the time of day, season, and other sensors.

## Absolutely necessary

- Lights should not stay on when no one is in the room/area, especially at night.
- House should know the light intensity so as to start turning on/off lights based on the time of day and season.
- Lights should be dimmed when watching TV or a movie and brighten when paused.

Example:

It is near bed time, automation detects there are 2 people in the bedroom, and it detects they are both in bed. Lights are being turned off, but on the 1st and 2nd floor some lights are still on. The automation should turn off those lights on the 1st and 2nd floor.

Furthermore, If possible, the front door should be locked and the alarm should be armed.

In the morning, when both people are out of bed and people are detected on the 1st floor, the alarm should be disarmed and the front door should be unlocked.

## Timing

A `virtual switch` that turns ON 30 minutes before sunset. This can be used to trigger other automations.

## HomeKit Scenes


Every scene in HomeKit has an associated `virtual switch`, this can be used to trigger the scene. 


Example 3:

Each motion sensor will have a `virtual sensor` that is set to ON when motion is detected and set to OFF when no motion is detected. 

We can then have an automation that uses these `virtual sensor`s to determine if the kitchen lights should be ON or OFF.
This automation can also use more information like the time of day, season, and other sensors to determine if the kitchen lights should be ON or OFF.


## Appartment

- Main Bedroom
  - Main Light
  - Stand Jurgen
  - Stand Faith
  - Shower Light
- Sophia Room
  - Main Light
- Jennifer Room
  - Main Light
- Kitchen
  - Diner
  - Counter
- Living Room
  - Main
  - Stand
  - Chandelier
- Bathroom
  - Main Light
  - Shower
- Entrance
  - Frontdoor Light

So we need to create a `virtual switch` for each `light` in the house.
In total we have 13 `virtual switch`es.

We can also make N `moods`

