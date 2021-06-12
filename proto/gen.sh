#!/usr/bin/env bash

PROTO_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO_ROOT=${PROTO_ROOT}/../

function gen_go_proto() {
    protoc  -I${PROTO_ROOT} $1 \
        --go_out=plugins=grpc:./tmp/ \
        --oneofmapper_out=./tmp/
}

function gen_go_proto_validate() {
    protoc -I${PROTO_ROOT} $1 \
        --go_out=plugins=grpc:./tmp/ \
        --validate_out="lang=go:./tmp/"
}

mkdir -p ${REPO_ROOT}tmp/
gen_go_proto "common.proto"
gen_go_proto "messages.proto"
gen_go_proto "shard_events.proto"
gen_go_proto "envelope.proto"

cp -rf ${REPO_ROOT}tmp/github.com/alexykot/cncraft/* ${REPO_ROOT}
rm -rf ${REPO_ROOT}tmp
