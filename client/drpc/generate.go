package drpc

//xxgo:generate protoc -I./proto -I. --go_out=paths=source_relative:./proto --go-micro_out=components=micro|rpc,debug=true,paths=source_relative:./proto proto/test.proto

//go:generate go install storj.io/drpc/cmd/protoc-gen-go-drpc
//go:generate protoc -I./proto -I. --go_out=paths=source_relative:./proto --go-micro_out=components=micro|rpc,debug=true,paths=source_relative:./proto --go-drpc_out=json=false,paths=source_relative:./proto proto/test.proto
