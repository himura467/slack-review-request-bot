#!/usr/bin/env bash

set -e

ROOT_DIR=$(cd $(dirname $0)/..; pwd)

if [[ ! -f .env ]]; then
  echo ".env not found"
  exit 1
fi

source .env

docker build \
  --no-cache \
  --provenance=false \
  --progress=plain \
  --platform=linux/amd64 \
  --build-arg SLACK_OAUTH_TOKEN=$SLACK_OAUTH_TOKEN \
  --build-arg SLACK_SIGNING_SECRET=$SLACK_SIGNING_SECRET \
  --build-arg SLACK_REVIEWER_IDS=$SLACK_REVIEWER_IDS \
  -f Dockerfile -t slack-review-request-bot .
