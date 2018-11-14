#!/usr/bin/env sh

DIR="$( cd "$( dirname "${0}" )" && pwd )"
ROOT_DIR=${DIR}/../..

PROJECT_PACKAGE=github.com/spotahome/service-level-operator
IMAGE=quay.io/slok/kube-code-generator:v1.11.3

# Execute once per package because we want independent output specs per kind/version.
docker run -it --rm \
    -v ${ROOT_DIR}:/go/src/${PROJECT_PACKAGE} \
    -e CRD_PACKAGES=${PROJECT_PACKAGE}/pkg/apis/monitoring/v1alpha1 \
    -e OPENAPI_OUTPUT_PACKAGE=${PROJECT_PACKAGE}/pkg/apis/monitoring/v1alpha1 \
    ${IMAGE} ./update-openapi.sh