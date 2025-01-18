package combo

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go_out=paths=source_relative:./proto proto/proto.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components='micro',standalone=false,debug=true,paths=source_relative:./proto proto/proto.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components='grpc',standalone=true,debug=true,paths=source_relative:./mgpb proto/proto.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components='http',standalone=true,debug=true,paths=source_relative:./mhpb proto/proto.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components='openapiv3',openapi_file=./apidocs.swagger.yaml,standalone=true,debug=true,paths=source_relative:./proto proto/proto.proto"

//go:generate sh -c "protoc -I./ngpb -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go_out=paths=source_relative:./ngpb ngpb/ngpb.proto"
//go:generate sh -c "protoc -I./ngpb -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-grpc_out=paths=source_relative:./ngpb ngpb/ngpb.proto"

////go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go_out=paths=source_relative:./ndpb proto/proto.proto"
////go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-drpc_out=json=false,paths=source_relative:./ndpb proto/proto.proto"

////go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v3) --go-micro_out=components='micro|drpc',standalone=true,debug=true,paths=source_relative:./mdpb proto/proto.proto"

//go:generate sh -c "mkdir -p swagger-ui && cp proto/apidocs.swagger.yaml swagger-ui/swagger.yaml && curl -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v4.17.0.tar.gz -o v4.17.0.tar.gz && tar -C swagger-ui --strip-components=2 -zxvf v4.17.0.tar.gz swagger-ui-4.17.0/dist && rm -f v4.17.0.tar.gz && sed -i '' 's|https://petstore.swagger.io/v2/swagger.json|./swagger.yaml|g' swagger-ui/index.html swagger-ui/swagger-initializer.js && sed -i '' 's|deepLinking: true,|deepLinking: true, displayOperationId: true, tryItOutEnabled: true,|g' swagger-ui/swagger-initializer.js "
