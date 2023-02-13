# builder image
FROM golang:alpine as builder
USER root
RUN mkdir -m 777 /MikroTik-Speedtest/
COPY MikroTik-Speedtest/. /MikroTik-Speedtest/
WORKDIR /MikroTik-Speedtest/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -extldflags=-static" -a -o speedtest .

# generate clean, final image for end users
FROM cdhtlr/busybox

LABEL maintainer="Sidik Hadi Kurniadi" name="Mikrotik Speedtest" description="Base minimum MikroTik-Terminal friendly Download Speedtest" version="4.5"

ENV URL="https://jakarta.speedtest.telkom.net.id.prod.hosts.ooklaserver.net:8080/download?size=25000000"
ENV MAX_DLSIZE="1.0"
ENV MIN_THRESHOLD="1.0"
ENV CONCURENT_CONNECTION="4"
ENV ALLOW_MEMORY_BUFFER="YES"

COPY --from=builder /MikroTik-Speedtest/speedtest .

EXPOSE 80

CMD ["/speedtest"]
