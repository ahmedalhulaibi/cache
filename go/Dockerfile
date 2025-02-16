#syntax=docker/dockerfile:1.2
ARG PKG
ARG BUILDPKG

FROM golang:1.24 AS builder

RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

# Force fetching modules over SSH
RUN git config --system url."ssh://git@github.com/".insteadOf "https://github.com/"

ARG PKG
ARG BUILDPKG
WORKDIR /go/src/$PKG

# Setup goprivate to fetch private dependencies
ARG GOPRIVATE
ENV GOPRIVATE=$GOPRIVATE

COPY . .

RUN --mount=type=ssh \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=linux \
    go build -o app $PKG/$BUILDPKG

FROM gcr.io/distroless/static AS final

ENV ADDR ":8080"
ENV GRPC_ADDR ":8090"

WORKDIR /

ARG PKG
COPY --from=builder /go/src/$PKG/app app

CMD ["/app"]