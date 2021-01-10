package grpc

//go:generate protoc -I./proto -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.9.5/third_party/googleapis  -I. --go-grpc_out=paths=source_relative:./proto --go_out=paths=source_relative:./proto --micro_out=paths=source_relative:./proto proto/test.proto
