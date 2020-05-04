FROM golang:1.13 AS builder

ENV GO111MODULE on
ENV GOPRIVATE github.com/AlpacaLabs

ARG GITHUB_USER
ARG GITHUB_PASS

COPY go.mod go.sum /go/app/
WORKDIR /go/app
RUN go mod download

COPY . /go/app
RUN CGO_ENABLED=0 go build .

WORKDIR /go/app

FROM alpine:latest as app

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.0 \
 && wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
 && chmod +x /bin/grpc_health_probe

COPY --from=builder /go/app/app /app/app

RUN apk add --no-cache ca-certificates

WORKDIR /app
CMD ["./app"]