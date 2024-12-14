FROM golang:alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./bin/ ./...

FROM alpine:latest AS runner

RUN mkdir /app \
    mkdir /app/config

WORKDIR /app

COPY --from=builder --chmod=777 /usr/src/app/bin/bshc-prometheus-exporter ./bshc-prometheus-exporter

CMD ["/app/bshc-prometheus-exporter", "-c", "/app/config/config.yaml" ]