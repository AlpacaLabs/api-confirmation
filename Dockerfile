FROM golang:1.13 AS builder

ENV GO111MODULE on
ENV GOPRIVATE github.com/AlpacaLabs

ARG GITHUB_USER
ARG GITHUB_PASS

COPY go.mod go.sum /go/app/
WORKDIR /go/app

# 1. configure git to use ssh instead of https
# 2. download dependencies
RUN git config --global url.git@github.com:.insteadOf https://github.com/ && \
    go mod download

COPY . /go/app

CMD ["go", "run", "main.go"]