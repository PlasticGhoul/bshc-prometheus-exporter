FROM golang:latest as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

FROM scratch:latest as runner

RUN mkdir /app \
    mkdir /app/config

WORKDIR /app

COPY --from=builder /usr/local/bin/app/bshc-prometheus-exporter .

CMD ["/app/bshc-prometheus-exporter", "-c", "/app/config/config.yaml" ]