package exporter

import (
	"fmt"

	"github.com/kylegrantlucas/speedtest_exporter/speedtest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// Exporter collects speedtest stats and exports
// them using the prometheus metrics package.
type Exporter struct {
	Client *speedtest.Client
}

func NewExporter(serverID int) (*Exporter, error) {
	client, err := speedtest.NewClient(serverID)
	if err != nil {
		return nil, fmt.Errorf("creating client: %s", err)
	}

	return &Exporter{
		Client: client,
	}, nil
}

// Describe iterates over all the metrics offered by speedtest to
// describe themselves to satisfy prometheus.Collector
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range speedtest.PromMetrics() {
		ch <- metric
	}
}

// Collect fetches the stats using the speedtest package
// to satisfy prometheus.Collector
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if e.Client == nil {
		log.Errorf("speedtest client not configured")
		return
	}

	metrics, err := e.Client.NetworkMetrics()
	if err != nil {
		log.Errorf("can't run export: %s", err)
		return
	}

	for _, metric := range metrics.ToPrometheusMetrics() {
		ch <- metric
	}
}
