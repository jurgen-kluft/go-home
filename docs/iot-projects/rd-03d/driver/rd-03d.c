#define DT_DRV_COMPAT aithinker_rd03d

#include <zephyr/logging/log.h>
#include <zephyr/sys/byteorder.h>
#include <zephyr/drivers/sensor.h>
#include "rd-03d.h"

LOG_MODULE_REGISTER(rd03d, CONFIG_SENSOR_LOG_LEVEL);

// Protocol packet begin and end
#define RD03D_CMD_HEADER_BEGIN   0xFD, 0xFC, 0xFB, 0xFA
#define RD03D_CMD_HEADER_END     0x04, 0x03, 0x02, 0x01
#define RD03D_FRAME_BEGIN        0xF4, 0xF3, 0xF2, 0xF1
#define RD03D_FRAME_END          0xF8, 0xF7, 0xF6, 0xF5
#define RD03D_REPORT_FRAME_BEGIN 0xAA, 0xFF, 0x03, 0x00
#define RD03D_REPORT_FRAME_END   0x55, 0xCC

// clang-format off

// Endianness is Little Endian

// Note: So in the comments you may read Word which means that this is a register value. 
//       The byte order of the command is thus swapped compared to an indicated Word value.

// Send command protocol frame format
// |----------------------------------------------------------------------------------
// | Frame header | Intra-frame data length  |  Intra-frame data  |   End of frame   |
// | FD FC FB FA  |  2 bytes                 |    See table 4     |   04 03 02 01    |
// |----------------------------------------------------------------------------------

// Table 4 
// Send intra-frame data format
// |----------------------------------------------------------
// |  Command Word (2 bytes)   |    Command value (N bytes)  |
// |----------------------------------------------------------

// static const uint8_t rd03d_cmds[][18] = {
// 	[RD03D_PROTOCOL_CMD_IDX_OPEN_CMD_MODE]      = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0xFF, 0x00, 0x01, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_CLOSE_CMD_MODE]     = { RD03D_CMD_HEADER_BEGIN, 0x02, 0x00, 0xFE, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_DEBUGGING_MODE]     = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_REPORTING_MODE]     = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_RUNNING_MODE]       = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_SET_MIN_DISTANCE]   = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_SET_MAX_DISTANCE]   = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END }, 
// 	[RD03D_PROTOCOL_CMD_IDX_SET_MIN_FRAMES]     = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x07, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_SET_MAX_FRAMES]     = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x07, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_SET_DELAY_TIME]     = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x07, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_GET_MIN_DISTANCE]   = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0x08, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_GET_MAX_DISTANCE]   = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0x08, 0x00, 0x01, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_GET_MIN_FRAMES]     = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0x08, 0x00, 0x02, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_GET_MAX_FRAMES]     = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0x08, 0x00, 0x03, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_GET_DELAY_TIME]     = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0x08, 0x00, 0x04, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_SINGLE_TARGET_MODE] = { RD03D_CMD_HEADER_BEGIN, 0x02, 0x00, 0x80, 0x00, RD03D_CMD_HEADER_END },
// 	[RD03D_PROTOCOL_CMD_IDX_MULTI_TARGET_MODE]  = { RD03D_CMD_HEADER_BEGIN, 0x02, 0x00, 0x90, 0x00, RD03D_CMD_HEADER_END }, 
// };

// return the length of the command
static int prepare_cmd(rd03d_protocol_cmd_idx, uint8_t *cmd_buffer, int value) {
	// Assume the header is already set
	int l = 4; // Skip header

	if (cmd_idx  == RD03D_PROTOCOL_CMD_IDX_OPEN_CMD_MODE) 
	{
		cmd_buffer[l++] = 0x04; // Intra-frame data length
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0xFF; // Command word
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x01; // 
		cmd_buffer[l++] = 0x00; //
	}
	else if (cmd_idx == RD03D_PROTOCOL_CMD_IDX_CLOSE_CMD_MODE) 
	{
		cmd_buffer[l++] = 0x02; // Intra-frame data length
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0xFE; // Command word
		cmd_buffer[l++] = 0x00; //
	} 
	else if (cmd_idx >= RD03D_PROTOCOL_CMD_IDX_DEBUGGING_MODE && cmd_idx <= RD03D_PROTOCOL_CMD_IDX_RUNNING_MODE) 
	{
		cmd_buffer[l++] = 0x08; // Intra-frame data length
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x12; // Command word
		cmd_buffer[l++] = 0x00; //
		// 6 bytes of data
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x00; //
	}
	else if (cmd_idx >= RD03D_PROTOCOL_CMD_IDX_SET_MIN_DISTANCE && cmd_idx <= RD03D_PROTOCOL_CMD_IDX_SET_DELAY_TIME) 
	{
		cmd_buffer[l++] = 0x08; // Intra-frame data length
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x07; // Command word
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = cmd_index - RD03D_PROTOCOL_CMD_IDX_SET_MIN_DISTANCE; // Command value
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = value & 0xFF; // Set value, 32-bit
		cmd_buffer[l++] = (value >> 8) & 0xFF;
		cmd_buffer[l++] = (value >> 16) & 0xFF;
		cmd_buffer[l++] = (value >> 24) & 0xFF;
	}
	else if (cmd_idx >= RD03D_PROTOCOL_CMD_IDX_GET_MIN_DISTANCE && cmd_idx <= RD03D_PROTOCOL_CMD_IDX_GET_DELAY_TIME) 
	{
		cmd_buffer[l++] = 0x04; // Intra-frame data length
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x08; // Command word
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = cmd_index - RD03D_PROTOCOL_CMD_IDX_GET_MIN_DISTANCE; // Command value
		cmd_buffer[l++] = 0x00; //
	}
	else if (cmd_idx >= RD03D_PROTOCOL_CMD_IDX_SINGLE_TARGET_MODE && cmd_idx <= RD03D_PROTOCOL_CMD_IDX_MULTI_TARGET_MODE) 
	{
		cmd_buffer[l++] = 0x02; // Intra-frame data length
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = RD03D_PROTOCOL_CMD_IDX_SINGLE_TARGET_MODE ? 0x80 : 0x90; // Command word
		cmd_buffer[l++] = 0x00; //
	}

	// Write the end
	const uint8_t* cmd_header_end[] = RD03D_CMD_HEADER_END;
	cmd_buffer[l++] = cmd_header_end[0];
	cmd_buffer[l++] = cmd_header_end[1];
	cmd_buffer[l++] = cmd_header_end[2];
	cmd_buffer[l++] = cmd_header_end[3];

	return l;
}

// ACK command protocol frame format
// |--------------------------------------------------------------------------------
// | Frame header  | Intra-frame data length  |  Intra-frame data  |  End of frame |
// | FD FC FB FA   |       2bytes             |      See table 6   |   04 03 02 01 |
// |--------------------------------------------------------------------------------

// Table 6
// ACK intra-frame data format
// |------------------------------------------------------------------
// |  Send Command Word | 0x0100 (2 bytes)  | Return value (N bytes)  |
// |------------------------------------------------------------------


static const uint8_t rd03d_acks[][1 + 18] = {
    // Protocol, acks for open and close cmd
	[RD03D_PROTOCOL_CMD_ACK_IDX_OPEN_CMD_MODE]    = { 0, RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0xFF, 0x01, 0x00, 0x00, 0x01, 0x00, 0x40, 0x00, RD03D_CMD_HEADER_END },
	[RD03D_PROTOCOL_CMD_ACK_IDX_CLOSE_CMD_MODE]   = { 0, RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0xFF, 0x01, 0x00, 0x00, RD03D_CMD_HEADER_END },
};
// clang-format on

static void rd03d_uart_flush(const struct device *uart_dev)
{
	uint8_t c;

	while (uart_fifo_read(uart_dev, &c, 1) > 0) {
		continue;
	}
}

static uint8_t rd03d_checksum(const uint8_t *data)
{
	uint8_t cs = 0;

	for (uint8_t i = 1; i < RD03D_BUF_LEN - 1; i++) {
		cs += data[i];
	}

	return 0xff - cs + 1;
}

static int rd03d_open_cmd_mode(const struct device *dev)
{
	struct rd03d_data *data = dev->data;
	int ret;

	if ((data->operation_mode & RD03D_OPERATION_MODE_CMD) == RD03D_OPERATION_MODE_CMD) {
		return 0;
	}

	// Open command mode
	ret = rd03d_send_cmd(dev, RD03D_PROTOCOL_CMD_IDX_OPEN_CMD_MODE, true);
	if (ret < 0) {
		return ret;
	}

	// Check if the command was acknowledged
	if (data->rd_data[0] != RD03D_PROTOCOL_CMD_ACK_IDX_OPEN_CMD_MODE) {
		LOG_ERR("Failed to open command mode");
		return -EIO;
	}

	data->operation_mode |= RD03D_OPERATION_MODE_CMD;
	return 0;
}

static int rd03d_close_cmd_mode(const struct device *dev)
{
	struct rd03d_data *data = dev->data;
	int ret;

	if ((data->operation_mode & RD03D_OPERATION_MODE_CMD) == 0) {
		return 0;
	}

	// Close command mode
	data->operation_mode &= ~RD03D_OPERATION_MODE_CMD;

	// Close command mode
	ret = rd03d_send_cmd(dev, RD03D_PROTOCOL_CMD_IDX_CLOSE_CMD_MODE, true);
	if (ret < 0) {
		return ret;
	}

	// Check if the command was acknowledged
	if (data->rd_data[0] != RD03D_PROTOCOL_CMD_ACK_IDX_CLOSE_CMD_MODE) {
		LOG_ERR("Failed to close command mode");
		return -EIO;
	}

	return 0;
}

static int rd03d_send_cmd(const struct device *dev, enum rd03d_cmd_idx cmd_idx, int value,
			  bool has_rsp)
{
	struct rd03d_data *data = dev->data;
	const struct rd03d_cfg *cfg = dev->config;
	int ret;

	if (data->operation_mode != RD03D_OPERATION_MODE_CMD) {
		return -EPERM;
	}

	/* Make sure last command has been transferred */
	ret = k_sem_take(&data->tx_sem, RD03D_WAIT);
	if (ret) {
		return ret;
	}

	data->tx_data_len = prepare_cmd(cmd_idx, data->tx_data, value);
	data->cmd_idx = cmd_idx;
	data->has_rsp = has_rsp;
	k_sem_reset(&data->rx_sem);

	uart_irq_tx_enable(cfg->uart_dev);

	if (has_rsp) {
		uart_irq_rx_enable(cfg->uart_dev);
		ret = k_sem_take(&data->rx_sem, RD03D_WAIT);
	}

	return ret;
}

static inline int rd03d_get_attribute(const struct device *dev, enum rd03d_attribute attr)
{
	struct rd03d_data *data = dev->data;
	uint8_t checksum;
	int ret;

	// get attribute cmd, has a reponse
	ret = rd03d_send_cmd(dev, cmd_idx, 0, true);
	if (ret < 0) {
		return ret;
	}

	// Decode the response

	return 0;
}

static int rd03d_channel_get(const struct device *dev, enum sensor_channel chan,
			     struct sensor_value *val)
{
	struct rd03d_data *data = dev->data;

	// TODO, for any custom channels that need to read from the

	if (chan != SENSOR_CHAN_ALL) {
		return -ENOTSUP;
	}

	val->val1 = (int32_t)data->data;
	val->val2 = 0;

	return 0;
}

static int rd03d_set_attribute(const struct device *dev, rd03d_protocol_cmd_idx cmd_idx, int value)
{
	struct rd03d_data *data = dev->data;

	// This is always a mutable command, so we need to copy the command into
	// the transmit (tx) buffer, set the value, and then send the command.

	// Note: The user has to explicitly set the command mode to be able to
	//       set and get attributes/channel data.
	if (data->operation_mode != RD03D_OPERATION_MODE_CMD) {
		rd03d_open_cmd_mode(dev);
	}

	// Set the attribute value in the command
	rd03d_send_cmd(dev, cmd_idx, value, true);

	// Decode the response

	// Note: When setting RD03D_ATTR_OPERATION_MODE to anything other than
	//       RD03D_OPERATION_MODE_CMD, the sensor will close the command mode
	//       and enter the new mode. This means that the command mode will
	//       be closed and the sensor will not respond to any commands until
	if (attr == RD03D_ATTR_OPERATION_MODE && value != RD03D_OPERATION_MODE_CMD) {
		rd03d_close_cmd_mode(dev);
	}
}

static int rd03d_attr_set(const struct device *dev, enum sensor_channel chan,
			  enum sensor_attribute attr, const struct sensor_value *val)
{
	if (!(chan >= SENSOR_CHAN_RD03D_MIN_DISTANCE && chan <= SENSOR_CHAN_RD03D_OPERATION_MODE)) {
		return -ENOTSUP;
	}

	switch (chan) {
	case SENSOR_CHAN_RD03D_MIN_DISTANCE:
		return rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_SET_MIN_DISTANCE, val->val1);
	case SENSOR_CHAN_RD03D_MAX_DISTANCE:
		return rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_SET_MAX_DISTANCE, val->val1);
	case SENSOR_CHAN_RD03D_MIN_FRAMES:
		return rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_SET_MIN_FRAMES, val->val1);
	case SENSOR_CHAN_RD03D_MAX_FRAMES:
		return rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_SET_MAX_FRAMES, val->val1);
	case SENSOR_CHAN_RD03D_DELAY_TIME:
		return rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_SET_DELAY_TIME, val->val1);
	case SENSOR_CHAN_RD03D_DETECTION_MODE:
		switch (val->val1) {
		case RD03D_DETECTION_MODE_SINGLE_TARGET:
			return rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_SINGLE_TARGET_MODE,
						   0);
		case RD03D_DETECTION_MODE_MULTI_TARGET:
			return rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_MULTI_TARGET_MODE,
						   0);
		}
		return -EINVAL;
	case SENSOR_CHAN_RD03D_OPERATION_MODE:
		switch (val->val1) {
		case RD03D_OPERATION_MODE_DEBUG:
			return rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_DEBUGGING_MODE, 0);
		case RD03D_OPERATION_MODE_REPORT:
			return rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_REPORTING_MODE, 0);
		case RD03D_OPERATION_MODE_RUN:
			return rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_RUNNING_MODE, 0);
		}
		return -EINVAL;
	}
}

static int rd03d_attr_get(const struct device *dev, enum sensor_channel chan,
			  enum sensor_attribute attr, struct sensor_value *val)
{
	struct rd03d_data *data = dev->data;

	int ret;

	if (chan != SENSOR_CHAN_ALL ||
	    !(chan >= SENSOR_CHAN_RD03D_MIN_DISTANCE && chan <= SENSOR_CHAN_RD03D_OPERATION_MODE) ||
	    !(attr >= SENSOR_CHAN_RD03D_T0_POS && attr <= SENSOR_CHAN_RD03D_T2_POS) ||
	    !(attr >= SENSOR_CHAN_RD03D_T0_SPEED && attr <= SENSOR_CHAN_RD03D_T2_SPEED) ||
	    !(attr >= SENSOR_CHAN_RD03D_T0_DISTANCE && attr <= SENSOR_CHAN_RD03D_T2_DISTANCE)) {
		return -ENOTSUP;
	}

	if (!(chan >= SENSOR_CHAN_RD03D_MIN_DISTANCE && chan <= SENSOR_CHAN_RD03D_OPERATION_MODE)) {
		return -ENOTSUP;
	}

	ret = 0;

	switch (chan) {
	case SENSOR_CHAN_RD03D_MIN_DISTANCE:
		val->val1 = data->min_distance;
		break;
	case SENSOR_CHAN_RD03D_MAX_DISTANCE:
		val->val1 = data->max_distance;
		break;
	case SENSOR_CHAN_RD03D_MIN_FRAMES:
		val->val1 = data->min_frames;
		break;
	case SENSOR_CHAN_RD03D_MAX_FRAMES:
		val->val1 = data->max_frames;
		break;
	case SENSOR_CHAN_RD03D_DELAY_TIME:
		val->val1 = data->delay_time;
		break;
	case SENSOR_CHAN_RD03D_DETECTION_MODE:
		switch (data->detection_mode) {
		case RD03D_DETECTION_MODE_SINGLE_TARGET:
			val->val1 = RD03D_DETECTION_MODE_SINGLE_TARGET;
			break;
		case RD03D_DETECTION_MODE_MULTI_TARGET:
			val->val1 = RD03D_DETECTION_MODE_MULTI_TARGET;
			break;
		}
		break;
	case SENSOR_CHAN_RD03D_OPERATION_MODE:
		switch (data->operation_mode) {
		case RD03D_OPERATION_MODE_CMD:
			val->val1 = RD03D_OPERATION_MODE_CMD;
			break;
		case RD03D_OPERATION_MODE_DEBUG:
			val->val1 = RD03D_OPERATION_MODE_DEBUG;
			break;
		case RD03D_OPERATION_MODE_REPORT:
			val->val1 = RD03D_OPERATION_MODE_REPORT;
			break;
		case RD03D_OPERATION_MODE_RUN:
			val->val1 = RD03D_OPERATION_MODE_RUN;
			break;
		}
		break;

	case SENSOR_CHAN_RD03D_T0_POS:
	case SENSOR_CHAN_RD03D_T1_POS:
	case SENSOR_CHAN_RD03D_T2_POS: {
		int target_idx = chan - SENSOR_CHAN_RD03D_T0_POS;
		val->val1 = data->targets[target_idx].x;
		val->val2 = data->targets[target_idx].y;
		break;
	}
	case SENSOR_CHAN_RD03D_T0_SPEED:
	case SENSOR_CHAN_RD03D_T1_SPEED:
	case SENSOR_CHAN_RD03D_T2_SPEED: {
		int target_idx = chan - SENSOR_CHAN_RD03D_T0_SPEED;
		val->val1 = data->targets[target_idx].speed;
		val->val2 = 0;
		break;
	}
	case SENSOR_CHAN_RD03D_T0_DISTANCE:
	case SENSOR_CHAN_RD03D_T1_DISTANCE:
	case SENSOR_CHAN_RD03D_T2_DISTANCE: {
		int target_idx = chan - SENSOR_CHAN_RD03D_T0_DISTANCE;
		val->val1 = data->targets[target_idx].distance;
		val->val2 = 0;
		break;
	}
	}

	return ret;
}

static int rd03d_sample_fetch(const struct device *dev, enum sensor_channel chan)
{
	// When in 'reporting' mode, the sensor will send 'reports' very frequently
	// and data will be available in the RX buffer.
	struct rd03d_data *data = dev->data;

	// We decode the rx buffer into data->targets
	if (data->operation_mode == RD03D_OPERATION_MODE_REPORT) {

		uart_irq_rx_enable(cfg->uart_dev);

		ret = k_sem_take(&data->rx_sem, RD03D_WAIT);
		if (ret < 0) {
			return ret;
		}

		// Decode the response
		// RD03D_REPORT_FRAME_BEGIN 0xAA, 0xFF, 0x03, 0x00
		//   Target 1 { x, y, speed, distance }
		//   Target 2 { x, y, speed, distance }
		//   Target 3 { x, y, speed, distance }
		// RD03D_REPORT_FRAME_END   0x55, 0xCC

		uint8_t const *rx = data->rx_data;

		if (rx[0] != 0xAA || rx[1] != 0xFF || rx[2] != 0x03 || rx[3] != 0x00) {
			LOG_ERR("Invalid report frame");
			return -EINVAL;
		}

		// rx len should be 4 + (8 * num_targets) + 2

		// For multi-target, assume each target occupies 8 bytes and parse them
		// sequentially. Here we start from index 4 and step through the buffer.
		int ti = 0;
		for (int i = 4; i < (data->rx_bytes - 2) && ti < RD03D_MAX_TARGETS; i += 8) {

			data->targets[ti].x = (int16_t)(rx[i] | (rx[i + 1] << 8)) - 0x200;
			data->targets[ti].y = (int16_t)(rx[i + 2] | (rx[i + 3] << 8)) - 0x8000;
			data->targets[ti].speed = (int16_t)(rx[i + 4] | (rx[i + 5] << 8)) - 0x10;
			data->targets[ti].distance = (uint16_t)(rx[i + 6] | (rx[i + 7] << 8));
		}
	}

	return -ENOTSUP;
}

static DEVICE_API(sensor, rd03d_api_funcs) = {
	.attr_set = rd03d_attr_set,
	.attr_get = rd03d_attr_get,
	.sample_fetch = rd03d_sample_fetch,
	.channel_get = rd03d_channel_get,
};

static void rd03d_uart_isr(const struct device *uart_dev, void *user_data)
{
	const struct device *dev = user_data;
	struct rd03d_data *data = dev->data;

	ARG_UNUSED(user_data);

	if (uart_dev == NULL) {
		return;
	}

	if (!uart_irq_update(uart_dev)) {
		return;
	}

	if (uart_irq_rx_ready(uart_dev)) {
		data->xfer_bytes += uart_fifo_read(uart_dev, &data->rd_data[data->xfer_bytes],
						   RD03D_BUF_LEN - data->xfer_bytes);

		if (data->xfer_bytes == RD03D_BUF_LEN) {
			data->xfer_bytes = 0;
			uart_irq_rx_disable(uart_dev);
			k_sem_give(&data->rx_sem);
			if (data->has_rsp) {
				k_sem_give(&data->tx_sem);
			}
		}
	}

	if (uart_irq_tx_ready(uart_dev)) {
		data->xfer_bytes +=
			uart_fifo_fill(uart_dev, &rd03d_cmds[data->cmd_idx][data->xfer_bytes],
				       RD03D_BUF_LEN - data->xfer_bytes);

		if (data->xfer_bytes == RD03D_BUF_LEN) {
			data->xfer_bytes = 0;
			uart_irq_tx_disable(uart_dev);
			if (!data->has_rsp) {
				k_sem_give(&data->tx_sem);
			}
		}
	}
}

static int rd03d_init(const struct device *dev)
{
	struct rd03d_data *data = dev->data;
	const struct rd03d_cfg *cfg = dev->config;

	int ret;

	data->tx_data_len = 0;
	memset(data->tx_data, 0, RD03D_TX_BUF_MAX_LEN);
	const uint8_t *header_begin = RD03D_CMD_HEADER_BEGIN;
	data->tx_data[0] = header_begin[0];
	data->tx_data[1] = header_begin[1];
	data->tx_data[2] = header_begin[2];
	data->tx_data[3] = header_begin[3];

	data->rx_data_len = 0;
	memset(data->rx_data, 0, RD03D_RX_BUF_MAX_LEN);

	data->operation_mode = RD03D_OPERATION_MODE_CMD;
	data->detection_mode = RD03D_DETECTION_MODE_MULTI_TARGET;

	uart_irq_rx_disable(cfg->uart_dev);
	uart_irq_tx_disable(cfg->uart_dev);

	rd03d_uart_flush(cfg->uart_dev);

	uart_irq_callback_user_data_set(cfg->uart_dev, cfg->cb, (void *)dev);

	k_sem_init(&data->rx_sem, 0, 1);
	k_sem_init(&data->tx_sem, 1, 1);

	rd03d_open_cmd_mode(dev);

	/* Configure default min and max range */
	ret = rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_SET_MIN_DISTANCE, cfg->min_distance);
	if (ret != 0) {
		LOG_ERR("Error setting minimum range %d", cfg->range);
		return ret;
	}
	ret = rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_SET_MAX_DISTANCE, cfg->max_distance);
	if (ret != 0) {
		LOG_ERR("Error setting maximum range %d", cfg->range);
		return ret;
	}
	/* Configure default min and max frames */
	ret = rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_SET_MIN_FRAMES, cfg->min_frames);
	if (ret != 0) {
		LOG_ERR("Error setting minimum frames %d", cfg->min_frames);
		return ret;
	}
	ret = rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_SET_MAX_FRAMES, cfg->max_frames);
	if (ret != 0) {
		LOG_ERR("Error setting maximum frames %d", cfg->max_frames);
		return ret;
	}
	/* Configure default delay time */
	ret = rd03d_set_attribute(dev, RD03D_PROTOCOL_CMD_IDX_SET_DELAY_TIME, cfg->delay_time);
	if (ret != 0) {
		LOG_ERR("Error setting delay time %d", cfg->delay_time);
		return ret;
	}

	/* Configure operation mode */
	ret = rd03d_set_attribute(dev, RD03D_ATTR_OPERATION_MODE, RD03D_OPERATION_MODE_REPORTING);
	if (ret != 0) {
		LOG_ERR("Error setting default operation mode");
	}

	return ret;
}

#define RD03D_INIT(inst)                                                                           \
                                                                                                   \
	static struct rd03d_data rd03d_data_##inst;                                                \
                                                                                                   \
	static const struct rd03d_cfg rd03d_cfg_##inst = {                                         \
		.uart_dev = DEVICE_DT_GET(DT_INST_BUS(inst)),                                      \
		.min_distance = DT_INST_PROP(inst, minimum_distance),                              \
		.max_distance = DT_INST_PROP(inst, maximum_distance),                              \
		.min_frames = DT_INST_PROP(inst, minimum_frames),                                  \
		.max_frames = DT_INST_PROP(inst, maximum_frames),                                  \
		.delay_time = DT_INST_PROP(inst, delay_time),                                      \
		.cb = rd03d_uart_isr,                                                              \
	};                                                                                         \
                                                                                                   \
	SENSOR_DEVICE_DT_INST_DEFINE(inst, rd03d_init, NULL, &rd03d_data_##inst,                   \
				     &rd03d_cfg_##inst, POST_KERNEL, CONFIG_SENSOR_INIT_PRIORITY,  \
				     &rd03d_api_funcs);

DT_INST_FOREACH_STATUS_OKAY(RD03D_INIT)
