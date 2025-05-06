ARG GO_VERSION=1.23
FROM golang:${GO_VERSION}-bullseye as builder
# Build essentials

RUN apt update && apt install build-essential zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libsqlite3-dev libreadline-dev libffi-dev curl wget libbz2-dev -y

WORKDIR /opt/k8sload-exporter

ADD go.mod go.sum /opt/k8sload-exporter/

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod \
    --mount=type=cache,mode=0755,target=/root/.cache/go-build \
    go mod download

ADD . /opt/k8sload-exporter

RUN go mod tidy && go build -o k8sload-exporter


FROM debian:bookworm-slim
WORKDIR /opt/k8sload-exporter
COPY --from=builder /opt/k8sload-exporter /opt/k8sload-exporter

ENTRYPOINT ["/opt/k8sload-exporter/k8sload-exporter"]