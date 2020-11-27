module github.com/unistack-org/micro-tests

go 1.13

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/client_model v0.2.0
	github.com/stretchr/testify v1.6.1
	github.com/unistack-org/micro-broker-http v0.0.0-20201125231853-bb4bd204b8c0
	github.com/unistack-org/micro-broker-memory v0.0.2-0.20201105185131-5ff932308afd
	github.com/unistack-org/micro-client-grpc v0.0.2-0.20201125224558-067cf68d2312
	github.com/unistack-org/micro-client-http v0.0.0-20201125231021-64a08cf7fd5d
	github.com/unistack-org/micro-codec-grpc v0.0.0-20201126055537-b2e5c1ec2168
	github.com/unistack-org/micro-codec-json v0.0.0-20201125092251-38cf770f2eb0
	github.com/unistack-org/micro-codec-proto v0.0.0-20201125092414-f627fea89e7e
	github.com/unistack-org/micro-codec-segmentio v0.0.0-20201127144339-03740d564751
	github.com/unistack-org/micro-metrics-prometheus v0.0.2-0.20201125232532-93104a0ff374
	github.com/unistack-org/micro-registry-memory v0.0.2-0.20201105195351-bd57ee0e4bd6
	github.com/unistack-org/micro-router-registry v0.0.2-0.20201105175056-773128885d9e
	github.com/unistack-org/micro-server-grpc v0.0.3-0.20201125221721-36040a57659a
	github.com/unistack-org/micro-server-http v0.0.2-0.20201125222045-54ee918b278c
	github.com/unistack-org/micro-server-tcp v0.0.2-0.20201125222121-31fd93a07671
	github.com/unistack-org/micro-wrapper-opentracing v0.0.1
	github.com/unistack-org/micro/v3 v3.0.2-0.20201125221305-0d93b2c31c79
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0
)

//replace (
//	github.com/unistack-org/micro-client-grpc => ../done/micro-client-grpc
//	github.com/unistack-org/micro-client-http => ../done/micro-client-http
//	github.com/unistack-org/micro-server-grpc => ../done/micro-server-grpc
//	github.com/unistack-org/micro-server-http => ../done/micro-server-http
//	github.com/unistack-org/micro/v3 => ../micro
//)
