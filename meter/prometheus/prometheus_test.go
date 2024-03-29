package prometheus_test

import (
	"bytes"
	"context"
	"testing"

	prometheus "go.unistack.org/micro-meter-prometheus/v3"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/meter/wrapper"
)

func TestWrapper(t *testing.T) {
	m := prometheus.NewMeter() // meter.Labels("test_key", "test_val"))

	w := wrapper.NewClientWrapper(
		wrapper.ServiceName("svc1"),
		wrapper.ServiceVersion("0.0.1"),
		wrapper.ServiceID("12345"),
		wrapper.Meter(m),
	)

	ctx := context.Background()

	c := client.NewClient(client.Wrap(w))
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	rsp := &codec.Frame{}
	req := &codec.Frame{}
	err := c.Call(ctx, c.NewRequest("svc2", "Service.Method", req), rsp)
	_, _ = rsp, err
	buf := bytes.NewBuffer(nil)
	_ = m.Write(buf, meter.WriteProcessMetrics(false), meter.WriteFDMetrics(false))
	if !bytes.Contains(buf.Bytes(), []byte(`micro_client_request_inflight{micro_endpoint="svc2.Service.Method"} 0`)) {
		t.Fatalf("invalid metrics output: %s", buf.Bytes())
	}
}
