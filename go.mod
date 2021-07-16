module github.com/unistack-org/micro-tests

go 1.16

require (
	github.com/frankban/quicktest v1.11.3 // indirect
	github.com/google/uuid v1.3.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/client_model v0.2.0
	github.com/stretchr/testify v1.7.0
	github.com/unistack-org/micro-api-handler-rpc/v3 v3.3.0
	github.com/unistack-org/micro-api-router-register/v3 v3.2.2
	github.com/unistack-org/micro-api-router-static/v3 v3.2.1
	github.com/unistack-org/micro-broker-http/v3 v3.3.1
	//github.com/unistack-org/micro-client-drpc/v3 v3.0.0-00010101000000-000000000000
	github.com/unistack-org/micro-client-grpc/v3 v3.3.3
	github.com/unistack-org/micro-client-http/v3 v3.4.5
	github.com/unistack-org/micro-codec-grpc/v3 v3.2.1
	github.com/unistack-org/micro-codec-json/v3 v3.2.1
	github.com/unistack-org/micro-codec-jsonpb/v3 v3.2.2
	github.com/unistack-org/micro-codec-proto/v3 v3.2.2
	github.com/unistack-org/micro-codec-segmentio/v3 v3.2.2
	github.com/unistack-org/micro-codec-urlencode/v3 v3.0.0
	github.com/unistack-org/micro-codec-xml/v3 v3.2.2
	github.com/unistack-org/micro-config-env/v3 v3.4.0
	github.com/unistack-org/micro-config-vault/v3 v3.4.0
	github.com/unistack-org/micro-meter-victoriametrics/v3 v3.3.1
	github.com/unistack-org/micro-metrics-prometheus/v3 v3.1.1
	github.com/unistack-org/micro-proto v0.0.2
	github.com/unistack-org/micro-router-register/v3 v3.2.2
	github.com/unistack-org/micro-server-grpc/v3 v3.3.6
	github.com/unistack-org/micro-server-http/v3 v3.4.1
	github.com/unistack-org/micro-server-tcp/v3 v3.3.2
	github.com/unistack-org/micro-wrapper-trace-opentracing/v3 v3.2.0
	github.com/unistack-org/micro/v3 v3.4.9
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210716133855-ce7ef5c701ea
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
//replace github.com/unistack-org/micro-broker-segmentio/v3 => ../micro-broker-segmentio

//replace github.com/unistack-org/micro/v3 => ../micro

//replace github.com/unistack-org/micro-proto => ../micro-proto
