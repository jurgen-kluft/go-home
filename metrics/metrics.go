package metrics

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/jurgen-kluft/go-home/config"
)

// Metrics library, mainly for:
// Raspberry Pi system stats
// Sensors
// Presence

// Mac Mini:
// * Telegraf
// * InfluxDB (localhost:8083)
// * Grafana (localhost:3000)

type Metrics struct {
	client  client.Client
	metrics map[string]*Metric
}

type Metric struct {
	name   string
	tags   map[string]string
	fields map[string]interface{}
	bp     client.BatchPoints
}

func New() (*Metrics, error) {
	var err error
	metrics := &Metrics{}
	metrics.metrics = map[string]*Metric{}

	// Create a new HTTPClient
	metrics.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: config.InfluxSecrets["host"],
		//Username: config.InfluxSecrets["username"],
		//Password: config.InfluxSecrets["password"],
	})

	return metrics, err
}

func (m *Metrics) Close() {
	m.client.Close()
}

func (m *Metrics) Register(name string, tags map[string]string, fields map[string]interface{}) {
	metric, exists := m.metrics[name]
	if !exists {
		metric = &Metric{name: name}
		m.metrics[name] = metric
	}
	metric.tags = tags
	metric.fields = fields
}

func (m *Metrics) Begin(name string) error {
	metric, exists := m.metrics[name]
	if !exists {
		return fmt.Errorf("No metric registered with name %s", name)
	}

	// Create a new point batch
	var err error
	metric.bp, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.InfluxSecrets["database"],
		Precision: "s",
	})

	if err != nil {
		return err
	}
	return nil
}

func (m *Metrics) Set(name string, field string, value float64) error {
	metric, exists := m.metrics[name]
	if !exists {
		return fmt.Errorf("No metric registered with name %s", name)
	}

	metric.fields[field] = value
	return nil
}

func (m *Metrics) Send(name string) error {
	metric, exists := m.metrics[name]
	if !exists {
		return fmt.Errorf("No metric registered with name %s", name)
	}

	pt, err := client.NewPoint(name, metric.tags, metric.fields, time.Now())
	if err != nil {
		return err
	}

	metric.bp.AddPoint(pt)

	// Write the batch
	err = m.client.Write(metric.bp)
	return err
}
