# Basic Controller

Minimal Kubernetest controller example/showcase. Includes local pipeline commands (lint/test/build), Docker image, and a quick k8s (k3s via k3d) local environment.

The controller is monitoring changes on the `greeting` CRD objects in the cluster and look for the related `configMap` to update its content with the `greeting`'s `message` field.

## Docker

```
# Run linter
docker run --rm \
    -v "$(pwd)":/controller \
    -w /controller \
    golangci/golangci-lint:v2.4.0-alpine \
    golangci-lint run ./...

# Run test
docker run --rm \
    -v "$(pwd)":/controller \
    -w /controller \
    golang:1.25.0-alpine \
    go test -v -coverprofile=out.cover ./...

# Build local image
docker build -t controller:dev .

# Run container
docker run --rm -p 8080:8080 controller:dev
```

### Dockerfile

[BuildKit](https://github.com/moby/buildkit) replaces the legacy Docker builder. BuildKit is the default builder for users on Docker Desktop, and Docker Engine as of version 23.0. Enabled with `DOCKER_BUILDKIT=1` if that's not the case.

The Dockerfile provides syntax specific for BuildKit:

- `--mount=type=cache` allows you to mount a cache volume during the build step that can persist between builds, e.g.

```
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
```

the `go mod download` command uses a cache mount for `/go/pkg/mod` and `/root/.cache/go-build` directories. The cache mount is persisted across builds, so even if you end up rebuilding the layer, you only download new or changed packages.

The used image for run the controller is `gcr.io/distroless/static:nonroot`:

- [Distroless](https://github.com/GoogleContainerTools/distroless) images contain only your application and its runtime dependencies. This makes them much smaller and much harder to exploit (by an attacker).
- The `static` variant is designed for statically compiled binaries (like your Go binary with `CGO_ENABLED=0`).
- The`nonroot`tag means the image is set up with a predefined non-root user.

### go build

`go build` flags:

- `-trimpath`: Remove all file system paths from the resulting executable. Hence, reduces binary size and prevents leaking local paths.
- `-ldflags`: Arguments to pass on each go tool link invocation.

and environment variables:

- `CGO_ENABLED=0`: To not support the cgo command and produce a static binary.

### go tool link

`go tool link` flags:

- `-s`: Omit the symbol table and debug information. Implies the `-w` flag, which can be negated with `-w=0`. Hence, the binary is smaller.
- `-w`: Omit the DWARF symbol table. Hence, the binary is smaller.

## Notes

- About CRDs, see the [official documentation](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/).