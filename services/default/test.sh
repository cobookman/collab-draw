#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. $DIR/../../config.sh

GCLOUD_PROJECT_ID=$PROJECT_ID go test
