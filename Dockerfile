FROM --platform=$BUILDPLATFORM golang:1.23.5 AS builder

COPY . /app
WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build

FROM gcr.io/distroless/static:latest@sha256:3f2b64ef97bd285e36132c684e6b2ae8f2723293d09aae046196cca64251acac

COPY --from=builder /app/backup-google /

VOLUME /files

WORKDIR /files

ENTRYPOINT ["/backup-google"]
