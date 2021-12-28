module go.unistack.org/micro-tests

go 1.16

//replace go.unistack.org/micro-wrapper-sql/v3 => ../micro-wrapper-sql

require (
	github.com/jmoiron/sqlx v1.3.4
	github.com/opentracing/opentracing-go v1.2.0
	github.com/segmentio/kafka-go v0.4.25
	github.com/stretchr/testify v1.7.0
	github.com/twmb/franz-go v1.2.6
	go.unistack.org/micro-broker-kgo/v3 v3.8.2
	go.unistack.org/micro-broker-segmentio/v3 v3.8.0
	go.unistack.org/micro-client-grpc/v3 v3.8.1
	go.unistack.org/micro-client-http/v3 v3.8.3
	go.unistack.org/micro-codec-grpc/v3 v3.8.1
	go.unistack.org/micro-codec-json/v3 v3.8.0
	go.unistack.org/micro-codec-jsonpb/v3 v3.8.1
	go.unistack.org/micro-codec-proto/v3 v3.8.1
	go.unistack.org/micro-codec-segmentio/v3 v3.8.1
	go.unistack.org/micro-codec-urlencode/v3 v3.8.1
	go.unistack.org/micro-codec-xml/v3 v3.8.1
	go.unistack.org/micro-config-consul/v3 v3.8.1
	go.unistack.org/micro-config-env/v3 v3.8.2
	go.unistack.org/micro-config-vault/v3 v3.8.3
	go.unistack.org/micro-meter-prometheus/v3 v3.8.0
	go.unistack.org/micro-meter-victoriametrics/v3 v3.8.4
	go.unistack.org/micro-proto/v3 v3.1.0
	go.unistack.org/micro-router-register/v3 v3.8.1
	go.unistack.org/micro-server-grpc/v3 v3.8.0
	go.unistack.org/micro-server-http/v3 v3.8.1
	go.unistack.org/micro-server-tcp/v3 v3.8.0
	go.unistack.org/micro-wrapper-recovery/v3 v3.8.0
	go.unistack.org/micro-wrapper-sql/v3 v3.0.1
	go.unistack.org/micro-wrapper-trace-opentracing/v3 v3.8.0
	go.unistack.org/micro/v3 v3.8.13
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211223182754-3ac035c7e7cb // indirect
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	modernc.org/sqlite v1.14.3
	storj.io/drpc v0.0.26
)
