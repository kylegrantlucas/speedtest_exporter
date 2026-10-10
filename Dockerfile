# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.26-bookworm AS builder
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/speedtest_exporter .

# Fetch the official Ookla CLI tarball (static binary) for the target arch.
FROM --platform=$BUILDPLATFORM debian:bookworm-slim AS speedtest
ARG TARGETARCH
ARG SPEEDTEST_VERSION=1.2.0
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*
RUN case "$TARGETARCH" in \
      amd64) arch=x86_64 ;; \
      arm64) arch=aarch64 ;; \
      arm) arch=armhf ;; \
      *) echo "unsupported arch: $TARGETARCH" >&2; exit 1 ;; \
    esac \
    && curl -fsSL "https://install.speedtest.net/app/cli/ookla-speedtest-${SPEEDTEST_VERSION}-linux-${arch}.tgz" \
       | tar -xz -C /usr/local/bin speedtest

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=speedtest /usr/local/bin/speedtest /usr/local/bin/speedtest
COPY --from=builder /out/speedtest_exporter /usr/local/bin/speedtest_exporter
EXPOSE 9112
ENTRYPOINT ["/usr/local/bin/speedtest_exporter"]
