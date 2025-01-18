package grpc

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go_out=paths=source_relative:./proto --go-micro_out=components='micro|grpc',debug=true,paths=source_relative:./proto proto/test.proto"
