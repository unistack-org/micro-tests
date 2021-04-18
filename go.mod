module github.com/unistack-org/micro-tests

go 1.16

require (
	github.com/google/uuid v1.2.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/client_model v0.2.0
	github.com/stretchr/testify v1.7.0
	github.com/unistack-org/micro-api-handler-rpc/v3 v3.3.0
	github.com/unistack-org/micro-api-router-register/v3 v3.2.2
	github.com/unistack-org/micro-api-router-static/v3 v3.2.1
	github.com/unistack-org/micro-broker-http/v3 v3.3.0
	github.com/unistack-org/micro-client-grpc/v3 v3.3.1
	github.com/unistack-org/micro-client-http/v3 v3.3.8
	github.com/unistack-org/micro-codec-grpc/v3 v3.2.1
	github.com/unistack-org/micro-codec-json/v3 v3.2.1
	github.com/unistack-org/micro-codec-jsonpb/v3 v3.2.2
	github.com/unistack-org/micro-codec-proto/v3 v3.2.2
	github.com/unistack-org/micro-codec-segmentio/v3 v3.2.2
	github.com/unistack-org/micro-codec-urlencode/v3 v3.0.0
	github.com/unistack-org/micro-config-env/v3 v3.2.5
	github.com/unistack-org/micro-config-vault/v3 v3.2.9
	github.com/unistack-org/micro-meter-victoriametrics/v3 v3.3.0
	github.com/unistack-org/micro-metrics-prometheus/v3 v3.1.1
	github.com/unistack-org/micro-proto v0.0.2-0.20210227213711-77c7563bd01e
	github.com/unistack-org/micro-router-register/v3 v3.2.2
	github.com/unistack-org/micro-server-grpc/v3 v3.3.2
	github.com/unistack-org/micro-server-http/v3 v3.3.11
	github.com/unistack-org/micro-server-tcp/v3 v3.3.0
	github.com/unistack-org/micro-wrapper-trace-opentracing/v3 v3.2.0
	github.com/unistack-org/micro/v3 v3.3.14
	google.golang.org/genproto v0.0.0-20210416161957-9910b6c460de
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)

//replace github.com/unistack-org/micro-wrapper-trace-opentracing/v3 => ../micro-wrapper-trace-opentracing
//replace github.com/unistack-org/micro-client-grpc/v3 => ../micro-client-grpc
//replace github.com/unistack-org/micro-server-grpc/v3 => ../micro-server-grpc
//replace github.com/unistack-org/micro-server-http/v3 => ../micro-server-http
//replace github.com/unistack-org/micro-client-http/v3 => ../micro-client-http
//replace github.com/unistack-org/micro/v3 => ../micro
//replace github.com/unistack-org/micro-proto => ../micro-proto
