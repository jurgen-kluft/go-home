# I2C

- SCL; This is the clock line.
- SDA; This is the data line.

To have multiple sensors on the same bus, you need to make sure that each sensor has a different address. The address is usually set by the manufacturer and can be found in the datasheet. If you have multiple sensors with the same address, you can use a multiplexer to switch between them.

For connecting multiple sensors, you need to connect them in parallel. This means that you connect the SCL and SDA lines of all sensors together. You also need to connect the ground and power lines of all sensors together. The power line is usually 3.3V or 5V, depending on the sensor. The ground line is usually connected to the ground of the microcontroller.

