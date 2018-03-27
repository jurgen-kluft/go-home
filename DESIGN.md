# Mist subscriptions

Why not replace the current tag matching with something more flexible, like:

Subscribe: messages tagged with: ``(('hue' OR 'flux') AND state') OR 'lighting'``  

``$flux $hue | $state & $lighting |``

Subscribe: messages tagged with: ``('xiaomi' AND state') OR 'xiaomi'``

``$xiaomi $state & $xiaomi |``

Sub

Optimization, we can replace the string with something more optimum, like:
```
package main

import "fmt"

const (
	var_false = iota
	var_true
	var_xiaomi
	var_state
	var_hue
	var_yee
	var_flux
	var_lighting
	var_wemo
	var_mask = 0x0FFF
)
const (
	op_mask = 0xF000
	op_and  = 0x1000
	op_or   = 0x2000
	op_xor  = 0x3000
	op_not  = 0x4000
)

// ('xiaomi' AND state') OR 'xiaomi' OR 'hue'
func automation_message_filter() []int16 {
	expr := make([]int16, 0, 0)
	expr = append(expr, var_xiaomi)
	expr = append(expr, var_state)
	expr = append(expr, op_and)
	expr = append(expr, var_xiaomi)
	expr = append(expr, op_or)
	expr = append(expr, var_hue)
	expr = append(expr, op_or)
	return expr
}

func evaluate_filter(vars []bool, expr []int16) bool {
	r := []bool{}
	for _, s := range expr {
		i := len(r)
		switch s {
		case op_and:
			r[i-2] = r[i-2] && r[i-1]
			r = r[:len(r)-1]
			break
		case op_or:
			r[i-2] = r[i-2] || r[i-1]
			r = r[:len(r)-1]
			break
		case op_xor:
			r[i-2] = (r[i-2] && !r[i-1]) || (!r[i-2] && r[i-1])
			r = r[:len(r)-1]
			break
		case op_not:
			r[i-1] = !r[i-1]
			break
		default:
			r = append(r, vars[s])
			break
		}
	}
	return r[0]
}

func main() {
	vars := make([]bool, 16, 16)
	vars[var_yee] = true
	expr := automation_message_filter()
	result := evaluate_filter(vars, expr)
	fmt.Println("Result", result, " variable")
}

```

So in mist if we change the message and subscription into:
```
	Message struct {
		Command string   `json:"command"`
		Tags    []int16  `json:"tags,omitempty"`
		Data    string   `json:"data,omitempty"`
		Error   string   `json:"error,omitempty"`
	}

```

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

  Publish: message tagged as: ``'config'``

  ```
  {
      "name": "presence"
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

  Publish: message tagged as: ``'config'``

  ```
  {
      "name": "flux"
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

  Publish: message tagged as: ``'config'``

  ```
  {
      "name": "aqi"
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

  Publish: message tagged as: ``'config'``

  ```
  {
      "name": "suncalc"
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

  Publish: message tagged as: ``'config'``

  ```
  {
      "name": "timeofday"
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

