#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. $DIR/../../config.sh

echo "Testing dev"
WS_ADDR="localhost:65080"
GCLOUD_DATASET_ID=${PROJECT_ID} IS_DEV="1" bash -c "go test --wsAddr=$WS_ADDR"

