FROM golang:1.13 AS builder

ENV GO111MODULE on
ENV GOPRIVATE github.com/AlpacaLabs

ARG GITHUB_USER
ARG GITHUB_PASS

COPY go.mod go.sum /go/app/
WORKDIR /go/app

# 1. add credentials on build
# 2. make sure your domain is accepted
# 3. configure git to use ssh instead of https
# 4. download dependencies
ARG SSH_PRIVATE_KEY
RUN mkdir /root/.ssh/ && \
    echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa && \
    chmod 600 /root/.ssh/id_rsa && \
    touch /root/.ssh/known_hosts && \
    ssh-keyscan github.com >> /root/.ssh/known_hosts && \
    git config --global url.git@github.com:.insteadOf https://github.com/ && \
    go mod download

COPY . /go/app
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:latest as app

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.0 \
 && wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
 && chmod +x /bin/grpc_health_probe

COPY --from=builder /go/app/app /app/app

RUN apk add --no-cache ca-certificates

WORKDIR /app
CMD ["./app"]