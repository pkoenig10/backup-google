FROM --platform=$BUILDPLATFORM golang:1.22.3 AS builder

COPY . /app
WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build

FROM gcr.io/distroless/static:latest@sha256:41972110a1c1a5c0b6adb283e8aa092c43c31f7c5d79b8656fbffff2c3e61f05

COPY --from=builder /app/backup-google /

VOLUME /files

WORKDIR /files

ENTRYPOINT ["/backup-google"]
