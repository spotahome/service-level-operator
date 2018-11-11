#!/bin/sh

DIR="$(dirname "$(readlink -f $0)")"
ROOT_DIR="${DIR}/../.."
ALERTS=${ROOT_DIR}/alerts

promtool test rules ${ALERTS}/*_test.yaml