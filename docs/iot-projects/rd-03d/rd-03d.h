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

enum sensor_channel_rd03d {
	/**
	 * Channels to set/get specifics of the sensor
	 */
	SENSOR_CHAN_RD03D_MIN_DISTANCE = SENSOR_CHAN_PRIV_START,
	SENSOR_CHAN_RD03D_MAX_DISTANCE,
	SENSOR_CHAN_RD03D_MIN_FRAMES,
	SENSOR_CHAN_RD03D_MAX_FRAMES,
	SENSOR_CHAN_RD03D_DELAY_TIME,
	SENSOR_CHAN_RD03D_DETECTION_MODE,
	SENSOR_CHAN_RD03D_OPERATION_MODE,

	/**
	 * Channels to get position (X/Y) for target 0, 1 or 2
	 */
	SENSOR_CHAN_RD03D_T0_POS,
	SENSOR_CHAN_RD03D_T1_POS,
	SENSOR_CHAN_RD03D_T2_POS,

	/**
	 * Channels to get speed (cm/s) for target 0, 1 or 2
	 */
	SENSOR_CHAN_RD03D_T0_SPEED,
	SENSOR_CHAN_RD03D_T1_SPEED,
	SENSOR_CHAN_RD03D_T2_SPEED,

	/**
	 * Channels to get distance (mm) for target 0, 1 or 2
	 */
	SENSOR_CHAN_RD03D_T0_DISTANCE,
	SENSOR_CHAN_RD03D_T1_DISTANCE,
	SENSOR_CHAN_RD03D_T2_DISTANCE,

};


#ifdef __cplusplus
}
#endif

#endif /* ZEPHYR_INCLUDE_DRIVERS_SENSOR_RD_03D_H_ */
