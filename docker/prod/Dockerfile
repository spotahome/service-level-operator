FROM golang:1.11-alpine AS build-stage

RUN apk --no-cache add \
    g++ \
    git \
    make

ARG VERSION
ENV VERSION=${VERSION}
WORKDIR /src
COPY . .
RUN ./hack/scripts/build-binary.sh

# Final image.
FROM alpine:latest
RUN apk --no-cache add \
  ca-certificates
COPY --from=build-stage /src/bin/service-level-operator /usr/local/bin/service-level-operator
ENTRYPOINT ["/usr/local/bin/service-level-operator"]
