FROM golang:1.14-alpine as build-stage

RUN apk --no-cache add \
  g++ \
  git \
  make \
  bash

ARG VERSION
ENV VERSION=${VERSION}

WORKDIR /src
COPY . .
RUN ./scripts/build/build.sh

# Final image.
FROM alpine:latest
RUN apk --no-cache add \
  ca-certificates
COPY --from=build-stage /src/bin/service-level-operator /usr/local/bin/service-level-operator
ENTRYPOINT ["/usr/local/bin/service-level-operator"]