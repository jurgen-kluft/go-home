package metrics

import (
	"github.com/alexcesaro/statsd"
)

// Metrics library, mainly for:
// Raspberry Pi system stats
// Motion sensors
// Presence

// Mac Mini:
// * InfluxDB (localhost:8083)
// * Telegraf (configured as StatsD server)
// * Grafana (localhost:3000)

// StatsD client library:
// * https://github.com/peterbourgon/g2s
// * https://github.com/alexcesaro/statsd/tree/v2.0.0

type Metrics struct {
	client *statsd.Client
}

type Timing struct {
	name   string
	timing statsd.Timing
}

func New() (*Metrics, error) {
	var err error
	metrics := &Metrics{}
	metrics.client, err = statsd.New()
	metrics.client.NewTiming()
	return metrics, err
}

func (m *Metrics) Close() {
	m.client.Close()
}

func (m *Metrics) Measure(name string, value int64) {
	m.client.Gauge(name, value)
}

func (m *Metrics) Increment(name string) {
	m.client.Increment(name)
}

func (m *Metrics) BeginTiming(name string) Timing {
	t := Timing{name: name, timing: m.client.NewTiming()}
	return t
}

func (t Timing) End() {
	t.timing.Send(t.name)
}
