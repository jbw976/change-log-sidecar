FROM alpine:latest

WORKDIR /app

COPY change-log-sidecar .

CMD ["./change-log-sidecar"]