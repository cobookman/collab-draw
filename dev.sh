#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. $DIR/config.sh

cd $DIR

# Build go binary, quit if this fails
# This step isn't necessary, but its useful for
# debugging to fail fast here vs deal with parsing
# docker output when the go build tool would have more
# useful output
for d in $DIR/services/* ; do
  cd $d
  go build
  if [ $? -ne 0 ]; then
    echo "Failed to build go binary"
    exit
  fi
done

# for each service...
for d in $DIR/services/* ; do
  echo "Starting up service: $service"
  $d/dev.sh
done
