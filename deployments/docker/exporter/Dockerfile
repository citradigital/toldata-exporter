FROM alpine

ENV NATS_URL ${NATS_URL}

COPY deployments/docker/build-results/exporter /
COPY deployments/docker/exporter/start.sh /

RUN apk add --update bash

RUN chmod +x /start.sh 
ENTRYPOINT [ "/start.sh" ]
