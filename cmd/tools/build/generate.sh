#!/usr/bin/env bash

declare -r target="${1:-all}"

case $target in
    proto)
      ./proto/gen.sh
      ;;
    go)
    	go generate ./...
      ;;
    sql)
      ./core/db/gen.sh
      ;;
    all|*)
      ./proto/gen.sh
      go generate ./...
      ./core/db/gen.sh
      ;;
esac
