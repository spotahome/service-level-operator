#!/usr/bin/env sh

set -o errexit
set -o nounset

IMAGE_CLI_GEN=quay.io/slok/kube-code-generator:v1.17.3
IMAGE_CRD_GEN=quay.io/slok/kube-code-generator:v1.18.0
ROOT_DIRECTORY=$(dirname "$(readlink -f "$0")")/../
PROJECT_PACKAGE="github.com/spotahome/service-level-operator"

echo "Generating Kubernetes CRD clients..."
docker run -it --rm \
	-v ${ROOT_DIRECTORY}:/go/src/${PROJECT_PACKAGE} \
	-e PROJECT_PACKAGE=${PROJECT_PACKAGE} \
	-e CLIENT_GENERATOR_OUT=${PROJECT_PACKAGE}/pkg/kubernetes/gen \
	-e APIS_ROOT=${PROJECT_PACKAGE}/pkg/apis \
	-e GROUPS_VERSION="monitoring:v1alpha1" \
	-e GENERATION_TARGETS="deepcopy,client" \
	${IMAGE_CLI_GEN}

echo "Generating Kubernetes CRD manifests..."
docker run -it --rm \
	-v ${ROOT_DIRECTORY}:/src \
	-e GO_PROJECT_ROOT=/src \
	-e CRD_TYPES_PATH=/src/pkg/apis \
	-e CRD_OUT_PATH=/src/manifests/crd \
	${IMAGE_CRD_GEN} update-crd.sh