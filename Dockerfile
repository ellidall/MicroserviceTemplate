FROM golang:1.23-bookworm AS builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends --no-install-suggests \
        ca-certificates \
        tzdata && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean && \
    groupadd -g 1001 microuser && \
    useradd -u 1001 -r -g 1001 -s /sbin/nologin -c "go microservice user" microuser

ADD ./bin/microservicetemplate /app/bin/
ADD ./data /app/data
WORKDIR /app

USER microuser
# Run the application command by default when the container starts.
ENTRYPOINT [ "/app/bin/microservicetemplate" ]
# Pass `service` as argument to ENTRYPOINT if no argument passed on container run
CMD ["service"]
