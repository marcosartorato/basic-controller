# syntax=docker/dockerfile:1.7

############################
# Stage 1: deps (cache go mod)
############################

FROM golang:1.25-alpine AS deps
WORKDIR /src

COPY go.mod go.sum ./
# Cache modules & build cache for faster rebuilds
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    go mod download

############################
# Stage 2: build
############################

FROM golang:1.25-alpine AS builder
WORKDIR /src

# Copy source
COPY . .

# Build the binary
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    go build -trimpath \
      -ldflags="-s -w" \
      -o /bin/greeting-manager ./cmd/greeting-manager

############################
# Stage 3: runtime
############################

FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /bin/greeting-manager /greeting-manager
USER nonroot:nonroot

# Expose HTTP port for metrics/healthz
EXPOSE 8080
EXPOSE 8081

# Run manager binary
CMD ["/greeting-manager"]