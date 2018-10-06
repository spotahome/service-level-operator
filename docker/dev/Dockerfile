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
    chown service-level-operator:service-level-operator -R /src

USER service-level-operator
WORKDIR /src
