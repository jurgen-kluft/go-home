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

Lowest In-Bed: 10 K Ohm
Highest In-Bed: 30 K Ohm

Out-of-Bed: The resistance is very high, so the FSR is not conductive. The resistance is in the MOhm range.

R1 = SQRT(10 K Ohm * 3000 K Ohm) = SQRT(30000 K Ohm) = 170 K Ohm

## ESP32 Code

The following code can read the ADC value from the FSR and print it to the serial monitor.

- URL: https://github.com/espressif/arduino-esp32/blob/master/libraries/ESP32/examples/AnalogRead/AnalogRead.ino


