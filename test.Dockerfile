#syntax=docker/dockerfile:1.2

FROM golang:1.24 AS builder

RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

# Force fetching modules over SSH
RUN git config --system url."ssh://git@github.com/".insteadOf "https://github.com/"

WORKDIR /go/src/github.com/ahmedalhulaibi/hello-world-api

# Setup goprivate to fetch ahmedalhulaibi dependencies
ENV GOPRIVATE="github.com/ahmedalhulaibi"

COPY . .

RUN --mount=type=ssh \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=linux \
    go test -v -cover ./...

