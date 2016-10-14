#!/bin/bash
PROJECT_ID="strong-moose"
WS_ADDR="localhost:65080"
HTTP_ADDR="localhost:8080"
GCLOUD_DATASET_ID=${PROJECT_ID} IS_DEV="1" bash -c "go test --wsAddr=${WS_ADDR} --httpAddr=${HTTP_ADDR}"

