#!/bin/bash

set -e
set -u
set -o pipefail

case ${1:-} in
  build)
    docker build -t clanker-jail .
    ;;
  run)
    cd ..
    docker run -it --rm \
      -v "$(pwd):/home/clanker/workspace:z" \
      -e "DEEPSEEK_API_KEY=$DEEPSEEK_API_KEY" \
      clanker-jail
    ;;
  *)
    echo "Usage: ./docker.sh [build|run]"
    exit 1
    ;;
esac
