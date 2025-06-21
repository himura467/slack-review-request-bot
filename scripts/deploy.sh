#!/usr/bin/env bash

set -e

if [[ $# -ne 2 ]]; then
  echo "Usage: $0 <aws-profile> <gcloud-config-name>"
  exit 1
fi

export AWS_PROFILE=$1
export CLOUDSDK_ACTIVE_CONFIG_NAME=$2

ROOT_DIR=$(cd "$(dirname "$0")"/..; pwd)

cd "$ROOT_DIR/terraform"

terraform init
terraform apply
