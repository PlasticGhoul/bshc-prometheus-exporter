FROM alpine:latest AS runner

ARG ARTIFACT_NAME
ARG ARTIFACT_VERSION

ARG TARGETOS
ARG TARGETARCH

RUN mkdir /app \
    mkdir /app/config

WORKDIR /app

RUN chown -R 10000:10000 /app

COPY --chmod=754 --chown=10000:10000 ./${ARTIFACT_NAME}-${ARTIFACT_VERSION}_${TARGETOS}-${TARGETARCH} /app/bshc-prometheus-exporter

CMD ["/app/bshc-prometheus-exporter", "-c", "/app/config/config.yaml" ]