#!/usr/bin/env sh

set -o errexit
set -o nounset

go test -race -v -tags='integration' ./...