#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. $DIR/../../config.sh

aedeploy gcloud app deploy --project ${PROJECT_ID}
