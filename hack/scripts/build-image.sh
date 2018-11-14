#!/usr/bin/env sh

set -e

if [ -z ${VERSION} ]; then
    echo "VERSION env var needs to be set"
    exit 1
fi

REPOSITORY="quay.io/spotahome/"
IMAGE="service-level-operator"
TARGET_IMAGE=${REPOSITORY}${IMAGE}


docker build \
    --build-arg VERSION=${VERSION} \
    -t ${TARGET_IMAGE}:${VERSION} \
    -t ${TARGET_IMAGE}:latest \
    -f ./docker/prod/Dockerfile .

if [ -n "${PUSH_IMAGE}" ]; then
    echo "pushing ${TARGET_IMAGE} images..."
    docker push ${TARGET_IMAGE}
fi