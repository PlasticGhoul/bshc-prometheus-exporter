FROM golang:alpine AS builder

ARG ARTIFACT_NAME

WORKDIR /usr/src/app

COPY . .

RUN cd src \
    go mod download \
    go mod verify

RUN cd src \
	go build -v -o ./bin/ ./...
	
	
FROM alpine:latest AS runner

RUN mkdir /app \
    mkdir /app/config

WORKDIR /app

COPY --from=builder --chmod=777 /usr/src/app/src/bin/$ARTIFACT_NAME ./bshc-prometheus-exporter

CMD ["/app/bshc-prometheus-exporter", "-c", "/app/config/config.yaml" ]