# Resistance Setup

See `resistance_setup.png` for the schematic.

The choice of R1 is very important for the largest possible detection range. To determine its value, measure the resistance of the FSR with a multimeter when you are in and out of bed. R1 can then be found using this formula:

`R1 = SQRT( Rin_bed * Rout_of_bed)`

Make the measurements a few times and take averages for each of the two values.

The result will depend on the thickness, weight and resilience of your mattress. For one bed I used 300 Ohms, for another 3K Ohms!. Do calculate this value properly.

The Vout should be connected to an ADC pin of the ESP32. The Vcc should be connected to 3.3V and GND to GND.

## Hardware

- ESP32 S3 DevKit
- 0.5 mm wire, black, 5 meter
- 0.75 mm wire, blue, 5 meter
- 2 x FSR (Force Sensitive Resistor) 60 cm long
- 2 x 100K Ohm resistor and 2 x 60K Ohm resistor

## Functionality

A sliding window of 10 seconds is used to determine presence ON or OFF. The ESP32 will read the ADC value every 100 ms and if the value is above a certain threshold, it will be considered as presence ON. If the value is below the threshold for 10 seconds, it will be considered as presence OFF.

The presence will be posted directly to the MQTT broker. The ESP32 will be connected to the WiFi network and will send the data to the MQTT broker.

## Measurement

--------------------------------------------------------------------

Side of Bed = Jurgen (15th May 2025)
     Lowest In-Bed:      ~15 K Ohm  /   -20 K Ohm
     Out-Of-Bed:         60 K Ohm
R1 = SQRT(15 K Ohm * 60 K Ohm) = SQRT(900 K Ohm) = 30 K Ohm = 20 K Ohm + 10 K Ohm

Side of Bed = Jurgen (4th June 2025)
     Lowest In-Bed:      ~40 K Ohm  /   -50 K Ohm
     Out-Of-Bed:         600 K Ohm
R1 = SQRT(40 K Ohm * 600 K Ohm) = SQRT(24000 K Ohm) = 155 K Ohm = 100 K Ohm + 50 K Ohm + 5 K Ohm


--------------------------------------------------------------------

Side of Bed = Faith (15th May 2025)
     Lowest In-Bed:     ~5 K Ohm    /    -10 K Ohm
     Out-Of-Bed:        25 K Ohm    /     35 K Ohm

R1 = SQRT(5 K Ohm * 35 K Ohm) = SQRT(175 K Ohm) = 13.3 K Ohm = 10 K Ohm + 1 K Ohm + 1 K Ohm + 1 K Ohm


Side of Bed = Faith (4th June 2025)
     Lowest In-Bed:     ~5 K Ohm   /    -10 K Ohm
     Out-Of-Bed:       200 K Ohm   /    300 K Ohm

R1 = SQRT(5 K Ohm * 300 K Ohm) = SQRT(1500 K Ohm) = 38.7 K Ohm = 30 K Ohm + 8.2 K Ohm


--------------------------------------------------------------------

## ESP32 Code

The following code can read the ADC value from the FSR and print it to the serial monitor.

- URL: https://github.com/espressif/arduino-esp32/blob/master/libraries/ESP32/examples/AnalogRead/AnalogRead.ino


