package prometheus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	cli "github.com/unistack-org/micro-client-grpc/v3"
	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	promwrapper "github.com/unistack-org/micro-metrics-prometheus/v3"
	rrouter "github.com/unistack-org/micro-router-register/v3"
	srv "github.com/unistack-org/micro-server-grpc/v3"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/register"
	"github.com/unistack-org/micro/v3/router"
	"github.com/unistack-org/micro/v3/server"
)

type Test interface {
	Method(ctx context.Context, in *TestRequest, opts ...client.CallOption) (*TestResponse, error)
}

type TestRequest struct {
	IsError bool
}
type TestResponse struct{}

type testHandler struct{}

func (t *testHandler) Method(ctx context.Context, req *TestRequest, rsp *TestResponse) error {
	if req.IsError {
		return fmt.Errorf("test error")
	}
	return nil
}

func TestPrometheusMetrics(t *testing.T) {
	client.DefaultRetries = 0
	// setup
	reg := register.NewRegister()
	brk := broker.NewBroker(broker.Register(reg))

	name := "test"
	id := "id-1234567890"
	version := "1.2.3.4"
	rt := rrouter.NewRouter(router.Register(reg))

	c := cli.NewClient(
		client.Codec("application/grpc+json", jsoncodec.NewCodec()),
		client.Codec("application/json", jsoncodec.NewCodec()),
		client.Router(rt),
	)
	s := srv.NewServer(
		server.Codec("application/grpc+json", jsoncodec.NewCodec()),
		server.Codec("application/json", jsoncodec.NewCodec()),
		server.Name(name),
		server.Version(version),
		server.Id(id),
		server.Register(reg),
		server.Broker(brk),
		server.WrapHandler(
			promwrapper.NewHandlerWrapper(
				promwrapper.ServiceName(name),
				promwrapper.ServiceVersion(version),
				promwrapper.ServiceID(id),
			),
		),
	)

	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	type Test struct {
		*testHandler
	}

	s.Handle(
		s.NewHandler(&Test{new(testHandler)}),
	)

	if err := s.Start(); err != nil {
		t.Fatalf("Unexpected error starting server: %v", err)
	}
	defer s.Stop()

	req := c.NewRequest(name, "Test.Method", &TestRequest{IsError: false}, client.RequestContentType("application/json"))
	rsp := TestResponse{}

	assert.NoError(t, c.Call(context.TODO(), req, &rsp))

	req = c.NewRequest(name, "Test.Method", &TestRequest{IsError: true}, client.RequestContentType("application/json"))
	assert.Error(t, c.Call(context.TODO(), req, &rsp))

	list, _ := prometheus.DefaultGatherer.Gather()

	metric := findMetricByName(list, dto.MetricType_SUMMARY, "micro_server_latency_microseconds")

	if metric == nil || metric.Metric == nil || len(metric.Metric) == 0 {
		t.Fatalf("no metrics returned")
	}

	for _, v := range metric.Metric[0].Label {
		switch *v.Name {
		case "micro_version":
			assert.Equal(t, version, *v.Value)
		case "micro_id":
			assert.Equal(t, id, *v.Value)
		case "micro_name":
			assert.Equal(t, name, *v.Value)
		case "micro_endpoint":
			assert.Equal(t, "Test.Method", *v.Value)
		default:
			t.Fatalf("unknown %v with %v", *v.Name, *v.Value)
		}
	}

	assert.Equal(t, uint64(2), *metric.Metric[0].Summary.SampleCount)
	assert.True(t, *metric.Metric[0].Summary.SampleSum > 0)

	metric = findMetricByName(list, dto.MetricType_HISTOGRAM, "micro_server_request_duration_seconds")

	for _, v := range metric.Metric[0].Label {
		switch *v.Name {
		case "micro_version":
			assert.Equal(t, version, *v.Value)
		case "micro_id":
			assert.Equal(t, id, *v.Value)
		case "micro_name":
			assert.Equal(t, name, *v.Value)
		case "micro_endpoint":
			assert.Equal(t, "Test.Method", *v.Value)
		default:
			t.Fatalf("unknown %v with %v", *v.Name, *v.Value)
		}
	}

	assert.Equal(t, uint64(2), *metric.Metric[0].Histogram.SampleCount)
	assert.True(t, *metric.Metric[0].Histogram.SampleSum > 0)

	metric = findMetricByName(list, dto.MetricType_COUNTER, "micro_server_request_total")

	for _, v := range metric.Metric[0].Label {
		switch *v.Name {
		case "micro_version":
			assert.Equal(t, version, *v.Value)
		case "micro_id":
			assert.Equal(t, id, *v.Value)
		case "micro_name":
			assert.Equal(t, name, *v.Value)
		case "micro_endpoint":
			assert.Equal(t, "Test.Method", *v.Value)
		case "micro_status":
			assert.Equal(t, "failure", *v.Value)
		}
	}
	assert.Equal(t, *metric.Metric[0].Counter.Value, float64(1))

	for _, v := range metric.Metric[1].Label {
		switch *v.Name {
		case "micro_version":
			assert.Equal(t, version, *v.Value)
		case "micro_id":
			assert.Equal(t, id, *v.Value)
		case "micro_name":
			assert.Equal(t, name, *v.Value)
		case "micro_endpoint":
			assert.Equal(t, "Test.Method", *v.Value)
		case "micro_status":
			assert.Equal(t, "success", *v.Value)
		}
	}

	assert.Equal(t, *metric.Metric[1].Counter.Value, float64(1))
}

func findMetricByName(list []*dto.MetricFamily, tp dto.MetricType, name string) *dto.MetricFamily {
	for _, metric := range list {
		if *metric.Name == name && *metric.Type == tp {
			return metric
		}
	}

	return nil
}
