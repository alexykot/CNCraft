#!/usr/bin/env bash

PROTO_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO_ROOT=${PROTO_ROOT}/../
GO_PROTO_PLUGIN="go_out"

function gen_go_proto() {
    protoc --${GO_PROTO_PLUGIN}=plugins=grpc:./tmp/ \
        -I${PROTO_ROOT} \
        $1
}

function gen_go_proto_validate() {
    protoc --${GO_PROTO_PLUGIN}=plugins=grpc:./tmp/ \
        -I${PROTO_ROOT} \
        --validate_out="lang=go:./tmp/" \
        $1
}

mkdir -p ${REPO_ROOT}tmp/
gen_go_proto "common.proto"
gen_go_proto "messages.proto"
gen_go_proto "shard_events.proto"
gen_go_proto "envelope.proto"

cp -rf ${REPO_ROOT}tmp/github.com/alexykot/cncraft/* ${REPO_ROOT}
rm -rf ${REPO_ROOT}tmp
