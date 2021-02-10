module github.com/unistack-org/micro-tests

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.2.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/client_model v0.2.0
	github.com/stretchr/testify v1.7.0
	github.com/unistack-org/micro-api-handler-rpc/v3 v3.1.1
	github.com/unistack-org/micro-api-router-register/v3 v3.2.2
	github.com/unistack-org/micro-api-router-static/v3 v3.2.1
	github.com/unistack-org/micro-broker-http/v3 v3.2.2
	github.com/unistack-org/micro-broker-memory/v3 v3.2.3
	github.com/unistack-org/micro-client-grpc/v3 v3.2.3
	github.com/unistack-org/micro-client-http/v3 v3.2.3
	github.com/unistack-org/micro-codec-grpc/v3 v3.1.1
	github.com/unistack-org/micro-codec-json/v3 v3.1.1
	github.com/unistack-org/micro-codec-jsonpb/v3 v3.1.1
	github.com/unistack-org/micro-codec-proto/v3 v3.1.1
	github.com/unistack-org/micro-codec-segmentio/v3 v3.1.1
	github.com/unistack-org/micro-config-env/v3 v3.2.4
	github.com/unistack-org/micro-config-vault/v3 v3.2.4
	github.com/unistack-org/micro-metrics-prometheus/v3 v3.1.1
	github.com/unistack-org/micro-register-memory/v3 v3.2.2
	github.com/unistack-org/micro-router-register/v3 v3.2.2
	github.com/unistack-org/micro-server-grpc/v3 v3.2.2
	github.com/unistack-org/micro-server-http/v3 v3.2.4
	github.com/unistack-org/micro-server-tcp/v3 v3.2.2
	github.com/unistack-org/micro-wrapper-trace-opentracing/v3 v3.1.1
	github.com/unistack-org/micro/v3 v3.2.8
	google.golang.org/genproto v0.0.0-20210207032614-bba0dbe2a9ea
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
)

//replace github.com/unistack-org/micro-server-http/v3 => ../micro-server-http

//replace github.com/unistack-org/micro-client-http/v3 => ../micro-client-http

//replace github.com/unistack-org/micro/v3 => ../micro
