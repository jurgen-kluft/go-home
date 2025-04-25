# Automations

## States

Mostly we want to compute states based local information and/or other states. 

For example, we want to know if the kitchen lights should be ON or OFF based on the time of day, season, and other sensors.

## NFC Tags

Let's say we want a `virtual switch` that is set to ON when a NFC tag is scanned, we can create a `virtual switch` that is triggered, and then we can use this `virtual switch` in our automations. Other automations could turn off the `virtual switch`.

## Timing

A `virtual switch` that turns ON 30 minutes before sunset. This can be used to trigger other automations.

## HomeKit Scenes

`
Every scene in HomeKit has an associated `virtual switch`, this can be used to trigger the scene. 
`

Example 3:

`
Each motion sensor will have a `virtual sensor` that is set to ON when motion is detected and set to OFF when no motion is detected. 

We can then have an automation that uses these `virtual sensor`s to determine if the kitchen lights should be ON or OFF.
This automation can also use more information like the time of day, season, and other sensors to determine if the kitchen lights should be ON or OFF.

`

## Appartment

- Bedroom
  - Main
  - Stand Jurgen
  - Stand Faith
  - Shower
- Sophia
  - Main
- Jennifer
  - Main
- Kitchen
  - Diner
  - Counter
- Living Room
  - Main
  - Stand
  - Chandelier
- Bathroom
  - Shower
- Entrance
  - Frontdoor Light

So we need to create a `virtual switch` for each `light` in the house.
In total we have 13 `virtual switch`es.

We can also make N `moods`