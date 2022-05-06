# builder image
FROM golang:1.18.1-alpine3.15 as builder
USER root
RUN mkdir -m 777 /speedtest-go/
COPY speedtest-go/. /speedtest-go/
WORKDIR /speedtest-go/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o speedtest-cli .

# generate clean, final image for end users
FROM alpine:3.15

LABEL maintainer="Sidik Hadi Kurniadi" name="Mikrotik Speedtest" description="Base minimum MikroTik-Terminal friendly Ookla Speedtest.net" version="2.0"

USER root

ENV PING_IP=""
ENV PING_DOMAIN=""
ENV TZ=""
ENV CRON_FIELD=""

RUN apk --no-cache add php8 php8-cli tzdata curl sudo && \
    printf "Build of cdhtlr/mikrotik-speedtest, date: %s\n"  `date -u +"%Y-%m-%dT%H:%M:%SZ"` >> /etc/BUILD && \
    rm -rf /var/cache/apk/*

COPY --from=builder /speedtest-go/speedtest-cli .

COPY webpage/. /var/www/localhost/htdocs/

ADD entrypoint.sh /

RUN mkdir -m 777 /var/www/localhost/htdocs/results/data/
RUN chmod 755 /entrypoint.sh

WORKDIR /var/www/localhost/

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
