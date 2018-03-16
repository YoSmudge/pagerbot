FROM golang:1.9-alpine as golang
WORKDIR /go/src/github.com/yosmudge/pagerbot
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags '-w -s' -a -installsuffix cgo -o main

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM scratch
# the program:
COPY --from=golang /go/src/github.com/yosmudge/pagerbot/main /main
# the config:
COPY --from=golang /go/src/github.com/yosmudge/pagerbot/config.yml /config.yml
# the timezone data:
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/main"]
