# Stage 1: build Go binary
FROM golang:1.24 AS builder

WORKDIR /tmp/builder
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY ./cmd ./cmd
COPY ./pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(go env GOARCH) go build -ldflags="-s -w" ./cmd/main.go


## Stage 2: copy the binary to alpine image
FROM alpine:3.21.3
RUN apk update && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=builder /tmp/builder/main /app/pokedex

ENTRYPOINT [ "/app/pokedex"]