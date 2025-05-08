#ifndef ZEPHYR_DRIVERS_SENSOR_RD_03D_H
#define ZEPHYR_DRIVERS_SENSOR_RD_03D_H

#include <zephyr/kernel.h>
#include <zephyr/device.h>
#include <zephyr/drivers/uart.h>
#include <zephyr/drivers/sensor/rd-03d.h>

#define RD03D_TX_BUF_MAX_LEN 18
#define RD03D_RX_BUF_MAX_LEN 30
#define RD03D_UART_BAUD_RATE 115200

/* Arbitrary max duration to wait for the response */
#define RD03D_WAIT K_SECONDS(1)

enum rd03d_protocol_cmd_idx {
	RD03D_PROTOCOL_CMD_IDX_OPEN_CMD_MODE,
	RD03D_PROTOCOL_CMD_IDX_CLOSE_CMD_MODE,
	RD03D_PROTOCOL_CMD_IDX_DEBUGGING_MODE,
	RD03D_PROTOCOL_CMD_IDX_REPORTING_MODE,
	RD03D_PROTOCOL_CMD_IDX_RUNNING_MODE,

	RD03D_PROTOCOL_CMD_IDX_SET_MIN_DISTANCE,
	RD03D_PROTOCOL_CMD_IDX_SET_MAX_DISTANCE,
	RD03D_PROTOCOL_CMD_IDX_SET_MIN_FRAMES, // Minimum number of frames for being considered as
					       // appearing
	RD03D_PROTOCOL_CMD_IDX_SET_MAX_FRAMES, // Maximum number of frames for being considered as
					       // dissapearing
	RD03D_PROTOCOL_CMD_IDX_SET_DELAY_TIME, // The delay time to use for considering the object
					       // as dissapearing

	RD03D_PROTOCOL_CMD_IDX_GET_MIN_DISTANCE,
	RD03D_PROTOCOL_CMD_IDX_GET_MAX_DISTANCE,
	RD03D_PROTOCOL_CMD_IDX_GET_MIN_FRAMES,
	RD03D_PROTOCOL_CMD_IDX_GET_MAX_FRAMES,
	RD03D_PROTOCOL_CMD_IDX_GET_DELAY_TIME,

	RD03D_PROTOCOL_CMD_IDX_SINGLE_TARGET_MODE, // Detection mode single target
	RD03D_PROTOCOL_CMD_IDX_MULTI_TARGET_MODE,  // Detection mode multiple targets

	RD03D_PROTOCOL_CMD_IDX_MAX,
};

enum rd03d_protocol_ack_idx {
	RD03D_PROTOCOL_CMD_ACK_IDX_OPEN_CMD_MODE,
	RD03D_PROTOCOL_CMD_ACK_IDX_CLOSE_CMD_MODE,
	RD03D_PROTOCOL_CMD_ACK_IDX_MAX,
};

enum rd03d_detection_mode {
	RD03D_DETECTION_MODE_SINGLE_TARGET,
	RD03D_DETECTION_MODE_MULTI_TARGET,
};

enum rd03d_operation_mode {
	RD03D_OPERATION_MODE_CMD = 0x80,
	RD03D_OPERATION_MODE_DEBUG = 0x0,
	RD03D_OPERATION_MODE_REPORT = 0x01,
	RD03D_OPERATION_MODE_RUN = 0x02,
};

enum rd03d_property {
	RD03D_PROP_MIN_DISTANCE,
	RD03D_PROP_MAX_DISTANCE,
	RD03D_PROP_MIN_FRAMES,
	RD03D_PROP_MAX_FRAMES,
	RD03D_PROP_DELAY_TIME,
	RD03D_PROP_DETECTION_MODE,
	RD03D_PROP_OPERATION_MODE,
};

#define RD03D_MAX_TARGETS 3

/** @brief Information collected from the sensor on each fetch. */
struct rd03d_report {

	struct rd03d_target {
		uint16_t x;        /**< X coordinate in mm */
		uint16_t y;        /**< Y coordinate in mm */
		uint16_t distance; /**< Distance in mm, 0 means target is invalid */
		uint16_t speed;    /**< Speed of the target in cm/s */
	};

	rd03d_target targets[RD03D_MAX_TARGETS]; /**< Array of targets detected */
};

struct rd03d_data {
	uint8_t tx_bytes;    /* Number of bytes that are send */
	uint8_t tx_data_len; /* Length of the buffer to send */
	uint8_t tx_data[RD03D_TX_BUF_MAX_LEN];

	uint16_t rx_bytes;
	uint8_t rx_len;   /* Length of the buffer to receive (header begin + 2 + intra-frame data length) + header end*/
	uint8_t rx_data[RD03D_RX_BUF_MAX_LEN];

	uint8_t has_rsp;

	uint8_t operation_mode;
	uint8_t detection_mode;

	struct k_sem tx_sem;
	struct k_sem rx_sem;

	rd03d_report report;
};

struct rd03d_cfg {
	const struct device *uart_dev;
	uint16_t min_distance;
	uint16_t max_distance;
	uint16_t min_frames;
	uint16_t max_frames;
	uint16_t delay_time;
	uart_irq_callback_user_data_t cb;
};

#endif /* ZEPHYR_DRIVERS_SENSOR_RD_03D_H */
