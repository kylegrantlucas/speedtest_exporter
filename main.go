package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/kylegrantlucas/speedtest_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	prom_version "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	prometheus.MustRegister(
		prom_version.NewCollector("speedtest_exporter"),
	)
}

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9112", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		serverID      = flag.Int("server-id", 0, "Speedtest server ID to use.")
	)
	flag.Parse()

	exporter, err := exporter.NewExporter(*serverID)
	if err != nil {
		log.Printf("creating exporter: %s", err)
		os.Exit(1)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Speedtest Exporter</title></head>
             <body>
             <h1>Speedtest Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.Printf("listening on: %v", *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Printf("listen error: %s", err)
		os.Exit(1)
	}
}
