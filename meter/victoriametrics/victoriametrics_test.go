package victoriametrics_test

import (
	"bytes"
	"context"
	"testing"

	victoriametrics "go.unistack.org/micro-meter-victoriametrics/v3"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/meter"
)

func TestWrapper(t *testing.T) {
	t.Skip()
	m := victoriametrics.NewMeter() // meter.Labels("test_key", "test_val"))
	_ = m.Init()

	ctx := context.Background()

	c := client.NewClient()
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	rsp := &codec.Frame{}
	req := &codec.Frame{}
	err := c.Call(ctx, c.NewRequest("svc2", "Service.Method", req), rsp)
	_, _ = rsp, err
	buf := bytes.NewBuffer(nil)
	_ = m.Write(buf, meter.WriteProcessMetrics(false), meter.WriteFDMetrics(false))
	if !bytes.Contains(buf.Bytes(), []byte(`micro_client_request_inflight{endpoint="svc2.Service.Method"} 0`)) {
		t.Fatalf("invalid metrics output: %s", buf.Bytes())
	}
}
