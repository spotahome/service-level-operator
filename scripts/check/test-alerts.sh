#!/usr/bin/env sh

set -o errexit
set -o nounset

DIR="$(dirname "$(readlink -f $0)")"
ROOT_DIR="${DIR}/../.."
ALERTS=${ROOT_DIR}/alerts

promtool test rules ${ALERTS}/*_test.yaml