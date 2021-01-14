module github.com/unistack-org/micro-tests

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/client_model v0.2.0
	github.com/stretchr/testify v1.6.1
	github.com/unistack-org/micro-api-handler-rpc v0.0.0-20210113163127-3f36e7b3d99a
	github.com/unistack-org/micro-api-router-registry v0.0.0-20210110113004-b9ccb4324370
	github.com/unistack-org/micro-api-router-static v0.0.0-20210110113147-58f8ed2f7347
	github.com/unistack-org/micro-broker-http v0.0.0-20201125231853-bb4bd204b8c0
	github.com/unistack-org/micro-broker-memory v0.0.2-0.20201105185131-5ff932308afd
	github.com/unistack-org/micro-client-grpc v0.0.2-0.20201228123319-bbd07bb0914a
	github.com/unistack-org/micro-client-http v0.0.0-20210110114810-59d77f7b8cf9
	github.com/unistack-org/micro-codec-grpc v0.0.0-20201220205513-cad30014cbf2
	github.com/unistack-org/micro-codec-json v0.0.0-20201220205604-ed33fab21d87
	github.com/unistack-org/micro-codec-proto v0.0.0-20201220205718-066176ab59b7
	github.com/unistack-org/micro-codec-segmentio v0.0.0-20201220210027-bc88e5dad1c2
	github.com/unistack-org/micro-config-env v0.0.0-20201219213431-afab7aa1d69f
	github.com/unistack-org/micro-metrics-prometheus v0.0.2-0.20201125232532-93104a0ff374
	github.com/unistack-org/micro-registry-memory v0.0.2-0.20210110004413-c27422bc489a
	github.com/unistack-org/micro-router-registry v0.0.2-0.20201105175056-773128885d9e
	github.com/unistack-org/micro-server-grpc v0.0.3-0.20201228125110-b2aa849c1e7b
	github.com/unistack-org/micro-server-http v0.0.2-0.20201125222045-54ee918b278c
	github.com/unistack-org/micro-server-tcp v0.0.2-0.20201125222121-31fd93a07671
	github.com/unistack-org/micro-wrapper-opentracing v0.0.1
	github.com/unistack-org/micro/v3 v3.0.2-0.20210110005504-270ad1b88914
	google.golang.org/genproto v0.0.0-20200904004341-0bd0a958aa1d
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
)

//replace github.com/unistack-org/micro-client-http => ../done/micro-client-http
//replace github.com/unistack-org/micro/v3 => ../micro
