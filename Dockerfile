###############################
# Builder container
###############################

FROM golang:1.11.3-alpine3.8 AS builder
ENV GO111MODULE=on

WORKDIR /go/src/github.com/crowdworks/slacts
COPY . .

RUN set -x \
  && apk add --no-cache git="2.18.1-r0" build-base="0.5-r1" \
  && go mod download \
  && go mod verify \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo  -o ./slacts ./cmd/slacts/main.go

###############################
# Exec container
###############################

FROM alpine:3.8

ENV APP_DIR /usr/src/app

RUN set -x \
    && apk add --no-cache ca-certificates="20171114-r3" \
    && adduser -S slacts \
    && echo "slacts:slacts" | chpasswd \
    && addgroup -S slacts \
    && addgroup slacts slacts \
    && mkdir -p $APP_DIR

COPY --from=builder /go/src/github.com/crowdworks/slacts/slacts /usr/bin
RUN chown -R slacts:slacts $APP_DIR

USER slacts
WORKDIR $APP_DIR

ENTRYPOINT ["slacts"]
CMD ["slacts", "--help"]
