package grpc

//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
//go:generate sh -c "protoc -I./proto -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) -I. --go-grpc_out=paths=source_relative:./proto --go_out=paths=source_relative:./proto --go-micro_out=components='micro|grpc',debug=true,standalone=true,paths=source_relative:./gproto proto/test.proto"
