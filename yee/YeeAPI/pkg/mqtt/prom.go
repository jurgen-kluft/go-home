package mqtt

import "github.com/prometheus/client_golang/prometheus"

var (
	connectedBulbs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "yeelight_connected_bulbs",
			Help: "asdasdasdasdd",
		})
	devicesDiscovered = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "yeelight_discover_processed",
			Help: " asasd",
		})
)

func init() {
	prometheus.MustRegister(connectedBulbs, devicesDiscovered)

}
