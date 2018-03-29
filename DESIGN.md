# Mist subscriptions

Why not replace the current tag matching with something more flexible, like:

Subscribe: messages tagged with: ``(('hue' OR 'flux') AND state') OR 'lighting'``  

``flux hue | state & lighting |``

Subscribe: messages tagged with: ``('xiaomi' AND state') OR 'xiaomi'``

``xiaomi state & xiaomi |``


# Mist client connections

  Clients can connect and disconnect at any time, but a client might need state from others that might not have connected yet.
  We do have the message to request a list of clients from the server.


# Log
  Subscribe: messages tagged with: ``'log'``
  Action: Log messages to console

# Presence
  Subscribe: messages tagged with ``'presence'``

  ```
  {
      "type": "config"
  }
    
  ```
  Publish: messages tagged with: ``'presence'``
  
  ```
  {
      "type": "presence"
      "devices": [
          {
              "name": "A mobile phone",
              "state": "home",
          },
          {
              "name": "A kindle",
              "state": "away",
          }

      ]
  }
  ```

# Flux
  Subscribe: messages tagged with: 'flux' OR (('suncalc' OR 'weather') AND 'state')

  Publish: messages tagged with: ``'flux', 'state'``
  
  ```
  {
      "type": "flux"
      "ct": 100.0
      "bri": 100.0
  }
  ```

# AQI
  Subscribe: messages tagged with: ``'aqi'``

  Publish: messages tagged with: 'aqi-state'
  
  ```
  {
      "type": "aqi"
      "pm2.5": 100.0
  }
  ```

# Suncalc
  Subscribe: messages tagged with: ``'suncalc'``

  Publish: messages tagged with: ``'suncalc', 'state'``
  
  ```
  {
      "type": "suncalc",
      "name": "night.dawn",
      "descr": "midnight to twilight, 2nd part of the night",
      "begin": "12:00AM",
      "end": "05:00AM"
  }
  ```

# TimeOfDay
  Subscribe: messages tagged with: ``'timeofday'``

  Publish: messages tagged with: ``'timeofday', 'state'``
  
  ```
  {
      "type": "timeofday",
      "name": "breakfast",
      "descr": "early morning from 6:00AM to 9:00AM",
      "begin": "6:00AM",
      "end": "9:00AM"
  }
  ```


# Lighting
  Subscribe: messages tagged with: ``(('hue' OR 'yee' OR 'flux' OR 'weather') AND state') OR 'lighting'``

  Action: Turn On/Off, Set Ct/Bri


# Xiaomi Aqara
  Subscribe: messages tagged with: ``'xiaomi'``


# Automation
  Subscribe: messages tagged with: ``'state'``

  Receive events from:
  - Xiaomi Aqara WirelessSwitch
  - Xiaomi Aqara Motion Sensor
  - Calendar
  - 

