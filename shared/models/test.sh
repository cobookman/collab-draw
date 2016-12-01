#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. $DIR/../../config.sh

# Setup
EMULATOR_PIDS=()

$(gcloud beta emulators datastore env-init)
gcloud beta emulators datastore start > /dev/null 2>&1 &
EMULATOR_PIDS+=($!)

# Disown simply stops us from getting messages when we kill the emulator.
# This must be the first command after the emulator start command
disown

$(gcloud beta emulators pubsub env-init)
gcloud beta emulators pubsub start > /dev/null 2>&1 &
EMULATOR_PIDS+=($!)
disown

# Unit test
GCLOUD_PROJECT_ID=$PROJECT_ID UPSTREAM_DRAWING_TOPIC=$UPSTREAM_DRAWING_TOPIC go test

# Tear down emulators
for i in "${!EMULATOR_PIDS[@]}"; do
  EMULATOR_PID=${EMULATOR_PIDS[$i]}
  kill -TERM $EMULATOR_PID > /dev/null 2>&1
  while [[ $(awk '$1=="process_id" {print $0}' <(top -n 1 -b)) ]]; do
    sleep 1
    kill -9 $EMULATOR_PID > /dev/null 2>&1
  done
done

gcloud beta emulators datastore env-unset | eval

