FROM alpine:latest AS runner

ARG ARTIFACT_NAME

RUN mkdir /app \
    mkdir /app/config

WORKDIR /app

COPY --chmod=777 ./bin/${ARTIFACT_NAME} /app/bshc-prometheus-exporter

CMD ["/app/bshc-prometheus-exporter", "-c", "/app/config/config.yaml" ]