# Remote Rendered UI

## Constraints

- ESP32 is a display client
- Mac Mini M4 is the host that renders the UI and sends it to the ESP32 for display
- The transport between the ESP32 and the host is TCP/IP, so the ESP32 can be connected to the host through WiFi
- Upon TCP connect, the ESP32 display buffer is considered 'empty'
- We can have the ESP32 firmware already have 'screens' for certain states like: "WIFI disconnected", "Cannot connect to server". These screens can be available as files on a filesystem.
- ESP32 upon connecting to the server will send a packet that identifies: attached screen type, resolution, etc.. so that the HOST can know the frame buffer format of the ESP32.
- HOST rendering is fully deterministic
- ESP32 can notify host that display is OFF since the module has a mmWave sensor to detect human presence and distance, no presence and/or within a certain distance means that the display can be turned OFF.
- ESP32 can instruct the HOST for requesting a particular 'UI page' to display; for example, if the ESP32 request to display a local weather UI page, the HOST can collect weather information and render a local weather UI page. The main idea here is that it is potentially possible to display anything, even latest headline news.
- We do not mind UI partial updates, so we can always push the current state of the frame-buffer, no matter if there is a transition going on. This means that we can have the TCP connection push blocks at a certain rate instead of going full speed.

## Roles & authority

ESP32 (“Display Client”) owns:

- Hardware capabilities
- Power/display state
- UI intent (which page)
- Local fallback screens

Does NOT own:

- Rendering logic
- UI layout
- Sensor aggregation logic for visualizations

HOST (“UI Renderer”) owns:

- Rendering
- Framebuffer state
- Screen diffs
- External data (weather, news, etc.)

Assumes:

- ESP32 display buffer starts empty on connect
- Display remains correct as long as TCP stays connected

## Connection lifecycle (important change)

Because TCP gives you:

- Ordered bytes
- Hard disconnect detection
- Clean reconnect semantics

You should treat every TCP connection as a new session.

### On TCP connect

#### ESP32

Clears display or shows a local “CONNECTING” screen
Sends HELLO

#### HOST

Validates display characteristics
Resets its framebuffer state to “unknown”
Sends an initial full frame (block stream, paced)

#### During connection

No resync logic needed beyond reconnect.

## ESP32-S3 with Display

### Air Quality Monitor

These are ESP32-S3 air quality monitors (sensors; BME280, BH1750, SCD41, SC7A20H) with a small display can be used to display the air quality in the house. They can also be used to display the weather forecast, time, date, etc. 

### Bed Presence Sensor

These are ESP32-S3 bed presence sensors (sensors; RD03D, BH1750) with a small display can be used to detect if there is someone in the bed or not. They can also be used to display the time, date, etc.

## Mac Mini M4

The Mac Mini M4 can be used run the UI server, which renders the UI for the devices.

