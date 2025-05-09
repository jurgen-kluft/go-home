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

// 	RD03D_CMD_IDX_OPEN_CMD_MODE      = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0xFF, 0x00, 0x01, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_CLOSE_CMD_MODE     = { RD03D_CMD_HEADER_BEGIN, 0x02, 0x00, 0xFE, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_DEBUGGING_MODE     = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_REPORTING_MODE     = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_RUNNING_MODE       = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_SET_MIN_DISTANCE   = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_SET_MAX_DISTANCE   = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END },
// 	RD03D_CMD_IDX_SET_MIN_FRAMES     = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x07, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_SET_MAX_FRAMES     = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x07, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_SET_DELAY_TIME     = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x07, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_GET_MIN_DISTANCE   = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0x08, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_GET_MAX_DISTANCE   = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0x08, 0x00, 0x01, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_GET_MIN_FRAMES     = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0x08, 0x00, 0x02, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_GET_MAX_FRAMES     = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0x08, 0x00, 0x03, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_GET_DELAY_TIME     = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0x08, 0x00, 0x04, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_SINGLE_TARGET_MODE = { RD03D_CMD_HEADER_BEGIN, 0x02, 0x00, 0x80, 0x00, RD03D_CMD_HEADER_END }
// 	RD03D_CMD_IDX_MULTI_TARGET_MODE  = { RD03D_CMD_HEADER_BEGIN, 0x02, 0x00, 0x90, 0x00, RD03D_CMD_HEADER_END },

// return the length of the command
static int prepare_cmd(rd03d_protocol_cmd_idx, uint8_t *cmd_buffer, int value) {
	// Assume the header is already set
	int l = 4; // Skip header

	if (cmd_idx  == RD03D_CMD_IDX_OPEN_CMD_MODE) 
	{
		cmd_buffer[l++] = 0x04; // Intra-frame data length
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0xFF; // Command word
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x01; // 
		cmd_buffer[l++] = 0x00; //
	}
	else if (cmd_idx == RD03D_CMD_IDX_CLOSE_CMD_MODE) 
	{
		cmd_buffer[l++] = 0x02; // Intra-frame data length
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0xFE; // Command word
		cmd_buffer[l++] = 0x00; //
	} 
	else if (cmd_idx >= RD03D_CMD_IDX_DEBUGGING_MODE && cmd_idx <= RD03D_CMD_IDX_RUNNING_MODE) 
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
	else if (cmd_idx >= RD03D_CMD_IDX_SET_MIN_DISTANCE && cmd_idx <= RD03D_CMD_IDX_SET_DELAY_TIME) 
	{
		cmd_buffer[l++] = 0x08; // Intra-frame data length
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x07; // Command word
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = cmd_index - RD03D_CMD_IDX_SET_MIN_DISTANCE; // Command value
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = value & 0xFF; // Set value, 32-bit
		cmd_buffer[l++] = (value >> 8) & 0xFF;
		cmd_buffer[l++] = (value >> 16) & 0xFF;
		cmd_buffer[l++] = (value >> 24) & 0xFF;
	}
	else if (cmd_idx >= RD03D_CMD_IDX_GET_MIN_DISTANCE && cmd_idx <= RD03D_CMD_IDX_GET_DELAY_TIME) 
	{
		cmd_buffer[l++] = 0x04; // Intra-frame data length
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = 0x08; // Command word
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = cmd_index - RD03D_CMD_IDX_GET_MIN_DISTANCE; // Command value
		cmd_buffer[l++] = 0x00; //
	}
	else if (cmd_idx >= RD03D_CMD_IDX_SINGLE_TARGET_MODE && cmd_idx <= RD03D_CMD_IDX_MULTI_TARGET_MODE) 
	{
		cmd_buffer[l++] = 0x02; // Intra-frame data length
		cmd_buffer[l++] = 0x00; //
		cmd_buffer[l++] = RD03D_CMD_IDX_SINGLE_TARGET_MODE ? 0x80 : 0x90; // Command word
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

static const uint8_t RD03D_CMD_ACK_IDX_OPEN_CMD_MODE[]   = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0xFF, 0x01, 0x00, 0x00, 0x01, 0x00, 0x40, 0x00, RD03D_CMD_HEADER_END };
static const uint8_t RD03D_CMD_ACK_IDX_CLOSE_CMD_MODE[]  = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0xFE, 0x01, 0x00, 0x00, RD03D_CMD_HEADER_END };

// Protocol, ACKS we get for set and get commands, the ACK related to the get cmd contains a 4 byte variable at (rx_buffer[10] to rx_buffer[13])
// static const uint8_t RD03D_CMD_ACK_IDX_SET_XXX[] = { RD03D_CMD_HEADER_BEGIN, 0x04, 0x00, 0x07, 0x01, 0x00, 0x00, RD03D_CMD_HEADER_END },
// static const uint8_t RD03D_CMD_ACK_IDX_GET_XXX[] = { RD03D_CMD_HEADER_BEGIN, 0x08, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, RD03D_CMD_HEADER_END },

// clang-format on

static int verify_open_cmd_ack(const uint8_t *data, int data_len)
{
	if (data_len != sizeof(RD03D_CMD_ACK_IDX_OPEN_CMD_MODE)) {
		return -1;
	}

	return memcmp(data, RD03D_CMD_ACK_IDX_OPEN_CMD_MODE,
		      sizeof(RD03D_CMD_ACK_IDX_OPEN_CMD_MODE)) == 0
		       ? 0
		       : -1;
}

static int verify_close_cmd_ack(const uint8_t *data, int data_len)
{
	if (data_len != sizeof(RD03D_CMD_ACK_IDX_CLOSE_CMD_MODE)) {
		return -1;
	}

	return memcmp(data, RD03D_CMD_ACK_IDX_CLOSE_CMD_MODE,
		      sizeof(RD03D_CMD_ACK_IDX_CLOSE_CMD_MODE)) == 0
		       ? 0
		       : -1;
}

static int verify_set_cmd_ack(const uint8_t *data, int data_len)
{
	if (data_len != sizeof(RD03D_CMD_ACK_IDX_SET_XXX)) {
		return -1;
	}

	return memcmp(data, RD03D_CMD_ACK_IDX_SET_XXX, sizeof(RD03D_CMD_ACK_IDX_SET_XXX)) == 0 ? 0
											       : -1;
}

static int verify_get_cmd_ack(const uint8_t *data, int data_len, int *value)
{
	*value = 0;
	if (data_len != sizeof(RD03D_CMD_ACK_IDX_GET_XXX)) {
		return -1;
	}

	if (memcmp(data, RD03D_CMD_ACK_IDX_GET_XXX, sizeof(RD03D_CMD_ACK_IDX_GET_XXX)) == 0) {
		// read the 4 byte variable
		*value = data[10] | (data[11] << 8) | (data[12] << 16) | (data[13] << 24);
		return 0;
	}

	return -1;
}

static void rd03d_uart_flush(const struct device *uart_dev)
{
	uint8_t c;

	while (uart_fifo_read(uart_dev, &c, 1) > 0) {
		continue;
	}
}

static int rd03d_send_cmd(const struct device *dev, enum rd03d_cmd_idx cmd_idx, int value)
{
	struct rd03d_data *data = dev->data;
	const struct rd03d_cfg *cfg = dev->config;
	int ret;

	if (data->operation_mode != RD03D_OPERATION_MODE_CMD) {
		return -EPERM;
	}

	/* Make sure last command has been transferred */
	ret = k_sem_take(&data->tx_sem, RD03D_SEMA_MAX_WAIT);
	if (ret) {
		return ret;
	}

	data->tx_data_len = prepare_cmd(cmd_idx, data->tx_data, value);

	k_sem_reset(&data->rx_sem);

	// all the rd03d commands have a response
	uart_irq_tx_enable(cfg->uart_dev);
	uart_irq_rx_enable(cfg->uart_dev);

	// wait for the tx and rx to be done
	ret = k_sem_take(&data->rx_sem, RD03D_SEMA_MAX_WAIT);

	return ret;
}

static int rd03d_open_cmd_mode(const struct device *dev)
{
	struct rd03d_data *data = dev->data;
	int ret;

	// Open command mode
	ret = rd03d_send_cmd(dev, RD03D_CMD_IDX_OPEN_CMD_MODE, 1);
	if (ret < 0) {
		return ret;
	}

	// Check if the command was acknowledged
	if (verify_open_cmd_ack(data->rx_data, data->rx_data_len) != 0) {
		LOG_ERR("Failed to open command mode");
		return -EIO;
	}

	return 0;
}

static int rd03d_close_cmd_mode(const struct device *dev)
{
	struct rd03d_data *data = dev->data;
	int ret;

	// Close command mode
	ret = rd03d_send_cmd(dev, RD03D_CMD_IDX_CLOSE_CMD_MODE, 1);
	if (ret < 0) {
		return ret;
	}

	// Verify the command was acknowledged successfully
	if (verify_close_cmd_ack(data->rx_data, data->rx_data_len) != 0) {
		LOG_ERR("Failed to close command mode");
		return -EIO;
	}

	return 0;
}

static int rd03d_set_attribute(const struct device *dev, enum rd03d_protocol_cmd_idx cmd_idx,
			       int value)
{
	struct rd03d_data *data = dev->data;

	int ret;

	// This is always a mutable command, so we need to copy the command into
	// the transmit (tx) buffer, set the value, and then send the command.

	// Note: The user has to explicitly set the command mode to be able to
	//       set and get attributes/channel data.
	if ((data->operation_mode & RD03D_OPERATION_MODE_CMD) == 0) {
		ret = rd03d_open_cmd_mode(dev);
		if (ret < 0) {
			LOG_ERR("Failure, open command mode");
			return -EINVAL;
		}
		data->operation_mode |= RD03D_OPERATION_MODE_CMD;
	}

	// Set the attribute value in the command
	ret = rd03d_send_cmd(dev, cmd_idx, value, 1) if (ret < 0)
	{
		LOG_ERR("Failure, set attribute command (%d)", cmd_idx);
		return ret;
	}

	// Verify the command was acknowledged successfully
	ret = verify_set_cmd_ack(data->rx_data, data->rx_data_len);
	if (ret < 0) {
		LOG_ERR("Failure, set attribute command (%d) did not get valid ACK", cmd_idx);
		return ret;
	}

	// Note: When setting RD03D_ATTR_OPERATION_MODE to any of the reporting
	//       modes, the sensor will close the command mode.
	//       This means that the command mode will be closed automatically
	//       and the sensor will enter the reporting mode.
	if ((cmd_idx >= RD03D_CMD_IDX_DEBUGGING_MODE && cmd_idx <= RD03D_CMD_IDX_RUNNING_MODE)) {
		if ((data->operation_mode & RD03D_OPERATION_MODE_CMD) == RD03D_OPERATION_MODE_CMD) {
			ret = rd03d_close_cmd_mode(dev);
			if (ret < 0) {
				return -EINVAL;
			}
			data->operation_mode &= ~RD03D_OPERATION_MODE_CMD;
		}
	}

	return ret;
}

/*
enum sensor_channel_rd03d {
	SENSOR_CHAN_RD03D_CONFIG_DISTANCE = SENSOR_CHAN_PRIV_START,
	SENSOR_CHAN_RD03D_CONFIG_FRAMES,
	SENSOR_CHAN_RD03D_CONFIG_DELAY_TIME,
	SENSOR_CHAN_RD03D_CONFIG_DETECTION_MODE,
	SENSOR_CHAN_RD03D_CONFIG_OPERATION_MODE,
	SENSOR_CHAN_RD03D_POS,
	SENSOR_CHAN_RD03D_SPEED,
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

*/

static int rd03d_attr_set(const struct device *dev, enum sensor_channel chan,
			  enum sensor_attribute attr, const struct sensor_value *val)
{
	struct rd03d_data *data = dev->data;
	struct rd03d_cfg *cfg = dev->config;
	int ret;

	if (!(chan >= SENSOR_CHAN_RD03D_CONFIG_DISTANCE &&
	      chan <= SENSOR_CHAN_RD03D_CONFIG_OPERATION_MODE)) {
		return -ENOTSUP;
	}
	if (!(chan >= SENSOR_ATTR_RD03D_CONFIG_VALUE && chan <= SENSOR_ATTR_RD03D_CONFIG_MAXIMUM)) {
		return -ENOTSUP;
	}

	switch (chan) {
	case SENSOR_CHAN_RD03D_CONFIG_DISTANCE:
		switch (attr) {
		case SENSOR_ATTR_RD03D_CONFIG_MINIMUM:
			if (rd03d_set_attribute(dev, RD03D_CMD_IDX_SET_MIN_DISTANCE, val->val1) >=
			    0) {
				cfg->min_distance = val->val1;
			}
			break;
		case SENSOR_ATTR_RD03D_CONFIG_MAXIMUM:
			if (rd03d_set_attribute(dev, RD03D_CMD_IDX_SET_MAX_DISTANCE, val->val1) >=
			    0) {
				cfg->max_distance = val->val1;
			}
			break;
		}
		break;
	case SENSOR_CHAN_RD03D_CONFIG_FRAMES:
		switch (attr) {
		case SENSOR_ATTR_RD03D_CONFIG_MINIMUM:
			if (rd03d_set_attribute(dev, RD03D_CMD_IDX_SET_MIN_FRAMES, val->val1) >=
			    0) {
				cfg->min_frames = val->val1;
			}
			break;
		case SENSOR_ATTR_RD03D_CONFIG_MAXIMUM:
			if (rd03d_set_attribute(dev, RD03D_CMD_IDX_SET_MAX_FRAMES, val->val1) >=
			    0) {
				cfg->max_frames = val->val1;
			}
			break;
		}
		break;
	case SENSOR_CHAN_RD03D_CONFIG_DELAY_TIME:
		if (rd03d_set_attribute(dev, RD03D_CMD_IDX_SET_DELAY_TIME, val->val1) >= 0) {
			cfg->delay_time = val->val1;
		}
		break;
	case SENSOR_CHAN_RD03D_DETECTION_MODE:
		switch (val->val1) {
		case RD03D_DETECTION_MODE_SINGLE_TARGET:
			if (rd03d_set_attribute(dev, RD03D_CMD_IDX_SINGLE_TARGET_MODE, val->val1) >=
			    0) {
				data->detection_mode = RD03D_DETECTION_MODE_SINGLE_TARGET;
			}
			break;
		case RD03D_DETECTION_MODE_MULTI_TARGET:
			if (rd03d_set_attribute(dev, RD03D_CMD_IDX_MULTI_TARGET_MODE, val->val1) >=
			    0) {
				data->detection_mode = RD03D_DETECTION_MODE_MULTI_TARGET;
			}
			break;
		default:
			ret = -EINVAL;
			break;
		}
		break;
	case SENSOR_CHAN_RD03D_OPERATION_MODE:
		switch (val->val1) {
		case RD03D_OPERATION_MODE_DEBUG:
			if (rd03d_set_attribute(dev, RD03D_CMD_IDX_DEBUGGING_MODE, val->val1) >=
			    0) {
				data->operation_mode = RD03D_OPERATION_MODE_DEBUG;
			}
			break;
		case RD03D_OPERATION_MODE_REPORT:
			if (rd03d_set_attribute(dev, RD03D_CMD_IDX_REPORTING_MODE, val->val1) >=
			    0) {
				data->operation_mode = RD03D_OPERATION_MODE_REPORT;
			}
			break;
		case RD03D_OPERATION_MODE_RUN:
			if (rd03d_set_attribute(dev, RD03D_CMD_IDX_RUNNING_MODE, val->val1) >= 0) {
				data->operation_mode = RD03D_OPERATION_MODE_RUN;
			}
			break;
		default:
			ret = -EINVAL;
			break;
		}
		break;
	}

	return ret;
}

static inline int rd03d_get_attribute(const struct device *dev, rd03d_protocol_cmd_idx cmd_idx,
				      int *value)
{
	struct rd03d_data *data = dev->data;
	int ret;

	// get attribute cmd, has a reponse
	ret = rd03d_send_cmd(dev, cmd_idx, 0);
	if (ret < 0) {
		return ret;
	}

	// Decode the response and write 'value'

	return 0;
}

static int rd03d_attr_get(const struct device *dev, enum sensor_channel chan,
			  enum sensor_attribute attr, struct sensor_value *val)
{
	struct rd03d_data *data = dev->data;
	int ret;

	if (!(chan >= SENSOR_CHAN_RD03D_CONFIG_DISTANCE &&
	      chan <= SENSOR_CHAN_RD03D_CONFIG_OPERATION_MODE)) {
		return -ENOTSUP;
	}
	if (!(chan >= SENSOR_ATTR_RD03D_CONFIG_VALUE && chan <= SENSOR_ATTR_RD03D_CONFIG_MAXIMUM)) {
		return -ENOTSUP;
	}

	int ti = 0;

	switch (chan) {
	case SENSOR_CHAN_RD03D_POS:
		ti = attr - SENSOR_ATTR_RD03D_TARGET_0;
		if (ti < 0 || ti >= RD03D_MAX_TARGETS) {
			return -EINVAL;
		}
		val->val1 = data->targets[ti].x;
		val->val2 = data->targets[ti].y;
		break;
	case SENSOR_CHAN_RD03D_SPEED:
		ti = attr - SENSOR_ATTR_RD03D_TARGET_0;
		if (ti < 0 || ti >= RD03D_MAX_TARGETS) {
			return -EINVAL;
		}
		val->val1 = data->targets[ti].speed;
		val->val2 = 0;
		break;
	case SENSOR_CHAN_RD03D_DISTANCE:
		ti = attr - SENSOR_ATTR_RD03D_TARGET_0;
		if (ti < 0 || ti >= RD03D_MAX_TARGETS) {
			return -EINVAL;
		}
		val->val1 = data->targets[ti].distance;
		val->val2 = 0;
		break;

	case SENSOR_CHAN_PROX:
		val->val1 = 0;
		for (int ti = 0; ti < RD03D_MAX_TARGETS; ti++) {
			if (data->targets[ti].x != 0 && data->targets[ti].y != 0) {
				val->val1 |= (1 << ti);
			}
		}
		val->val2 = 0;
		break;
	case SENSOR_CHAN_DISTANCE:
		val->val1 = data->targets[target].distance;
		val->val2 = 0;
		break;

	case SENSOR_CHAN_RD03D_CONFIG_DISTANCE:
		switch (attr) {
		case SENSOR_ATTR_RD03D_CONFIG_MINIMUM:
			ret = rd03d_get_attribute(dev, RD03D_CMD_IDX_GET_MIN_DISTANCE, &val->val1);
			if (ret >= 0) {
				data->min_distance = val->val1;
			}
			break;
		case SENSOR_ATTR_RD03D_CONFIG_MAXIMUM:
			ret = rd03d_get_attribute(dev, RD03D_CMD_IDX_GET_MAX_DISTANCE, &val->val1);
			if (ret >= 0) {
				data->max_distance = val->val1;
			}
			break;
		}
		break;
	case SENSOR_CHAN_RD03D_CONFIG_FRAMES:
		switch (attr) {
		case SENSOR_ATTR_RD03D_CONFIG_MINIMUM:
			ret = rd03d_get_attribute(dev, RD03D_CMD_IDX_GET_MIN_FRAMES, &val->val1);
			if (ret >= 0) {
				data->min_frames = val->val1;
			}
			break;
		case SENSOR_ATTR_RD03D_CONFIG_MAXIMUM:
			ret = rd03d_get_attribute(dev, RD03D_CMD_IDX_GET_MAX_FRAMES, &val->val1);
			if (ret >= 0) {
				data->max_frames = val->val1;
			}
			break;
		}
		break;
	case SENSOR_CHAN_RD03D_CONFIG_DELAY_TIME:
		ret = rd03d_get_attribute(dev, RD03D_CMD_IDX_GET_DELAY_TIME, &val->val1);
		if (ret >= 0) {
			data->delay_time = val->val1;
		}
		break;
	case SENSOR_CHAN_RD03D_DETECTION_MODE:
		if (data->detection_mode == RD03D_DETECTION_MODE_SINGLE_TARGET) {
			val->val1 = RD03D_DETECTION_MODE_SINGLE_TARGET;
		} else {
			val->val1 = RD03D_DETECTION_MODE_MULTI_TARGET;
		}
		break;
	case SENSOR_CHAN_RD03D_OPERATION_MODE:
		switch (data->operation_mode & ~RD03D_OPERATION_MODE_CMD) {
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
		if ((data->operation_mode & RD03D_OPERATION_MODE_CMD) == RD03D_OPERATION_MODE_CMD) {
			val->val2 = RD03D_OPERATION_MODE_CMD;
		}
		break;
	}

	return ret;
}

static int rd03d_channel_get(const struct device *dev, enum sensor_channel chan,
			     struct sensor_value *val)
{
	struct rd03d_data *data = dev->data;
	struct rd03d_cfg *cfg = dev->config;

	int ret;
	const int target = val->val1;

	if (chan >= SENSOR_CHAN_RD03D_POS && chan <= SENSOR_CHAN_RD03D_DISTANCE) {
		if (target < 0 || target >= RD03D_MAX_TARGETS) {
			return -EINVAL;
		}
	} else if (chan == SENSOR_CHAN_PROX || chan == SENSOR_CHAN_DISTANCE) {
		if (target < 0 || target >= RD03D_MAX_TARGETS) {
			return -EINVAL;
		}
	}

	switch (chan) {
	case SENSOR_CHAN_RD03D_POS:
		val->val1 = data->targets[target].x;
		val->val2 = data->targets[target].y;
		break;
	case SENSOR_CHAN_RD03D_SPEED:
		val->val1 = data->targets[target].speed;
		break;
	case SENSOR_CHAN_RD03D_DISTANCE:
		val->val1 = data->targets[target].distance;
		break;
	case SENSOR_CHAN_RD03D_CONFIG_DISTANCE:
		val->val1 = cfg->min_distance;
		val->val2 = cfg->max_distance;
		break;

	case SENSOR_CHAN_PROX:
		val->val1 = 0;
		for (int ti = 0; ti < RD03D_MAX_TARGETS; ti++) {
			if (data->targets[ti].x != 0 && data->targets[ti].y != 0) {
				val->val1 |= (1 << ti);
			}
		}
		val->val2 = 0;
		break;
	case SENSOR_CHAN_DISTANCE:
		val->val1 = data->targets[target].distance;
		val->val2 = 0;
		break;

	case SENSOR_CHAN_RD03D_CONFIG_FRAMES:
		val->val1 = cfg->min_frames;
		val->val2 = cfg->max_frames;
		break;
	case SENSOR_CHAN_RD03D_CONFIG_DELAY_TIME:
		val->val1 = cfg->delay_time;
		val->val2 = 0;
		break;
	case SENSOR_CHAN_RD03D_CONFIG_DETECTION_MODE:
		if (data->detection_mode == RD03D_DETECTION_MODE_SINGLE_TARGET) {
			val->val1 = RD03D_DETECTION_MODE_SINGLE_TARGET;
		} else {
			val->val1 = RD03D_DETECTION_MODE_MULTI_TARGET;
		}
		val->val2 = 0;
		break;
	case SENSOR_CHAN_RD03D_CONFIG_OPERATION_MODE:
		switch (data->operation_mode & ~RD03D_OPERATION_MODE_CMD) {
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
		if ((data->operation_mode & RD03D_OPERATION_MODE_CMD) == RD03D_OPERATION_MODE_CMD) {
			val->val2 = RD03D_OPERATION_MODE_CMD;
		}
		break;
	default:
		LOG_ERR("Unsupported channel %d", chan);
		return -ENOTSUP;
	}

	return ret;
}

static int rd03d_sample_fetch(const struct device *dev, enum sensor_channel chan)
{
	// When in 'reporting' mode, the sensor will send 'reports' continuously
	// and data will become available in the RX buffer.
	struct rd03d_data *data = dev->data;

	// We decode the rx buffer into data->targets
	if (data->operation_mode == RD03D_OPERATION_MODE_REPORT) {

		uart_irq_rx_enable(cfg->uart_dev);

		ret = k_sem_take(&data->rx_sem, RD03D_SEMA_MAX_WAIT);
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

	uint8_t *rxb = &data->rd_data[0];
	if (uart_irq_rx_ready(uart_dev)) {

		const int byreq = RD03D_RX_BUF_MAX_LEN - data->rx_bytes; /* Avoid buffer overrun */
		const int byread = uart_fifo_read(uart_dev, rxb[data->rx_bytes], byreq);

		data->rx_bytes += byread;

determine_rx_data_len:

		/* The minimum data frame length is 14 bytes, and the maximum
		   data frame length is a report which is 30 bytes.
		   Our receive buffer has a size of 64 bytes, so we should be
		   able to receive a full data frame within the buffer, if
		   not then something is incorrect regarding the protocol.

		   The command + ACK protocol should not pose any out-of-sync
		   issues, as the ACK always follows the CMD.

		   The report stream is a bit more tricky, as the sensor
		   continuously sends frame data, and we need to be able to
		   detect the start of a new report frame.
		   The report frame starts with 0xAA, 0xFF, 0x03, 0x00 and
		   ends with 0x55, 0xCC. So we might start receiving a report
		   frame in the middle of a report frame, which we should
		   ignore and continue to receive until we find the start
		   of a new report frame.
		*/
		if (data->rx_data_len == 0 && (data->rx_bytes >= (4 + 2 + 4 + 4))) {
			if (rxb[0] == 0xFD && rxb[1] == 0xFC && rxb[2] == 0xFB && rxb[3] == 0xFA) {
				data->rx_data_len = 4 + 2 + rxb[4] + 4;
			} else if (rxb[0] == 0xAA && rxb[1] == 0xFF && rxb[2] == 0x03 &&
				   rxb[3] == 0x00) {
				data->rx_data_len = 30;
			} else {
				// Scan rx-buffer until 'FD FC FB FA' or 'AA FF 03 00'
				int i = 0;
				for (i = 0; i < data->rx_bytes - 4; i++) {
					if (rxb[i] == 0xFD && rxb[i + 1] == 0xFC &&
					    rxb[i + 2] == 0xFB && rxb[i + 3] == 0xFA) {
						data->rx_bytes -= i;
						// TODO we might be able to avoid the memmove, by
						// introducing data->rx_offset which indicates the
						// start of the data frame in the buffer.
						memmove(rxb, &rxb[i], data->rx_bytes);
						goto determine_rx_data_len;
						break;
					} else if (rxb[i] == 0xAA && rxb[i + 1] == 0xFF &&
						   rxb[i + 2] == 0x03 && rxb[i + 3] == 0x00) {
						data->rx_bytes -= i;
						// TODO we might be able to avoid the memmove, by
						// introducing data->rx_offset which indicates the
						// start of the data frame in the buffer.
						memmove(rxb, &rxb[i], data->rx_bytes);
						goto determine_rx_data_len;
						break;
					}
				}

				// TODO How to recover from this, reset rx_bytes ?
				data->LOG_ERR("Critical: invalid response!");
				data->rx_bytes = 0;
			}
		}

		/* keep reading until the end of the message */
		if (data->rx_bytes == data->rx_data_len) {
			data->rx_bytes = 0;
			uart_irq_rx_disable(uart_dev);
			k_sem_give(&data->rx_sem);
			if (rxb[0] == 0xFD && rxb[1] == 0xFC) {
				/* Receiving an ACK, this means a command was send. Signal
				 */
				/* that the command + ACK is done. */
				k_sem_give(&data->tx_sem);
			}
		}
	}

	if (uart_irq_tx_ready(uart_dev)) {
		data->tx_bytes += uart_fifo_fill(uart_dev, &data->tx_data[data->tx_bytes],
						 data->tx_data_len - data->tx_bytes);

		if (data->tx_bytes == data->tx_data_len) {
			data->tx_bytes = 0;
			uart_irq_tx_disable(uart_dev);
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

	for (int i = 0; i < RD03D_MAX_TARGETS; i++) {
		data->targets[i].x = 0;
		data->targets[i].y = 0;
		data->targets[i].distance = 0;
		data->targets[i].speed = 0;
	}

	uart_irq_rx_disable(cfg->uart_dev);
	uart_irq_tx_disable(cfg->uart_dev);

	rd03d_uart_flush(cfg->uart_dev);

	uart_irq_callback_user_data_set(cfg->uart_dev, cfg->cb, (void *)dev);

	k_sem_init(&data->rx_sem, 0, 1);
	k_sem_init(&data->tx_sem, 1, 1);

	rd03d_open_cmd_mode(dev);

	bool read_config = true;

	if (read_config) {
		int value = 0;

		/* Configure default min and max range */
		if (rd03d_get_attribute(dev, RD03D_CMD_IDX_GET_MIN_DISTANCE, &value) < 0) {
			LOG_ERR("Error getting minimum range %d", cfg->range);
			return ret;
		}
		cfg->min_distance = (uint16_t)value;

		if (rd03d_get_attribute(dev, RD03D_CMD_IDX_GET_MAX_DISTANCE, &value) < 0) {
			LOG_ERR("Error getting maximum range %d", cfg->range);
			return ret;
		}
		cfg->max_distance = (uint16_t)value;

		/* Configure default min and max frames */
		if (rd03d_get_attribute(dev, RD03D_CMD_IDX_GET_MIN_FRAMES, &value) < 0) {
			LOG_ERR("Error getting minimum frames %d", cfg->min_frames);
			return ret;
		}
		cfg->min_frames = (uint16_t)value;

		if (rd03d_get_attribute(dev, RD03D_CMD_IDX_GET_MAX_FRAMES, &value) < 0) {
			LOG_ERR("Error getting maximum frames %d", cfg->max_frames);
			return ret;
		}
		cfg->max_frames = (uint16_t)value;

		/* Configure default delay time */
		if (rd03d_get_attribute(dev, RD03D_CMD_IDX_GET_DELAY_TIME, &value) < 0) {
			LOG_ERR("Error getting delay time %d", cfg->delay_time);
			return ret;
		}
		cfg->delay_time = (uint16_t)value;

	} else {
		/* Configure default min and max range */
		if (rd03d_set_attribute(dev, RD03D_CMD_IDX_SET_MIN_DISTANCE, cfg->min_distance) <
		    0) {
			LOG_ERR("Error setting minimum range %d", cfg->range);
		}
		if (rd03d_set_attribute(dev, RD03D_CMD_IDX_SET_MAX_DISTANCE, cfg->max_distance) <
		    0) {
			LOG_ERR("Error setting maximum range %d", cfg->range);
			return ret;
		}
		/* Configure default min and max frames */
		if (rd03d_set_attribute(dev, RD03D_CMD_IDX_SET_MIN_FRAMES, cfg->min_frames) < 0) {
			LOG_ERR("Error setting minimum frames %d", cfg->min_frames);
			return ret;
		}
		if (rd03d_set_attribute(dev, RD03D_CMD_IDX_SET_MAX_FRAMES, cfg->max_frames) < 0) {
			LOG_ERR("Error setting maximum frames %d", cfg->max_frames);
			return ret;
		}
		/* Configure default delay time */
		if (rd03d_set_attribute(dev, RD03D_CMD_IDX_SET_DELAY_TIME, cfg->delay_time) < 0) {
			LOG_ERR("Error setting delay time %d", cfg->delay_time);
			return ret;
		}
	}

	/* Activate report mode */
	if (rd03d_attr_set(dev, SENSOR_CHAN_RD03D_CONFIG_OPERATION_MODE,
			   SENSOR_ATTR_RD03D_CONFIG_VALUE, RD03D_OPERATION_MODE_REPORT) < 0) {
		LOG_ERR("Error setting default operation mode");
	}

	return ret;
}

static DEVICE_API(sensor, rd03d_api_funcs) = {
	.attr_set = rd03d_attr_set,
	.attr_get = rd03d_attr_get,
	.sample_fetch = rd03d_sample_fetch,
	.channel_get = rd03d_channel_get,
};

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
