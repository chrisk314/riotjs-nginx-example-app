#!/usr/bin/env bash

set -euo pipefail

REGISTRY=gcr.io/arcane-storm-147316
IMAGE_BASENAME=riotjs-nginx-example
GITBRANCH=$(git rev-parse --abbrev-ref HEAD)
GITHASH=$(git rev-parse --short HEAD)

while getopts "s:pl" OPTION; do
  case "${OPTION}" in
    s) SERVICES=("${OPTARG}");;
    p) PUSH=true;;
    l) LATEST=true;;
  esac
done

DEFAULT_SERVICES=("backend" "frontend")
SERVICES=("${SERVICES[@]-${DEFAULT_SERVICES[@]}}")
PUSH=${PUSH:-false}
LATEST=${LATEST:-false}

build_and_push () {
  SERVICE=${1}
  IMAGE=${IMAGE_BASENAME}-${SERVICE}
  REG_IMAGE=${REGISTRY}/${IMAGE}
  docker build ${SERVICE} \
    -t ${IMAGE}:${GITBRANCH} -t ${IMAGE}:${GITHASH} \
    -t ${REG_IMAGE}:${GITBRANCH} -t ${REG_IMAGE}:${GITHASH}
  echo "Successfully built image ${REG_IMAGE}"
  if ( ${LATEST} ); then
    docker tag ${IMAGE}:${GITHASH} ${IMAGE}:latest
    docker tag ${IMAGE}:${GITHASH} ${REG_IMAGE}:latest
    echo "Successfully added tag latest"
  fi
  if [[ "${SERVICE}" == "frontend" ]]; then
    docker build --target dev-server ${SERVICE} \
      -t ${IMAGE}-dev:${GITBRANCH} -t ${IMAGE}-dev:${GITHASH} -t ${IMAGE}-dev:latest
  fi
  if ( ${PUSH} ); then
    docker push ${REG_IMAGE}
    echo "Successfully pushed image ${REG_IMAGE}"
  fi
}

for service in "${SERVICES[@]}"; do
  build_and_push ${service}
done
