FROM golang:1.18.3 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

ADD . ${GOPATH}/src/app/
WORKDIR ${GOPATH}/src/app

RUN go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/speedtest_exporter

FROM ubuntu:focal

COPY --from=builder /go/bin/speedtest_exporter /usr/bin/speedtest_exporter

RUN apt-get update \
    && apt-get install -y curl 

RUN curl -s https://install.speedtest.net/app/cli/install.deb.sh | bash
RUN apt-get install speedtest

EXPOSE 9112

ENTRYPOINT [ "/usr/bin/speedtest_exporter" ]
