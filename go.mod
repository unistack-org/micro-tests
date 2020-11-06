module github.com/unistack-org/micro-tests

go 1.13

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.5.1
	github.com/unistack-org/micro-broker-http v0.0.0-20201106084013-bff50fb8c334
	github.com/unistack-org/micro-broker-memory v0.0.2-0.20201105185131-5ff932308afd
	github.com/unistack-org/micro-client-grpc v0.0.2-0.20201028070730-15a5d7d2cde8
	github.com/unistack-org/micro-registry-memory v0.0.2-0.20201105195351-bd57ee0e4bd6
	github.com/unistack-org/micro-router-registry v0.0.2-0.20201105175056-773128885d9e
	github.com/unistack-org/micro-server-grpc v0.0.2-0.20201105204550-241e452ecf38
	github.com/unistack-org/micro-server-http v0.0.2-0.20201104225538-7d3dc63ae435
	github.com/unistack-org/micro-server-tcp v0.0.2-0.20201104231236-b12d45f45cbc
	github.com/unistack-org/micro-wrapper-opentracing v0.0.1
	github.com/unistack-org/micro/v3 v3.0.0-gamma.0.20201106081812-be8d09c66352
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0
)

//replace github.com/unistack-org/micro/v3 => ../micro
//replace github.com/unistack-org/micro-registry-memory => ../done/micro-registry-memory
