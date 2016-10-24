#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
. $DIR/config.sh

# test each service
for d in $DIR/services/* ; do
  cd $d
  ./test.sh
done
