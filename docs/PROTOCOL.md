# Protocol Specification

The protocol is designed to be both simple and efficient.

It allows for flexible data formats and can be easily extended to support new sensor types in the future. It is also a performant protocol, as it allows for direct memory copying of the data on the receiving end, eliminating the need for additional parsing or processing.

This is `key-value` based, where the key is a uint16_t ID (2 bytes) + MAC address (6 bytes) and the value is a variable-length binary blob. The ID identifies the type of data being sent, and the MAC address identifies the location, while the value contains the actual data.

The `[ID, MAC]` pair can be used to directly find a memory mapped file on the receiving end, allowing for efficient data transfer without the need for additional parsing or processing.

The size of the data to write to the memory mapped file is determined by the `length` field, which specifies the length of the value in bytes. This allows for variable-length data to be sent without the need for a fixed-size buffer.

A network packet can contain multiple `msg_t` structures, allowing for multiple sensor readings to be sent in a single packet. 

```c
struct msg_t
{
    uint16_t key;       // e.g. ID_TEMPERATURE
    uint8_t  mac[6];    // source MAC address (location)
    uint16_t length;    // length of data (e.g. 2 bytes for int16_t temperature)
    uint8_t  value[];   // e.g. int16_t temperature in degrees Celsius
};

```


