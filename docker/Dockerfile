FROM golang:alpine AS builder

WORKDIR /usr/src/app

COPY . .

RUN cd src \
    go mod download \
    go mod verify

RUN cd src \
	go build -v -o ../bin/ ./... \
	ls \
	ls bin

FROM alpine:latest AS runner

RUN mkdir /app \
    mkdir /app/config

WORKDIR /app

COPY --from=builder --chmod=777 /usr/src/app/bin/bshc-prometheus-exporter ./bshc-prometheus-exporter

CMD ["/app/bshc-prometheus-exporter", "-c", "/app/config/config.yaml" ]