FROM --platform=$BUILDPLATFORM golang:1.19.4 AS builder

COPY . /app
WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build

FROM gcr.io/distroless/static

COPY --from=builder /app/backup-google /

VOLUME /files

WORKDIR /files

ENTRYPOINT ["/backup-google"]
