package http

//go:generate protoc -I./proto -I. -I/home/vtolstov/.cache/go-path/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.1.0  -I/home/vtolstov/.cache/go-path/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.1.0/third_party/googleapis --go-grpc_out=paths=source_relative:./proto --go_out=paths=source_relative:./proto --micro_out=components=micro|http,debug=true,paths=source_relative:./proto proto/github.proto
