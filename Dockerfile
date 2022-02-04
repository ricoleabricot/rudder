# Build the manager binary
FROM golang:1.17 as builder

# Copy everything not ignored to workspace
WORKDIR /workspace
COPY . /workspace

# Set arg from CI
ARG GIT_COMMIT
ARG BUILD_DATE

# Configure env
ENV CGO_ENABLED=0 \
    GOOS="linux" \
    GOARCH="amd64" \
    GO111MODULE="on" \
    GO_APP_PKG="github.com/ricoleabricot/rudder" \
    GO111MODULE=on

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Build the controller manager
RUN go build \
    -a -installsuffix cgo \
    -o rudder \
    -ldflags "-X ${GO_APP_PKG}.gitCommit=${GIT_COMMIT} -X ${GO_APP_PKG}.buildDate=${BUILD_DATE}"

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy the manager binary previously built
COPY --from=builder /workspace/rudder .
USER nonroot:nonroot

ENTRYPOINT ["/rudder"]
