package speedtest

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Result struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Ping      struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
	} `json:"ping"`
	Download struct {
		Bandwidth float64 `json:"bandwidth"`
		Bytes     float64 `json:"bytes"`
		Elapsed   float64 `json:"elapsed"`
	} `json:"download"`
	Upload struct {
		Bandwidth float64 `json:"bandwidth"`
		Bytes     float64 `json:"bytes"`
		Elapsed   float64 `json:"elapsed"`
	} `json:"upload"`
	Isp       string `json:"isp"`
	Interface struct {
		InternalIP string `json:"internalIp"`
		Name       string `json:"name"`
		MacAddr    string `json:"macAddr"`
		IsVpn      bool   `json:"isVpn"`
		ExternalIP string `json:"externalIp"`
	} `json:"interface"`
	Server struct {
		ID       int    `json:"id"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
		IP       string `json:"ip"`
	} `json:"server"`
	Result struct {
		ID        string `json:"id"`
		URL       string `json:"url"`
		Persisted bool   `json:"persisted"`
	} `json:"result"`
}

func (result Result) ToMetricsMap() map[string]interface{} {
	return map[string]interface{}{
		"download_bandwidth": result.Download.Bandwidth / 100000,
		"download_bytes":     result.Download.Bytes,
		"download_elapsed":   result.Download.Elapsed,
		"upload_bandwidth":   result.Upload.Bandwidth / 100000,
		"upload_bytes":       result.Upload.Bytes,
		"upload_elapsed":     result.Upload.Elapsed,
		"ping_latency":       result.Ping.Latency,
		"ping_jitter":        result.Ping.Jitter,
		// "type":                  result.Type,
		// "isp":                   result.Isp,
		// "interface_name":        result.Interface.Name,
		// "interface_internal_ip": result.Interface.InternalIP,
		// "interface_mac_addr":    result.Interface.MacAddr,
		// "interface_is_vpn":      result.Interface.IsVpn,
		// "interface_external_ip": result.Interface.ExternalIP,
		// "server_id":             result.Server.ID,
		// "server_host":           result.Server.Host,
		// "server_port":           result.Server.Port,
		// "server_name":           result.Server.Name,
		// "server_location":       result.Server.Location,
		// "server_country":        result.Server.Country,
		// "server_ip":             result.Server.IP,
		// "result_id":             result.Result.ID,
		// "result_url":            result.Result.URL,
		// "result_persisted":      result.Result.Persisted,
	}
}

func (result Result) ToPrometheusMetrics() []prometheus.Metric {
	metrics := []prometheus.Metric{}
	for k, v := range result.ToMetricsMap() {
		metrics = append(metrics, prometheus.MustNewConstMetric(
			metricDescriptions[k],
			prometheus.GaugeValue,
			v.(float64),
		))
	}

	return metrics
}
