#!/bin/sh -e

INC=$(go list -f '{{ .Dir }}' -m github.com/unistack-org/micro-proto)
INC_CODEC=$(go list -f '{{ .Dir }}' -m github.com/unistack-org/micro/v3)
ARGS="-I${INC}"
CODEC_ARGS="-I${INC_CODEC}"

protoc $ARGS $CODEC_ARGS -Iproto --openapiv2_out=disable_default_errors=true,allow_merge=true:./proto/ --go_out=paths=source_relative:./proto/ proto/*.proto

protoc $ARGS $CODEC_ARGS -Iproto --go-micro_out=components="micro|http",debug=true,tag_path=./proto/,paths=source_relative:./proto/ proto/*.proto
