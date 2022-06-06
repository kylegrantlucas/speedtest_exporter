package speedtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "speedtest"
)

var metricDescriptions = map[string]*prometheus.Desc{
	"ping_latency": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ping_latency"),
		"Ping Latency (ms)",
		nil, nil,
	),
	"ping_jitter": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ping_jitter"),
		"Ping Jitter (ms)",
		nil, nil,
	),
	"download_bandwidth": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "download_bandwidth"),
		"Download Bandwidth (Mbps)",
		nil, nil,
	),
	"download_bytes": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "download_bytes"),
		"Download Bytes",
		nil, nil,
	),
	"download_elapsed": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "download_elapsed"),
		"Download Elapsed (ms)",
		nil, nil,
	),
	"upload_bandwidth": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "upload_bandwidth"),
		"Upload Bandwidth (Mbps)",
		nil, nil,
	),
	"upload_bytes": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "upload_bytes"),
		"Download Bytes",
		nil, nil,
	),
	"upload_elapsed": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "upload_elapsed"),
		"Download Elapsed (ms)",
		nil, nil,
	),
}

type Client struct {
	ServerID int
}

type Results struct {
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

// NewClient defines a new client for Speedtest
func NewClient(serverID int) (*Client, error) {
	return &Client{ServerID: serverID}, nil
}

func (client *Client) NetworkMetrics() (Results, error) {
	result := Results{}

	cmdArr := []string{"--accept-license", "-f", "json"}

	if client.ServerID > 0 {
		cmdArr = append(cmdArr, "-s", fmt.Sprintf("%d", client.ServerID))
	}

	cmd := exec.Command("speedtest", cmdArr...)
	var outBuf bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return result, fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}

	out := outBuf.String()
	if err != nil {
		return result, fmt.Errorf("reading speedtest output: %s", err)
	}

	err = json.Unmarshal([]byte(out), &result)
	if err != nil {
		return result, fmt.Errorf("unmarshaling speedtest output: %s", err)
	}

	return result, nil
}

func (result Results) ToMetricsMap() map[string]interface{} {
	return map[string]interface{}{
		"download_bandwidth": result.Download.Bandwidth,
		"download_bytes":     result.Download.Bytes,
		"download_elapsed":   result.Download.Elapsed,
		"upload_bandwidth":   result.Upload.Bandwidth,
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

func (result Results) ToPrometheusMetrics() []prometheus.Metric {
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

func PromMetrics() []*prometheus.Desc {
	descs := []*prometheus.Desc{}
	for _, v := range metricDescriptions {
		descs = append(descs, v)
	}
	return descs
}
