FROM golang:1.13 AS builder

ENV GO111MODULE on
ENV GOPRIVATE github.com/AlpacaLabs

ARG GITHUB_USER
ARG GITHUB_PASS

COPY go.mod go.sum /go/app/
WORKDIR /go/app

RUN go mod download

COPY . /go/app

CMD ["go", "run", "main.go"]