# syntax=docker/dockerfile:1

FROM aquasec/trivy:0.76.2 AS trivy

FROM --platform=$BUILDPLATFORM golang:1.27-alpine@sha256:70bd7ddea8b6cbe5d6943cde279d6b67feb10cc924fb9100d1d8757c1575b280 AS build
ARG TARGETARCH
ARG VERSION=dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /vordr ./cmd/vordr

FROM alpine:3.24
RUN apk add --no-cache ca-certificates
LABEL org.opencontainers.image.title="vordr" \
      org.opencontainers.image.description="Lightweight CVE scanner for homelab Docker stacks" \
      org.opencontainers.image.source="https://github.com/st0o0/vordr" \
      org.opencontainers.image.documentation="https://github.com/st0o0/vordr#readme" \
      org.opencontainers.image.licenses="MIT"
COPY --from=trivy /usr/local/bin/trivy /usr/local/bin/trivy
COPY --from=build /vordr /vordr
COPY LICENSE NOTICE /

HEALTHCHECK --interval=30s --timeout=10s --start-period=45s --retries=3 \
  CMD ["/vordr", "healthcheck"]

ENTRYPOINT ["/vordr"]
