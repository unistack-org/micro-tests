package grpc

//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
//go:generate protoc -I./proto -I. --go-grpc_out=paths=source_relative:./proto --go_out=paths=source_relative:./proto --go-micro_out=components=micro|rpc,standalone=true,paths=source_relative:./gproto proto/test.proto
