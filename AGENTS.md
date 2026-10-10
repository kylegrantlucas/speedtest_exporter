# speedtest_exporter

Prometheus exporter wrapping the Ookla `speedtest` CLI (runs it per scrape, parses `--format=json`).

## Layout
- `main.go` — flags, HTTP server, registers collectors.
- `exporter/` — `prometheus.Collector` that calls the client on Collect.
- `speedtest/` — CLI invocation (`Client.Args`, `NetworkMetrics`), JSON `Result`, metric descriptors. Tests in `speedtest_test.go`.

## Commands (verified)
- Dev shell: `nix develop` (go, gopls, golangci-lint, ookla-speedtest); `.envrc` = `use flake`.
- Test/vet/lint: `go test ./... && go vet ./... && golangci-lint run ./...`
- Nix package: `nix build` (buildGoModule, wraps ookla-speedtest on PATH). Bump `vendorHash` in `flake.nix` after any go.mod/go.sum change (set `lib.fakeHash`, build, copy `got:`).
- Docker: `docker build .`; multi-arch `docker buildx build --platform linux/amd64,linux/arm64 .` (needs a buildx/containerd builder, not the plain docker driver). Ookla CLI version pinned via `SPEEDTEST_VERSION` build arg.
- CI: `.github/workflows/ci.yml` (vet, test -race, golangci-lint, buildx amd64+arm64).

## Invariants / gotchas
- Module path is `github.com/kylegrantlucas/speedtest_exporter` (underscore, matches repo). Don't change.
- Ookla reports bandwidth in bytes/s; exported as Mbps (`*8/1e6`). Metric names are public API.
- `ookla-speedtest` is unfree; the flake allows only that package.
- No secrets in this repo.
