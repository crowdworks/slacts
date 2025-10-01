###############################
# Builder container
###############################

FROM golang:1.25.1-trixie AS builder

WORKDIR /go/src/github.com/crowdworks/slacts

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make install

###############################
# Exec container
###############################

FROM debian:trixie-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

ENV APP_DIR /usr/src/app

RUN set -x \
    && useradd -s /bin/bash slacts \
    && mkdir -p $APP_DIR \
    && chown -R slacts:slacts $APP_DIR

COPY --from=builder /go/bin/slacts /usr/bin

USER slacts
WORKDIR $APP_DIR

ENTRYPOINT ["slacts"]
CMD ["--help"]
