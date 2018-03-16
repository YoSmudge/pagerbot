FROM golang:1.9-alpine as golang
WORKDIR /go/src/github.com/yosmudge/pagerbot
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags '-w -s' -a -installsuffix cgo -o main

# we can't use the scratch container b/c ranch requires sh/sleep

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
# the program:
COPY --from=golang /go/src/github.com/yosmudge/pagerbot/main /main
# the config:
COPY --from=golang /go/src/github.com/yosmudge/pagerbot/config.yml /config.yml

CMD ["/main"]
