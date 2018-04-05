package config

import (
	"encoding/json"
	"time"
)

// SensorState holds all information of a sensor
// e.g. sensor/weather/aqi
type SensorState struct {
	Name            string            `json:"name"`
	Time            time.Time         `json:"time"`
	BoolSensors     *[]BoolSensor     `json:"bool_sensors,omitempty"`
	IntSensors      *[]IntSensor      `json:"int_sensors,omitempty"`
	FloatSensors    *[]FloatSensor    `json:"float_sensors,omitempty"`
	ValueSensors    *[]ValueSensor    `json:"value_sensors,omitempty"`
	TimeSlotSensors *[]TimeSlotSensor `json:"timeslot_sensors,omitempty"`
}

// SensorStateFromJSON deserializes a JSON string and returns a SensorState object
func SensorStateFromJSON(jsonstr string) (SensorState, error) {
	var s SensorState
	err := json.Unmarshal([]byte(jsonstr), &s)
	return s, err
}

// ToJSON serializes a SensorState object into a JSON string
func (s SensorState) ToJSON() (string, error) {
	data, err := json.Marshal(s)
	if err == nil {
		return string(data), err
	}
	return "", err
}

// GetBoolAttr returns the value of a FloatSensor with name 'name'
func (s SensorState) GetBoolAttr(name string, defaultvalue bool) bool {
	if s.BoolSensors != nil {
		for _, fs := range *s.BoolSensors {
			if fs.Name == name {
				return fs.Value
			}
		}
	}
	return defaultvalue
}

// GetIntAttr returns the value of a FloatSensor with name 'name'
func (s SensorState) GetIntAttr(name string, defaultvalue int64) int64 {
	if s.IntSensors != nil {
		for _, fs := range *s.IntSensors {
			if fs.Name == name {
				return fs.Value
			}
		}
	}
	return defaultvalue
}

// GetFloatAttr returns the value of a FloatSensor with name 'name'
func (s SensorState) GetFloatAttr(name string, defaultvalue float64) float64 {
	if s.FloatSensors != nil {
		for _, fs := range *s.FloatSensors {
			if fs.Name == name {
				return fs.Value
			}
		}
	}
	return defaultvalue
}

// GetValueAttr returns the value of a FloatSensor with name 'name'
func (s SensorState) GetValueAttr(name string, defaultvalue string) string {
	if s.ValueSensors != nil {
		for _, fs := range *s.ValueSensors {
			if fs.Name == name {
				return fs.Value
			}
		}
	}
	return defaultvalue
}

// ValueSensor is a sensor holding a string as value
type ValueSensor struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// FloatSensor is a sensor holding a float as value
type FloatSensor struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// BoolSensor is a sensor holding a boolean as value
type BoolSensor struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

// IntSensor is a sensor holding an integer as value
type IntSensor struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// TimeSlotSensor is a sensor holding a TinmeSlot as value
// e.g. TimeSlotSensorAsJSON("sun.rise", sunrise_begin, sunrise_end)
// or adding multiple sensors:
//    sensor := NewSensorState("suncalc")
//    sensor.AddTimeSlotSensor("night.dawn", nightdawn, nightdawn_end)
//    sensor.AddTimeSlotSensor("astronomical.dawn", astronomicaldawn, astronomicaldawn_end)
//    ...
//    jsonstr, err := sensor.ToJSON()
//    ...
type TimeSlotSensor struct {
	Name  string    `json:"name"`
	Begin time.Time `json:"begin"`
	End   time.Time `json:"end"`
}

// FloatSensorAsJSON can be called as FloatSensorAsJSON("state.sensor.clouds", "clouds", 0.2) and
// you will receive the JSON string or an error.
func FloatSensorAsJSON(sensorname string, name string, value float64) (string, error) {
	sensorstate := SensorState{}
	sensorstate.Name = sensorname
	sensorstate.Time = time.Now()
	sensorstate.FloatSensors = &[]FloatSensor{FloatSensor{Name: name, Value: value}}
	jsonstr, err := sensorstate.ToJSON()
	return jsonstr, err
}

// ValueSensorAsJSON can be called as ValueSensorAsJSON("state.sensor.motion", "motion_ba9825ae", "on") and
// you will receive the JSON string or an error.
func ValueSensorAsJSON(sensorname string, name string, value string) (string, error) {
	sensorstate := SensorState{}
	sensorstate.Name = sensorname
	sensorstate.Time = time.Now()
	sensorstate.ValueSensors = &[]ValueSensor{ValueSensor{Name: name, Value: value}}
	jsonstr, err := sensorstate.ToJSON()
	return jsonstr, err
}

// IntSensorAsJSON can be called as IntSensorAsJSON("state.sensor.aqi", "aqi", 120) and
// you will receive the JSON string or an error.
func IntSensorAsJSON(sensorname string, name string, value int64) (string, error) {
	sensorstate := SensorState{}
	sensorstate.Name = sensorname
	sensorstate.Time = time.Now()
	sensorstate.AddIntSensor(name, value)
	jsonstr, err := sensorstate.ToJSON()
	return jsonstr, err
}

// BoolSensorAsJSON can be called as BoolSensorAsJSON("state.sensor.motion", "motion_98AE7", true) and
// you will receive the JSON string or an error.
func BoolSensorAsJSON(sensorname string, name string, value bool) (string, error) {
	sensorstate := SensorState{}
	sensorstate.Name = sensorname
	sensorstate.Time = time.Now()
	sensorstate.AddBoolSensor(name, value)
	jsonstr, err := sensorstate.ToJSON()
	return jsonstr, err
}

// TimeSlotSensorAsJSON can be called as TimeSlotSensorAsJSON("state.sensor.sun", "sun.rise", sunrise_begin, sunrise_end) and
// you will receive the JSON string or an error.
func TimeSlotSensorAsJSON(sensorname string, name string, begin time.Time, end time.Time) (string, error) {
	sensorstate := SensorState{}
	sensorstate.Name = sensorname
	sensorstate.Time = time.Now()
	sensorstate.AddTimeSlotSensor(name, begin, end)
	jsonstr, err := sensorstate.ToJSON()
	return jsonstr, err
}

// NewSensorState returns a SensorState object initialized with 'name' and time.Now()
func NewSensorState(name string) SensorState {
	sensorstate := SensorState{}
	sensorstate.Name = name
	return sensorstate
}

// AddBoolSensor adds an BoolSensor to SensorState
func (s SensorState) AddBoolSensor(name string, value bool) {
	if s.BoolSensors == nil {
		s.BoolSensors = &[]BoolSensor{BoolSensor{Name: name, Value: value}}
	} else {
		*s.BoolSensors = append(*s.BoolSensors, BoolSensor{Name: name, Value: value})
	}
}

// AddIntSensor adds an IntSensor to SensorState
func (s SensorState) AddIntSensor(name string, value int64) {
	if s.IntSensors == nil {
		s.IntSensors = &[]IntSensor{IntSensor{Name: name, Value: value}}
	} else {
		*s.IntSensors = append(*s.IntSensors, IntSensor{Name: name, Value: value})
	}
}

// AddFloatSensor adds a TimeSlotSensor to SensorState
func (s SensorState) AddFloatSensor(name string, value float64) {
	if s.FloatSensors == nil {
		s.FloatSensors = &[]FloatSensor{FloatSensor{Name: name, Value: value}}
	} else {
		*s.FloatSensors = append(*s.FloatSensors, FloatSensor{Name: name, Value: value})
	}
}

// AddValueSensor adds a TimeSlotSensor to SensorState
func (s *SensorState) AddValueSensor(name string, value string) {
	if s.ValueSensors == nil {
		s.ValueSensors = &[]ValueSensor{ValueSensor{Name: name, Value: value}}
	} else {
		*s.ValueSensors = append(*s.ValueSensors, ValueSensor{Name: name, Value: value})
	}
}

// AddTimeSlotSensor adds a TimeSlotSensor to SensorState
func (s SensorState) AddTimeSlotSensor(name string, begin time.Time, end time.Time) {
	if s.TimeSlotSensors == nil {
		s.TimeSlotSensors = &[]TimeSlotSensor{TimeSlotSensor{Name: name, Begin: begin, End: end}}
	} else {
		*s.TimeSlotSensors = append(*s.TimeSlotSensors, TimeSlotSensor{Name: name, Begin: begin, End: end})
	}
}
