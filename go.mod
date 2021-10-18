module github.com/unistack-org/micro-tests

go 1.16

require (
	github.com/VictoriaMetrics/metrics v1.17.3 // indirect
	github.com/armon/go-metrics v0.3.9 // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/hashicorp/consul/api v1.9.1 // indirect
	github.com/hashicorp/go-hclog v0.16.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/miekg/dns v1.1.38 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/segmentio/kafka-go v0.4.21
	github.com/stretchr/testify v1.7.0
	github.com/twmb/franz-go v1.1.1
	github.com/unistack-org/micro-broker-kgo/v3 v3.7.3
	github.com/unistack-org/micro-broker-segmentio/v3 v3.4.3-0.20210804134048-7916dafb4dfe
	//github.com/unistack-org/micro-client-drpc/v3 v3.0.0-00010101000000-000000000000
	github.com/unistack-org/micro-client-grpc/v3 v3.7.8
	github.com/unistack-org/micro-client-http/v3 v3.7.3
	github.com/unistack-org/micro-codec-grpc/v3 v3.7.4
	github.com/unistack-org/micro-codec-json/v3 v3.7.8
	github.com/unistack-org/micro-codec-jsonpb/v3 v3.7.5
	github.com/unistack-org/micro-codec-proto/v3 v3.7.4
	github.com/unistack-org/micro-codec-segmentio/v3 v3.7.5
	github.com/unistack-org/micro-codec-urlencode/v3 v3.7.4
	github.com/unistack-org/micro-codec-xml/v3 v3.7.5
	github.com/unistack-org/micro-config-consul/v3 v3.6.0
	github.com/unistack-org/micro-config-env/v3 v3.5.0
	github.com/unistack-org/micro-config-vault/v3 v3.5.0
	github.com/unistack-org/micro-meter-victoriametrics/v3 v3.3.3
	github.com/unistack-org/micro-proto v0.0.9
	github.com/unistack-org/micro-router-register/v3 v3.7.8
	github.com/unistack-org/micro-server-grpc/v3 v3.7.4
	github.com/unistack-org/micro-server-http/v3 v3.7.6
	github.com/unistack-org/micro-server-tcp/v3 v3.3.2
	github.com/unistack-org/micro-wrapper-recovery/v3 v3.0.0
	github.com/unistack-org/micro-wrapper-trace-opentracing/v3 v3.3.0
	github.com/unistack-org/micro/v3 v3.7.8
	golang.org/x/net v0.0.0-20210929193557-e81a3d93ecf6 // indirect
	golang.org/x/sys v0.0.0-20210930212924-f542c8878de8 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210930144712-2e2e1008e8a3
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	storj.io/drpc v0.0.26
)

//replace github.com/unistack-org/micro-wrapper-trace-opentracing/v3 => ../micro-wrapper-trace-opentracing
//replace github.com/unistack-org/micro-client-grpc/v3 => ../micro-client-grpc
//replace github.com/unistack-org/micro-server-grpc/v3 => ../micro-server-grpc
//replace github.com/unistack-org/micro-server-http/v3 => ../micro-server-http
//replace github.com/unistack-org/micro-client-http/v3 => ../micro-client-http
//replace github.com/unistack-org/micro-client-drpc/v3 => ../micro-client-drpc
//replace github.com/unistack-org/micro-broker-segmentio/v3 => ../micro-broker-segmentio
//replace github.com/unistack-org/micro-proto => ../micro-proto
