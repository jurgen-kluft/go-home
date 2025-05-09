/*
 * Copyright (c) 2025
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * @file
 * @brief Extended public API for RD-03D Sensor
 *
 * Some capabilities and operational requirements for this sensor
 * cannot be expressed within the sensor driver abstraction.
 */

#ifndef ZEPHYR_INCLUDE_DRIVERS_SENSOR_RD_03D_H_
#define ZEPHYR_INCLUDE_DRIVERS_SENSOR_RD_03D_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <zephyr/drivers/sensor.h>

/*
Base channels that are working:
  - SENSOR_CHAN_PROX
    - attribute = SENSOR_ATTR_RD03D_TARGETS
      - val1 = (3 bits) 
	    - binary '0001', only tracking target 0
	    - binary '0011', tracking target 0 & 1
	    - binary '0110', tracking target 1 & 2
    - attribute = SENSOR_ATTR_RD03D_TARGET_0 / 1 / 2
      - val1 = (1 bit, 0 or 1)
	    - binary '0' target is unavailable
	    - binary '1' target is tracking
	  - if (val1 == 1) val2 = distance in mm
  - SENSOR_CHAN_DISTANCE
    - attribute = SENSOR_ATTR_RD03D_TARGET_0 / 1 / 2
    - val1 == 1, target is tracking
    - val1 == 0, target is unavailable
    - if (val1 == 1) val2 = distance in mm 
*/

enum sensor_channel_rd03d {
	/**
	 * Channels to configure the sensor
	 */
	SENSOR_CHAN_RD03D_CONFIG_DISTANCE = SENSOR_CHAN_PRIV_START,
	SENSOR_CHAN_RD03D_CONFIG_FRAMES,
	SENSOR_CHAN_RD03D_CONFIG_DELAY_TIME,
	SENSOR_CHAN_RD03D_CONFIG_DETECTION_MODE,
	SENSOR_CHAN_RD03D_CONFIG_OPERATION_MODE,
	/*
	 * Return the X (in val1), Y (in val2) of the target (in mm)
	 */
	SENSOR_CHAN_RD03D_POS,
	/*
	 * Return the speed of the target (val1, in cm/s)
	 */
	SENSOR_CHAN_RD03D_SPEED,
	/*
	 * Return the distance to the target (val1, in mm)
	 */
	SENSOR_CHAN_RD03D_DISTANCE,
};

enum sensor_attribute_rd03d {
	SENSOR_ATTR_RD03D_TARGETS = SENSOR_ATTR_PRIV_START,

	SENSOR_ATTR_RD03D_TARGET_0,
	SENSOR_ATTR_RD03D_TARGET_1,
	SENSOR_ATTR_RD03D_TARGET_2,

	SENSOR_ATTR_RD03D_CONFIG_VALUE,
	SENSOR_ATTR_RD03D_CONFIG_MINIMUM,
	SENSOR_ATTR_RD03D_CONFIG_MAXIMUM,
};

#ifdef __cplusplus
}
#endif

#endif /* ZEPHYR_INCLUDE_DRIVERS_SENSOR_RD_03D_H_ */
