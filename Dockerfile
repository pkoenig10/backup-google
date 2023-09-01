FROM --platform=$BUILDPLATFORM golang:1.20.6 AS builder

COPY . /app
WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build

FROM gcr.io/distroless/static:latest@sha256:e7e79fb2947f38ce0fab6061733f7e1959c12b843079042fe13f56ca7b9d178c

COPY --from=builder /app/backup-google /

VOLUME /files

WORKDIR /files

ENTRYPOINT ["/backup-google"]
