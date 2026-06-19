# Build the manager binary.
#
# NOTE: the build context is the repository ROOT, not this directory.
# The operator module is linked to the shared ./apis module through the
# repository-root go.work file, so the workspace and apis source must be
# present in the build:
#   docker build -f operators/account-operator/Dockerfile .
FROM --platform=$BUILDPLATFORM golang:1.26.4-bookworm@sha256:5f68ec6805843bd3981a951ffada82a26a0bd2631045c8f7dba483fa868f5ec5 AS builder
ARG TARGETARCH

WORKDIR /workspace

# Go workspace + the modules this image needs.
COPY go.work go.work
COPY go.work.sum go.work.sum
COPY apis/ apis/
COPY operators/account-operator/ operators/account-operator/

# Drop sibling workspace modules so the workspace resolves with only what we copied.
RUN go work edit -dropuse ./operators/security-operator

WORKDIR /workspace/operators/account-operator

# cache deps before building so that source changes don't invalidate the layer
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -ldflags '-w -s' -o /workspace/manager main.go

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV USER_UID=1001
ENV GROUP_UID=1001
COPY --from=builder --chown=${USER_UID}:${GROUP_UID}  /workspace/manager /operator/manager

USER ${USER_UID}:${GROUP_UID}
ENTRYPOINT ["/operator/manager"]