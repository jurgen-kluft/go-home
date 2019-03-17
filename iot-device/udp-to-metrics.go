package main

import (
	"encoding/binary"
	"time"

	"github.com/jurgen-kluft/go-home/metrics"
)

type movingAverage struct {
	history []float64
	index   int
}

func newFilter(sizeOfHistory int) *movingAverage {
	filter := &movingAverage{history: make([]float64, sizeOfHistory), index: -1}
	return filter
}

func (m *movingAverage) sample(sample float64) float64 {
	if m.index == -1 {
		for i := range m.history {
			m.history[i] = sample
		}
		m.index = 0
	}

	m.history[m.index] = sample
	m.index = (m.index + 1) % len(m.history)

	sum := 0.0
	for _, s := range m.history {
		sum += s
	}
	return sum / float64(len(m.history))
}

type context struct {
	name           string
	metrics        *metrics.Metrics
	mavTemperature *movingAverage
	mavHumidity    *movingAverage
	mavPressure    *movingAverage
	index          int64
	timestamp      time.Time
}

type sensordata struct {
	stride int

	temperature int
	humidity    int
	pressure    int

	magnetic     int
	acceleration int
	gyroscope    int

	data []byte
}

func (s *sensordata) init() {
	s.stride = (1 + 1 + 3 + 3*3) * 4

	s.temperature = 8
	s.humidity = 12
	s.pressure = 16

	s.magnetic = 20
	s.acceleration = 32
	s.gyroscope = 44
}

func (s *sensordata) readTemperature() float64 {
	t := int64(0)
	o := s.temperature
	for i := 0; i <= 16; i++ {
		t += int64(binary.LittleEndian.Uint32(s.data[o : o+4]))
		o += s.stride
	}
	return float64(t/16) / 100.0
}
func (s *sensordata) readHumidity() float64 {
	t := int64(0)
	o := s.humidity
	for i := 0; i <= 16; i++ {
		t += int64(binary.LittleEndian.Uint32(s.data[o : o+4]))
		o += s.stride
	}
	return float64(t/16) / 100.0
}
func (s *sensordata) readPressure() float64 {
	t := int64(0)
	o := s.pressure
	for i := 0; i <= 16; i++ {
		t += int64(binary.LittleEndian.Uint32(s.data[o : o+4]))
		o += s.stride
	}
	return float64(t/16) / 100.0
}

func read_int32(data []byte, o int) int32 {
	return int32(uint32(data[o+0]) + uint32(data[o+1])<<8 + uint32(data[o+2])<<16 + uint32(data[o+3])<<24)
}
func (s *sensordata) readMagnetic(i int) (X float64, Y float64, Z float64) {
	o := s.magnetic + i*s.stride
	x32 := read_int32(s.data, o)
	X = float64(x32)
	y32 := read_int32(s.data, o+4)
	Y = float64(y32)
	z32 := read_int32(s.data, o+8)
	Z = float64(z32)
	return X, Y, Z
}

func (s *sensordata) readAcceleration(i int) (X float64, Y float64, Z float64) {
	o := s.acceleration + i*s.stride
	x32 := read_int32(s.data, o)
	X = float64(x32)
	y32 := read_int32(s.data, o+4)
	Y = float64(y32)
	z32 := read_int32(s.data, o+8)
	Z = float64(z32)
	return X, Y, Z
}

func (s *sensordata) readGyroscope(i int) (X float64, Y float64, Z float64) {
	o := s.gyroscope + i*s.stride
	x32 := read_int32(s.data, o)
	X = float64(x32)
	y32 := read_int32(s.data, o+4)
	Y = float64(y32)
	z32 := read_int32(s.data, o+8)
	Z = float64(z32)
	return X, Y, Z
}

func new() *context {
	c := &context{}
	c.name = "az3166"
	c.metrics, _ = metrics.New()
	c.metrics.Register("Environment Sensors", map[string]string{"T": "Temperature", "H": "Humidity", "P": "Pressure"}, map[string]interface{}{"T": 20, "H": 20, "P": 1000})
	c.metrics.Register("Movement Sensors", map[string]string{"MX": "Magnetic-X", "MY": "Magnetic-Y", "MZ": "Magnetic-Z", "AX": "Acceleration-X", "AY": "Acceleration-Y", "AZ": "Acceleration-Z", "GX": "Gyroscope-X", "GY": "Gyroscope-Y", "GZ": "Gyroscope-Z"}, map[string]interface{}{"MX": 0, "MY": 0, "MZ": 0, "AX": 0, "AY": 0, "AZ": 0, "GX": 0, "GY": 0, "GZ": 0})
	c.mavTemperature = newFilter(2)
	c.mavHumidity = newFilter(2)
	c.mavPressure = newFilter(2)
	c.index = 0
	return c
}

func main() {

	inbound, inpktpool, err := Listen(":7331", nil)
	if err != nil {
		// handle err
	}

	c := new()
	p := &sensordata{}
	p.init()

	connected := true
	for connected {

		select {
		case inpkt := <-inbound:
			// Do something with UDP packet
			p.data = inpkt
			t := time.Now()

			em, _ := c.metrics.Begin("Environment Sensors")
			temp := p.readTemperature()
			temp = c.mavTemperature.sample(temp)
			hum := p.readHumidity()
			hum = c.mavHumidity.sample(hum)
			press := p.readPressure()
			press = c.mavPressure.sample(press)

			if c.index > 0 {
				em.Set("T", temp)
				em.Set("H", hum)
				em.Set("P", press)
				c.metrics.SendMetric(em, c.timestamp)

				//fmt.Printf("Temp: %v\n", temp)
				//fmt.Printf("Humidity: %v\n", hum)
				//fmt.Printf("Pressure: %v\n", press)

				deltatime := t.Sub(c.timestamp) / 16
				di := deltatime
				for i := 0; i < 16; i++ {
					m, e := c.metrics.Begin("Movement Sensors")
					if e == nil {
						x, y, z := p.readMagnetic(i)
						m.Set("MX", x)
						m.Set("MY", y)
						m.Set("MZ", z)

						x, y, z = p.readAcceleration(i)
						m.Set("AX", x)
						m.Set("AY", y)
						m.Set("AZ", z)

						if i == 0 {
							//fmt.Printf("Movement Acc: %v / %v / %v\n", x, y, z)
						}

						x, y, z = p.readGyroscope(i)
						m.Set("GX", x)
						m.Set("GY", y)
						m.Set("GZ", z)

						if i == 0 {
							//fmt.Printf("Movement Gyro: %v / %v / %v\n", x, y, z)
						}

						ti := c.timestamp.Add(di)
						c.metrics.SendMetric(m, ti)
					}
					di += deltatime
				}
			}

			c.index++
			c.timestamp = t

			// Release UDP packet back to pool
			inpktpool.Release(inpkt)

		case <-time.After(time.Second * 1):

		}
	}

	return
}
