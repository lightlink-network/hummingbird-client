# syntax=docker/dockerfile:1

## Builder
FROM golang:1.24.4-bookworm AS build
WORKDIR /app
ADD . .
ARG VERSION=dev
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-X main.Version=$VERSION" -o hummingbird cli/hb/main.go 

## Final Image
FROM alpine:3.19.0
WORKDIR /app/hummingbird
COPY --from=build /app/hummingbird .