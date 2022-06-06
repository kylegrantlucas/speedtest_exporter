package speedtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

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

// NewClient defines a new client for Speedtest
func NewClient(serverID int) (*Client, error) {
	return &Client{ServerID: serverID}, nil
}

func (client *Client) NetworkMetrics() (Result, error) {
	result := Result{}

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

func PromMetrics() []*prometheus.Desc {
	descs := []*prometheus.Desc{}
	for _, v := range metricDescriptions {
		descs = append(descs, v)
	}
	return descs
}
