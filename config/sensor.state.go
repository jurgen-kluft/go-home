package config

import (
	"encoding/json"
	"time"
)

// SensorState holds all information of a sensor
// e.g. sensor/weather/aqi
type SensorState struct {
	Name         string         `json:"name"`
	Time         time.Time      `json:"time"`
	BoolAttrs    *[]BoolAttr    `json:"boolattrs,omitempty"`
	IntAttrs     *[]IntAttr     `json:"intattrs,omitempty"`
	FloatAttrs   *[]FloatAttr   `json:"floatattrs,omitempty"`
	StringAttrs  *[]StringAttr  `json:"valueattrs,omitempty"`
	TimeWndAttrs *[]TimeWndAttr `json:"timeslotattrs,omitempty"`
}

// SensorStateFromJSON deserializes a JSON string and returns a SensorState object
func SensorStateFromJSON(jsonstr string) (*SensorState, error) {
	s := &SensorState{}
	err := json.Unmarshal([]byte(jsonstr), s)
	return s, err
}

// ToJSON serializes a SensorState object into a JSON string
func (s *SensorState) ToJSON() (string, error) {
	data, err := json.Marshal(s)
	if err == nil {
		return string(data), err
	}
	return "", err
}

// GetBoolAttr returns the value of a FloatAttr with name 'name'
func (s *SensorState) GetBoolAttr(name string, defaultvalue bool) bool {
	if s.BoolAttrs != nil {
		for _, fs := range *s.BoolAttrs {
			if fs.Name == name {
				return fs.Value
			}
		}
	}
	return defaultvalue
}

// GetIntAttr returns the value of a FloatAttr with name 'name'
func (s *SensorState) GetIntAttr(name string, defaultvalue int64) int64 {
	if s.IntAttrs != nil {
		for _, fs := range *s.IntAttrs {
			if fs.Name == name {
				return fs.Value
			}
		}
	}
	return defaultvalue
}

// GetFloatAttr returns the value of a FloatAttr with name 'name'
func (s *SensorState) GetFloatAttr(name string, defaultvalue float64) float64 {
	if s.FloatAttrs != nil {
		for _, fs := range *s.FloatAttrs {
			if fs.Name == name {
				return fs.Value
			}
		}
	}
	return defaultvalue
}

func (s *SensorState) ExecFloatAttr(name string, action func(float64)) bool {
	if s.FloatAttrs != nil {
		for _, fs := range *s.FloatAttrs {
			if fs.Name == name {
				action(fs.Value)
				return true
			}
		}
	}
	return false
}

// GetValueAttr returns the value of a FloatAttr with name 'name'
func (s *SensorState) GetValueAttr(name string, defaultvalue string) string {
	if s.StringAttrs != nil {
		for _, fs := range *s.StringAttrs {
			if fs.Name == name {
				return fs.Value
			}
		}
	}
	return defaultvalue
}

func (s *SensorState) ExecValueAttr(name string, action func(string)) bool {
	if s.StringAttrs != nil {
		for _, fs := range *s.StringAttrs {
			if fs.Name == name {
				action(fs.Value)
				return true
			}
		}
	}
	return false
}

// StringAttr is a sensor holding a string as value
type StringAttr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// FloatAttr is a sensor holding a float as value
type FloatAttr struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// BoolAttr is a sensor holding a boolean as value
type BoolAttr struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

// IntAttr is a sensor holding an integer as value
type IntAttr struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// TimeWndAttr is a sensor holding a TimeWnd as value
// e.g. TimeWndAttrAsJSON("sun.rise", sunrise_begin, sunrise_end)
// or adding multiple sensors:
//    sensor := NewSensorState("suncalc")
//    sensor.AddTimeWndAttr("night.dawn", nightdawn, nightdawn_end)
//    sensor.AddTimeWndAttr("astronomical.dawn", astronomicaldawn, astronomicaldawn_end)
//    ...
//    jsonstr, err := sensor.ToJSON()
//    ...
type TimeWndAttr struct {
	Name  string    `json:"name"`
	Begin time.Time `json:"begin"`
	End   time.Time `json:"end"`
}

// FloatAttrAsJSON can be called as FloatAttrAsJSON("state.sensor.clouds", "clouds", 0.2) and
// you will receive the JSON string or an error.
func FloatAttrAsJSON(sensorname string, name string, value float64) (string, error) {
	sensorstate := NewSensorState(sensorname)
	sensorstate.FloatAttrs = &[]FloatAttr{FloatAttr{Name: name, Value: value}}
	jsonstr, err := sensorstate.ToJSON()
	return jsonstr, err
}

// StringAttrAsJSON can be called as StringAttrAsJSON("state.sensor.motion", "motion_ba9825ae", "on") and
// you will receive the JSON string or an error.
func StringAttrAsJSON(sensorname string, name string, value string) (string, error) {
	sensorstate := NewSensorState(sensorname)
	sensorstate.StringAttrs = &[]StringAttr{StringAttr{Name: name, Value: value}}
	jsonstr, err := sensorstate.ToJSON()
	return jsonstr, err
}

// IntAttrAsJSON can be called as IntAttrAsJSON("state.sensor.aqi", "aqi", 120) and
// you will receive the JSON string or an error.
func IntAttrAsJSON(sensorname string, name string, value int64) (string, error) {
	sensorstate := NewSensorState(sensorname)
	sensorstate.AddIntAttr(name, value)
	jsonstr, err := sensorstate.ToJSON()
	return jsonstr, err
}

// BoolAttrAsJSON can be called as BoolAttrAsJSON("state.sensor.motion", "motion_98AE7", true) and
// you will receive the JSON string or an error.
func BoolAttrAsJSON(sensorname string, name string, value bool) (string, error) {
	sensorstate := NewSensorState(sensorname)
	sensorstate.AddBoolAttr(name, value)
	jsonstr, err := sensorstate.ToJSON()
	return jsonstr, err
}

// TimeWndAttrAsJSON can be called as TimeWndAttrAsJSON("state.sensor.sun", "sun.rise", sunrise_begin, sunrise_end) and
// you will receive the JSON string or an error.
func TimeWndAttrAsJSON(sensorname string, name string, begin time.Time, end time.Time) (string, error) {
	sensorstate := NewSensorState(sensorname)
	sensorstate.AddTimeWndAttr(name, begin, end)
	jsonstr, err := sensorstate.ToJSON()
	return jsonstr, err
}

// NewSensorState returns a SensorState object initialized with 'name' and time.Now()
func NewSensorState(name string) *SensorState {
	sensorstate := &SensorState{}
	sensorstate.Name = name
	sensorstate.Time = time.Now()
	return sensorstate
}

// AddBoolAttr adds an BoolAttr to SensorState
func (s *SensorState) AddBoolAttr(name string, value bool) {
	if s.BoolAttrs == nil {
		s.BoolAttrs = &[]BoolAttr{BoolAttr{Name: name, Value: value}}
	} else {
		*s.BoolAttrs = append(*s.BoolAttrs, BoolAttr{Name: name, Value: value})
	}
}

// AddIntAttr adds an IntAttr to SensorState
func (s SensorState) AddIntAttr(name string, value int64) {
	if s.IntAttrs == nil {
		s.IntAttrs = &[]IntAttr{IntAttr{Name: name, Value: value}}
	} else {
		*s.IntAttrs = append(*s.IntAttrs, IntAttr{Name: name, Value: value})
	}
}

// AddFloatAttr adds a TimeWndAttr to SensorState
func (s *SensorState) AddFloatAttr(name string, value float64) {
	if s.FloatAttrs == nil {
		s.FloatAttrs = &[]FloatAttr{FloatAttr{Name: name, Value: value}}
	} else {
		*s.FloatAttrs = append(*s.FloatAttrs, FloatAttr{Name: name, Value: value})
	}
}

// AddStringAttr adds a TimeWndAttr to SensorState
func (s *SensorState) AddStringAttr(name string, value string) {
	if s.StringAttrs == nil {
		s.StringAttrs = &[]StringAttr{StringAttr{Name: name, Value: value}}
	} else {
		*s.StringAttrs = append(*s.StringAttrs, StringAttr{Name: name, Value: value})
	}
}

// AddTimeWndAttr adds a TimeWndAttr to SensorState
func (s *SensorState) AddTimeWndAttr(name string, begin time.Time, end time.Time) {
	if s.TimeWndAttrs == nil {
		s.TimeWndAttrs = &[]TimeWndAttr{TimeWndAttr{Name: name, Begin: begin, End: end}}
	} else {
		*s.TimeWndAttrs = append(*s.TimeWndAttrs, TimeWndAttr{Name: name, Begin: begin, End: end})
	}
}
