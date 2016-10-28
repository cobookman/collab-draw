#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. $DIR/../../config.sh

echo "Testing dev"
HTTP_ADDR="localhost:8080"
GCLOUD_DATASET_ID=${PROJECT_ID} IS_DEV="1" bash -c "go test --httpAddr=$HTTP_ADDR"

echo "Testing project id"
HTTP_ADDR="$PROJECT_ID.appspot.com"
GCLOUD_DATASET_ID=${PROJECT_ID} IS_DEV="0" bash -c "go test --httpAddr=$HTTP_ADDR"
