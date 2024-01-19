# syntax=docker/dockerfile:1

ARG VERSION=dev

## Builder
FROM golang:1.21.5-bookworm AS build
WORKDIR /app
ADD . .
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-X main.Version=$VERSION" -o hummingbird cli/main.go 

## Final Image
FROM alpine:3.19.0
WORKDIR /app/hummingbird
COPY --from=build /app/hummingbird .