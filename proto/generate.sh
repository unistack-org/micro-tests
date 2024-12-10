#!/bin/sh -e

INC=$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3)
INC_CODEC=$(go list -f '{{ .Dir }}' -m go.unistack.org/micro/v3)
ARGS="-I${INC}"
CODEC_ARGS="-I${INC_CODEC}"

protoc $ARGS $CODEC_ARGS -I. --go_out=paths=source_relative:./ *.proto

