###############################
# Builder container
###############################

FROM golang:1.17.6-alpine AS builder

WORKDIR /go/src/github.com/crowdworks/slacts
COPY . .

RUN set -x \
  && apk add --no-cache git build-base \
  && go mod download \
  && go mod verify \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo  -o ./slacts ./cmd/slacts/main.go

###############################
# Exec container
###############################

FROM alpine:3.15

ENV APP_DIR /usr/src/app

RUN set -x \
    && apk add --no-cache ca-certificates \
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
