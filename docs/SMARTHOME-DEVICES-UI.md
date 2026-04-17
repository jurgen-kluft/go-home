# SmartHome Devices UI

In our smart home we will have a server running on a Mac Mini M4 that can render multiple UI instances using ImGui and Glfw. Every device will connect to the server through TCP and will so have its own UI instance. The UI screen will be rendered and then the current and previous frames will be diffed and dirty blocks will be identified and sent over the TCP connection, which will then update the display on the device. This allows for a very flexible and customizable UI for every device, as well as a very efficient way to render the UI on the device. It will also avoid iterating on the UI directly on the device, which can be very time consuming and inefficient, especially for devices with limited resources and long iteration times.

- ImGUI + Glfw
- Capture UI instance frame-buffer on the server side
- Diff and send dirty blocks to the device to update the display

## Threads

1. Server thread, rendering the UI instances and identifying dirty blocks to send to the devices
2. Network thread, managing the TCP connections and sending/receiving data to/from the devices

## Messages

- Upon TCP connection, the device will send a message to the server with its device type and capabilities, so that the server can create the appropriate UI instance for that device.
- Receive network messages from devices, e.g. button presses, touch events, etc.