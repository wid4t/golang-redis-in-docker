FROM golang:alpine3.18 as builder
WORKDIR /golang-redis-in-docker
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" .
FROM busybox
WORKDIR /golang-redis-in-docker
COPY --from=builder /golang-redis-in-docker /usr/bin/
ENTRYPOINT ["golang-redis-in-docker"]