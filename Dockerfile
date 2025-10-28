FROM golang:1.23-bookworm AS builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends --no-install-suggests \
        ca-certificates \
        tzdata && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean && \
    groupadd -g 1001 appuser && \
    useradd -u 1001 -r -g 1001 -s /sbin/nologin -c "app user" appuser

ADD ./bin/microservicetemplate /app/bin/
WORKDIR /app

ENTRYPOINT [ "/app/bin/microservicetemplate" ]

CMD ["service"]
