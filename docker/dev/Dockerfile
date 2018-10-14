FROM golang:1.11-alpine

RUN apk --no-cache add \
    bash \
    git \
    g++ \
    curl \
    openssl \
    openssh-client

# Mock creator
RUN go get -u github.com/vektra/mockery/.../ 

RUN mkdir /src

# Create user
ARG uid=1000
ARG gid=1000
RUN addgroup -g $gid service-level-operator && \
    adduser -D -u $uid -G service-level-operator service-level-operator && \
    chown service-level-operator:service-level-operator -R /src && \
    chown service-level-operator:service-level-operator -R /go
USER service-level-operator

# Fill go mod cache.
RUN mkdir /tmp/cache
COPY go.mod /tmp/cache
COPY go.sum /tmp/cache
RUN cd /tmp/cache && \
    go mod download

WORKDIR /src
