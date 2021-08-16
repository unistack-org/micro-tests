package grpc

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m github.com/unistack-org/micro-proto) --go_out=paths=source_relative:./proto --go-micro_out=components='micro|http',debug=true,paths=source_relative:./proto proto/test.proto"

