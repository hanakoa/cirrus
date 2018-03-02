FROM golang:1.10.0 as builder
ADD . /go/src/github.com/kevinmichaelchen/cirrus
WORKDIR /go/src/github.com/kevinmichaelchen/cirrus
RUN go get ./... && \
    CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/cirrus .

FROM alpine:latest
#ENV CIRRUS_USER="cirrus" \
#    CIRRUS_UID="8080" \
#    CIRRUS_GROUP="cirrus" \
#    CIRRUS_GID="8080"
RUN apk --no-cache add ca-certificates
#&& \
#    addgroup -S -g $CIRRUS_GID $CIRRUS_GROUP && \
#    adduser -S -u $CIRRUS_UID -G $CIRRUS_GROUP $CIRRUS_USER
WORKDIR /root/
COPY --from=builder /go/src/github.com/kevinmichaelchen/cirrus/bin/cirrus .
#USER $CIRRUS_USER
CMD ["./cirrus"]