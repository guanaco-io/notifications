# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD .. /go/src/github.com/guanaco-io/notifications

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/guanaco-io/notifications
RUN go install github.com/guanaco-io/notifications

COPY ../config/config.yaml /etc/notifications/config.yaml

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/github.com/guanaco-io/notifications /etc/notifications/config.yaml
