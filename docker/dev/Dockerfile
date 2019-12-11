FROM golang:1.13-alpine

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

# Fill go mod cache.
RUN mkdir /tmp/cache
COPY go.mod /tmp/cache
COPY go.sum /tmp/cache
RUN cd /tmp/cache && \
    go mod download

WORKDIR /src
