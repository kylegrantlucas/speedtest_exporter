package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	prom_version "github.com/prometheus/common/version"

	"github.com/kylegrantlucas/speedtest_exporter/exporter"
	"github.com/kylegrantlucas/speedtest_exporter/version"
)

func init() {
	prometheus.MustRegister(
		prom_version.NewCollector("speedtest_exporter"),
	)
}

func main() {
	var (
		showVersion   = flag.Bool("version", false, "Print version information.")
		listenAddress = flag.String("web.listen-address", ":9112", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		serverID      = flag.Int("server-id", 0, "Speedtest server ID to use.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("Ookla Speedtest Prometheus exporter. v%s\n", version.Version)
		os.Exit(0)
	}

	log.Infof("starting speedtest exporter: %v", prom_version.Info())
	log.Infof("build context: %v", prom_version.BuildContext())

	exporter, err := exporter.NewExporter(*serverID)
	if err != nil {
		log.Errorf("creating exporter: %s", err)
		os.Exit(1)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Speedtest Exporter</title></head>
             <body>
             <h1>Speedtest Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.Infof("listening on: %v", *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Errorf("listen error: %s", err)
		os.Exit(1)
	}
}
