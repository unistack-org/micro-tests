module github.com/unistack-org/micro-tests

go 1.16

require (
	github.com/opentracing/opentracing-go v1.2.0
	github.com/segmentio/kafka-go v0.4.17
	github.com/stretchr/testify v1.7.0
	github.com/unistack-org/micro-api-handler-rpc/v3 v3.3.0
	github.com/unistack-org/micro-api-router-register/v3 v3.2.2
	github.com/unistack-org/micro-api-router-static/v3 v3.2.1
	github.com/unistack-org/micro-broker-segmentio/v3 v3.4.1
	//github.com/unistack-org/micro-client-drpc/v3 v3.0.0-00010101000000-000000000000
	github.com/unistack-org/micro-client-grpc/v3 v3.4.0
	github.com/unistack-org/micro-client-http/v3 v3.4.8
	github.com/unistack-org/micro-codec-grpc/v3 v3.2.1
	github.com/unistack-org/micro-codec-json/v3 v3.2.1
	github.com/unistack-org/micro-codec-jsonpb/v3 v3.2.2
	github.com/unistack-org/micro-codec-proto/v3 v3.2.2
	github.com/unistack-org/micro-codec-segmentio/v3 v3.2.2
	github.com/unistack-org/micro-codec-urlencode/v3 v3.0.0
	github.com/unistack-org/micro-codec-xml/v3 v3.2.2
	github.com/unistack-org/micro-config-consul/v3 v3.5.1
	github.com/unistack-org/micro-config-env/v3 v3.4.0
	github.com/unistack-org/micro-config-vault/v3 v3.4.0
	github.com/unistack-org/micro-meter-victoriametrics/v3 v3.3.3
	github.com/unistack-org/micro-proto v0.0.2
	github.com/unistack-org/micro-router-register/v3 v3.2.2
	github.com/unistack-org/micro-server-grpc/v3 v3.3.6
	github.com/unistack-org/micro-server-http/v3 v3.4.1
	github.com/unistack-org/micro-server-tcp/v3 v3.3.2
	github.com/unistack-org/micro-wrapper-trace-opentracing/v3 v3.2.0
	github.com/unistack-org/micro/v3 v3.5.9
	google.golang.org/genproto v0.0.0-20210729151513-df9385d47c1b
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
	storj.io/drpc v0.0.24
)

//replace github.com/unistack-org/micro-wrapper-trace-opentracing/v3 => ../micro-wrapper-trace-opentracing
//replace github.com/unistack-org/micro-client-grpc/v3 => ../micro-client-grpc
//replace github.com/unistack-org/micro-server-grpc/v3 => ../micro-server-grpc
//replace github.com/unistack-org/micro-server-http/v3 => ../micro-server-http
//replace github.com/unistack-org/micro-client-http/v3 => ../micro-client-http
//replace github.com/unistack-org/micro-client-drpc/v3 => ../micro-client-drpc
replace github.com/unistack-org/micro-broker-segmentio/v3 => ../micro-broker-segmentio

//replace github.com/unistack-org/micro-proto => ../micro-proto
