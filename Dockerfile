# builder image
FROM golang:alpine as builder
USER root
RUN mkdir -m 777 /speedtest-go/
COPY speedtest-go/. /speedtest-go/
WORKDIR /speedtest-go/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o speedtest .

# generate clean, final image for end users
FROM alpine

LABEL maintainer="Sidik Hadi Kurniadi" name="Mikrotik Speedtest" description="Base minimum MikroTik-Terminal friendly Download Speedtest" version="3.0"

USER root

ENV MAX_KB="1000"
ENV THRESHOLD_MBPS="1.0"
ENV URL="https://jakarta.speedtest.telkom.net.id.prod.hosts.ooklaserver.net:8080/download?size=25000000"

COPY --from=builder /speedtest-go/speedtest .

EXPOSE 80

CMD ./speedtest -m ${MAX_KB} -t ${THRESHOLD_MBPS} -u ${URL}