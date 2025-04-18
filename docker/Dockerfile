FROM alpine:latest AS runner

ARG ARTIFACT_PATH
ARG ARTIFACT_NAME

ENV ARTIFACT_PATH=${ARTIFACT_PATH}
ENV ARTIFACT_NAME=${ARTIFACT_NAME}

RUN mkdir /app \
    mkdir /app/config

WORKDIR /app

COPY --chmod=777 $ARTIFACT_PATH/$ARTIFACT_NAME /app/bshc-prometheus-exporter

CMD ["/app/bshc-prometheus-exporter", "-c", "/app/config/config.yaml" ]