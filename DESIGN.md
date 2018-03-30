# Emitter.IO

Been looking into emitter.io, a high-performance pub/sub server that seems a lot more suitable
to our situation. It also uses MQTT as the message protocol which makes it interesting from
a IoT point of view.

An emitter client can subscribe to channels, so generally for all of our processes they should
subscribe to their config channel, state-request channel.

## Config Emitter Client
Listen to presence messages, when a subscriber registers we can send him the configuration.
Also when we detect that the configuration on disk has changed, we can hot-load it and send
it to the associated channel.

## Presence Emitter Client

``presenceEmitterClient.Subscribe(secret_key, "config/presence")`` \
``presenceEmitterClient.Subscribe(secret_key, "request/state/presence")``

## Automation Emitter Client

Subscribe to all state messages \
``automationEmitterClient.Subscribe(secret_key, "state/+")``


# Log
  Subscribe: 
  * ``Subscribe(secret_key, "log/+")``

  Action: Log messages to console

# Presence
  Subscribe: 
  * ``Subscribe(secret_key, "presence/+")``

  ```
  {
      "type": "config"
  }
  {
      "type": "pull"
  }
    
  ```
  Publish: 

  * ``Publish(secret_key, 'state/presence', json)``
  
  ```
  {
      "type": "state/presence"
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
  Subscribe: 
  * ``Subscribe(secret_key, "flux/+")``
  * ``Subscribe(secret_key, "state/sensor/calendar/season")``
  * ``Subscribe(secret_key, "state/sensor/calendar/tod")``
  * ``Subscribe(secret_key, "state/suncalc")``
  * ``Subscribe(secret_key, "state/weather")``


  Publish: 
  * ``Publish(secret_key, 'state/flux')``
  
  ```
  {
      "type": "state/flux"
      "ct": 100.0
      "bri": 100.0
  }
  ```

# AQI
  Subscribe: 
  * ``Subscribe(secret_key, "aqi/+")``

  Publish: 
  * ``Publish(secret_key, 'sensor/weather/aqi')``
  
  ```
  {
      "type": "state/aqi"
      "pm2.5": 100.0
  }
  ```

# Suncalc
  Subscribe: 
  * ``Subscribe(secret_key, "suncalc/+")``

  Publish: 
  * ``Publish(secret_key, 'state/suncalc')``
  
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
  Subscribe:
  * ``Subscribe(secret_key, "timeofday/+")``

  Publish: 
  * ``Publish(secret_key, 'state/timeofday')``
  
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
  Subscribe: 
  * ``Subscribe(secret_key, "lighting/+")``
  * ``Subscribe(secret_key, "state/light/+")``
  * ``Subscribe(secret_key, "state/flux")``

  Action: Turn On/Off, Set Ct/Bri


# Xiaomi Aqara
  Subscribe: 
  * ``Subscribe(secret_key, "xiaomi/+")``

  Publish:
  * ``Publish(secret_key, 'state/xiaomi/wireless-switch/id')``
  * ``Publish(secret_key, 'state/xiaomi/motion-sensor/id')``
  * ``Publish(secret_key, 'state/xiaomi/wireddualwallswitch/id')``

# Automation
  Subscribe: 
  * ``Subscribe(secret_key, "state/presence")``
  * ``Subscribe(secret_key, "state/sensor/calendar/jennifer")``
  * ``Subscribe(secret_key, "state/sensor/calendar/sophia")``
  * ``Subscribe(secret_key, "state/sensor/calendar/parents")``
  * ``Subscribe(secret_key, "state/sensor/calendar/alarm")``
  * ``Subscribe(secret_key, "state/sensor/calendar/tod")``
  * ``Subscribe(secret_key, "state/xiaomi/+/+")``

  Publish:
  * ``Publish(secret_key, 'state/light/hue')``
  * ``Publish(secret_key, 'state/light/yee')``
  * ``Publish(secret_key, 'state//hue')``

