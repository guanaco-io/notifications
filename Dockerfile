FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

RUN mv notifications /usr/bin/

LABEL org.opencontainers.image.source="https://github.com/guanaco-io/notifications"

ENTRYPOINT ["notifications", "/etc/notifications/config.yml"]