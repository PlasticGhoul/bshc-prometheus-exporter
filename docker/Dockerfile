FROM alpine:latest AS runner

ARG ARTIFACT_NAME

LABEL org.opencontainers.image.title="BSHC Prometheus Exporter"
LABEL org.opencontainers.image.description="Bosch Smart Home Controller Prometheus Exporter"
LABEL org.opencontainers.image.authors="PlasticGhoul"

RUN mkdir /app \
    mkdir /app/config

WORKDIR /app

COPY --chmod=777 ./bin/${ARTIFACT_NAME} /app/bshc-prometheus-exporter

CMD ["/app/bshc-prometheus-exporter", "-c", "/app/config/config.yaml" ]