#!/usr/bin/env sh

set -o errexit
set -o nounset

goos=linux
goarch=amd64
src=./cmd/service-level-operator
out=./bin/service-level-operator
ldf_cmp="-w -extldflags '-static'"
f_ver="-X main.Version=${VERSION:-dev}"

echo "Building binary at ${out}"

GOOS=${goos} GOARCH=${goarch} CGO_ENABLED=0 go build -o ${out} --ldflags "${ldf_cmp} ${f_ver}"  ${src}