package combo

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go_out=paths=source_relative:./proto proto/test.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components='micro|http',standalone=true,debug=true,paths=source_relative:./mhpb proto/test.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components='micro|grpc',standalone=true,debug=true,paths=source_relative:./mgpb proto/test.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components='micro|drpc',standalone=true,debug=true,paths=source_relative:./mdpb proto/test.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go_out=paths=source_relative:./ngpb proto/test.proto"
//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-grpc_out=paths=source_relative:./ngpb proto/test.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go_out=paths=source_relative:./ndpb proto/test.proto"
//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-drpc_out=json=false,paths=source_relative:./ndpb proto/test.proto"
