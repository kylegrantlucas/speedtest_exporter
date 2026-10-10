package speedtest

import (
	"math"
	"slices"
	"testing"
)

const sample = `{"type":"result","ping":{"jitter":1.5,"latency":10.25},
"download":{"bandwidth":12500000,"bytes":150000000,"elapsed":12000},
"upload":{"bandwidth":2500000,"bytes":30000000,"elapsed":9000}}`

func TestParseResultMetrics(t *testing.T) {
	r, err := ParseResult([]byte(sample))
	if err != nil {
		t.Fatal(err)
	}
	m := r.ToMetricsMap()
	want := map[string]float64{
		"download_bandwidth": 100, "upload_bandwidth": 20,
		"ping_latency": 10.25, "ping_jitter": 1.5,
		"download_bytes": 150000000, "upload_elapsed": 9000,
	}
	for k, v := range want {
		if got := m[k].(float64); math.Abs(got-v) > 1e-9 {
			t.Errorf("%s = %v, want %v", k, got, v)
		}
	}
	if n := len(r.ToPrometheusMetrics()); n != len(metricDescriptions) {
		t.Errorf("got %d metrics, want %d", n, len(metricDescriptions))
	}
}

func TestParseResultInvalid(t *testing.T) {
	if _, err := ParseResult([]byte("nope")); err == nil {
		t.Fatal("expected error")
	}
}

func TestArgs(t *testing.T) {
	c, _ := NewClient(1234)
	if !slices.Contains(c.Args(), "--server-id=1234") {
		t.Errorf("args missing server id: %v", c.Args())
	}
	c.ServerID = 0
	if slices.ContainsFunc(c.Args(), func(a string) bool { return a == "--server-id=0" }) {
		t.Error("unexpected server id")
	}
}
