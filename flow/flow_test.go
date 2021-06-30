package flow

import (
	"context"
	"testing"
	"time"

	pb "github.com/unistack-org/micro-tests/client/http/proto"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/flow"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/store"
	"github.com/unistack-org/micro/v3/tracer"
)

func TestFlow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := store.DefaultStore
	c := client.DefaultClient
	m := meter.DefaultMeter
	tr := tracer.DefaultTracer
	l := logger.DefaultLogger

	f := flow.NewFlow(flow.Context(ctx), flow.Store(s), flow.Client(c), flow.Meter(m), flow.Tracer(tr), flow.Logger(l))

	if err := f.Init(); err != nil {
		t.Fatal(err)
	}

	steps := []flow.Step{
		flow.NewCallStep("test", "Github.LookupUser", flow.StepID("test.Github.LookupUser")),
		flow.NewCallStep("test", "Github.UpdateUser", flow.StepRequires("test.Github.LookupUser")),
		flow.NewCallStep("test", "Github.RemoveUser", flow.StepRequires("test.Github.UpdateUser")),
		flow.NewCallStep("test", "Github.MailUser", flow.StepRequires("test.Github.UpdateUser")),
	}
	w, err := f.WorkflowCreate(ctx, "test", steps...)
	if err != nil {
		t.Fatal(err)
	}

	req := &pb.LookupUserReq{Username: "vtolstov"}
	id, err := w.Execute(ctx, req, flow.ExecuteTimeout(2*time.Second))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("execution id: %s", id)
}
