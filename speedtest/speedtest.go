package speedtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
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
		"Upload Bytes",
		nil, nil,
	),
	"upload_elapsed": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "upload_elapsed"),
		"Upload Elapsed (ms)",
		nil, nil,
	),
}

// Client runs the Ookla speedtest CLI.
type Client struct {
	// Binary is the speedtest CLI to execute (default "speedtest").
	Binary   string
	ServerID int
	Timeout  time.Duration
}

// NewClient defines a new client for Speedtest.
func NewClient(serverID int) (*Client, error) {
	return &Client{Binary: "speedtest", ServerID: serverID, Timeout: 2 * time.Minute}, nil
}

// Args returns the CLI arguments for a run.
func (client *Client) Args() []string {
	args := []string{"--accept-license", "--accept-gdpr", "--format=json", "--progress=no"}
	if client.ServerID > 0 {
		args = append(args, "--server-id="+strconv.Itoa(client.ServerID))
	}
	return args
}

// NetworkMetrics runs a speedtest and parses its JSON result.
func (client *Client) NetworkMetrics(ctx context.Context) (Result, error) {
	if client.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, client.Timeout)
		defer cancel()
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, client.Binary, client.Args()...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return Result{}, fmt.Errorf("running %s: %w: %s", client.Binary, err, stderr.String())
	}

	return ParseResult(stdout.Bytes())
}

// ParseResult decodes the JSON emitted by `speedtest --format=json`.
func ParseResult(data []byte) (Result, error) {
	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("unmarshaling speedtest output: %w", err)
	}
	return result, nil
}

// PromMetrics returns every metric description.
func PromMetrics() []*prometheus.Desc {
	descs := make([]*prometheus.Desc, 0, len(metricDescriptions))
	for _, v := range metricDescriptions {
		descs = append(descs, v)
	}
	return descs
}
