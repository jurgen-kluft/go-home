# ESP NOW and WIFI

Seems you can use ESP-NOW and WiFi at the same time, but you need to set the WiFi mode to AP_STA.
This means for WiFi you will have to have someone connect to the 'gateway' ESP32, and then the 
ESP32 will be able to send data out on WiFi.
All other ESP32s will be able to send data to the 'gateway' ESP32 using ESP-NOW.

- https://github.com/jonathanrandall/esp32-esp32-now-wifi

# OLED Display

We still have some lying around, from our split keyboard projects.

- SSD1306 OLED display
  - I2C address: 0x3C
  - Voltage: Seems to be able to handle 3.3V and 5V ??
  - 

# 10cm x 10cm x 3.2cm case

Sufficient to build all the sensor units:

- Air quality
- Air quality + Room Presence
- Bed presence

# Soldering 4P (2.54 mm) Connector

Every sensor has a 4P connector, so we can easily connect/disconnect.

# I2C

- SCL; This is the clock line.
- SDA; This is the data line.

To have multiple sensors on the same bus, you need to make sure that each sensor has a different address. The address is usually set by the manufacturer and can be found in the datasheet. If you have multiple sensors with the same address, you can use a multiplexer to switch between them.

For connecting multiple sensors, you need to connect them in parallel. This means that you connect the SCL and SDA lines of all sensors together. You also need to connect the ground and power lines of all sensors together. The power line is usually 3.3V or 5V, depending on the sensor. The ground line is usually connected to the ground of the microcontroller.

