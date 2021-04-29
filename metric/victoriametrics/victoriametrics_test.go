// +build ignore

package victoriametrics

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	metrics "github.com/VictoriaMetrics/metrics"
	"github.com/stretchr/testify/assert"
	gclient "github.com/unistack-org/micro-client-grpc/v3"
	rrouter "github.com/unistack-org/micro-router-register/v3"
	gserver "github.com/unistack-org/micro-server-grpc/v3"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/client"
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
		return fmt.Errorf("victoriametrics test error")
	}
	return nil
}

func TestVictoriametrics(t *testing.T) {
	// setup
	r := register.NewRegister()
	b := broker.NewBroker(broker.Register(r))

	name := "test"
	id := "id-1234567890"
	version := "1.2.3.4"

	c := gclient.NewClient(
		client.Router(rrouter.NewRouter(router.Register(r))),
		client.Broker(b),
	)
	s := gserver.NewServer(
		server.Broker(b),
		server.Name(name),
		server.Version(version),
		server.Id(id),
		server.Register(r),
		server.WrapHandler(
			NewHandlerWrapper(
				ServiceName(name),
				ServiceVersion(version),
				ServiceID(id),
			),
		),
	)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Stop()

	type Test struct {
		*testHandler
	}

	s.Handle(
		s.NewHandler(&Test{new(testHandler)}),
	)

	if err := s.Start(); err != nil {
		t.Fatalf("Unexpected error starting server: %v", err)
	}

	req := c.NewRequest(name, "Test.Method", &TestRequest{IsError: false}, client.WithContentType("application/json"))
	rsp := TestResponse{}

	assert.NoError(t, c.Call(context.TODO(), req, &rsp))

	req = c.NewRequest(name, "Test.Method", &TestRequest{IsError: true}, client.WithContentType("application/json"))
	assert.Error(t, c.Call(context.TODO(), req, &rsp))

	buf := bytes.NewBuffer(nil)
	metrics.WritePrometheus(buf, false)

	metric, err := findMetricByName(buf, "sum", "micro_server_request_total")
	if err != nil {
		t.Fatal(err)
	}

	labels := metric[0]["labels"].(map[string]string)
	for k, v := range labels {
		switch k {
		case "micro_version":
			assert.Equal(t, version, v)
		case "micro_id":
			assert.Equal(t, id, v)
		case "micro_name":
			assert.Equal(t, name, v)
		case "micro_endpoint":
			assert.Equal(t, "Test.Method", v)
		case "micro_status":
			continue
		default:
			t.Fatalf("unknown %v with %v", k, v)
		}
	}
}

func findMetricByName(buf io.Reader, tp string, name string) ([]map[string]interface{}, error) {
	var metrics []map[string]interface{}
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.HasPrefix(txt, name) {
			mt := make(map[string]interface{})
			v := txt[strings.LastIndex(txt, " "):]
			k := ""
			if idx := strings.Index(txt, "{"); idx > 0 {
				labels := make(map[string]string)
				lb := strings.Split(txt[idx+1:strings.Index(txt, "}")], ",")
				for _, l := range lb {
					p := strings.Split(l, "=")
					labels[strings.Trim(p[0], `"`)] = strings.Trim(p[1], `"`)
				}
				mt["labels"] = labels
				k = txt[:idx]
			} else {
				k = txt[:strings.Index(txt, " ")]
			}
			mt["name"] = k
			mt["value"] = v
			metrics = append(metrics, mt)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(metrics) == 0 {
		return nil, fmt.Errorf("%s %s not found", tp, name)
	}
	return metrics, nil
}
