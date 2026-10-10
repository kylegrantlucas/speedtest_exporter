# speedtest_exporter

Prometheus exporter that runs the official [Ookla Speedtest CLI](https://www.speedtest.net/apps/cli)
on every scrape and exposes the results.

## Metrics

| Metric | Unit |
| --- | --- |
| `speedtest_download_bandwidth`, `speedtest_upload_bandwidth` | Mbps |
| `speedtest_download_bytes`, `speedtest_upload_bytes` | bytes |
| `speedtest_download_elapsed`, `speedtest_upload_elapsed` | ms |
| `speedtest_ping_latency`, `speedtest_ping_jitter` | ms |

A speedtest takes ~30s and uses real bandwidth: set a long `scrape_interval` (e.g. `1h`)
and a `scrape_timeout` of at least `2m`.

## Flags

| Flag | Default | |
| --- | --- | --- |
| `-web.listen-address` | `:9112` | listen address |
| `-web.telemetry-path` | `/metrics` | metrics path |
| `-server-id` | `0` (auto) | pin an Ookla server ID |

Running the exporter accepts the Ookla license and GDPR terms on your behalf.

## Run

```sh
docker build -t speedtest_exporter .
docker run --rm -p 9112:9112 speedtest_exporter

# or with Nix (wraps the Ookla CLI; the flake allows that unfree package)
nix run github:kylegrantlucas/speedtest_exporter
```

Multi-arch image: `docker buildx build --platform linux/amd64,linux/arm64 .`

## Develop

```sh
nix develop   # or direnv allow
go test ./... && go vet ./...
```
