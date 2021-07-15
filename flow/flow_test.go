package flow

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	httpcli "github.com/unistack-org/micro-client-http/v3"
	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	jsonpbcodec "github.com/unistack-org/micro-codec-jsonpb/v3"
	httpsrv "github.com/unistack-org/micro-server-http/v3"
	pb "github.com/unistack-org/micro-tests/flow/proto"
	"github.com/unistack-org/micro/v3"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/flow"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/server"
	"github.com/unistack-org/micro/v3/store"
	"github.com/unistack-org/micro/v3/tracer"
)

type handler struct{}

func (h *handler) DeleteUser(ctx context.Context, req *pb.DeleteUserReq, rsp *pb.DeleteUserRsp) error {
	return nil
}

func (h *handler) UpdateUser(ctx context.Context, req *pb.UpdateUserReq, rsp *pb.UpdateUserRsp) error {
	return nil
}

func (h *handler) MailUser(ctx context.Context, req *pb.MailUserReq, rsp *pb.MailUserRsp) error {
	return nil
}

func (h *handler) LookupUser(ctx context.Context, req *pb.LookupUserReq, rsp *pb.LookupUserRsp) error {
	rsp.Birthday = "31.07.1986"
	return nil
}

func TestFlow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.DefaultLogger = logger.NewLogger(logger.WithLevel(logger.DebugLevel))

	s := store.DefaultStore
	if err := s.Init(store.Codec(jsoncodec.NewCodec())); err != nil {
		t.Fatal(err)
	}
	if err := s.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	c := client.NewClientCallOptions(
		httpcli.NewClient(
			client.ContentType("application/json"),
			client.Codec("application/json", jsonpbcodec.NewCodec()),
		),
		client.WithAddress("http://127.0.0.1:7989"),
	)
	m := meter.DefaultMeter
	tr := tracer.DefaultTracer
	l := logger.DefaultLogger

	f := flow.NewFlow(flow.Context(ctx), flow.Store(s), flow.Client(c), flow.Meter(m), flow.Tracer(tr), flow.Logger(l))

	if err := f.Init(); err != nil {
		t.Fatal(err)
	}

	options := append([]micro.Option{},
		micro.Server(
			httpsrv.NewServer(
				server.Codec("application/json", jsonpbcodec.NewCodec()),
				server.Address("127.0.0.1:7989"),
				httpsrv.RegisterRPCHandler(true),
			),
		),
		micro.Context(ctx),
	)

	svc := micro.NewService(options...)

	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}

	h := &handler{}

	if err := pb.RegisterTestServiceServer(svc.Server(), h); err != nil {
		t.Fatal(err)
	}

	go func() {
		if err := svc.Run(); err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(2 * time.Second)
	steps := []flow.Step{
		flow.NewCallStep("test", pb.TestServiceName, "LookupUser", flow.StepID("test.TestService.LookupUser")),
		flow.NewCallStep("test", pb.TestServiceName, "UpdateUser", flow.StepRequires("test.TestService.LookupUser")),
		flow.NewCallStep("test", pb.TestServiceName, "DeleteUser", flow.StepRequires("test.TestService.UpdateUser")),
		flow.NewCallStep("test", pb.TestServiceName, "MailUser", flow.StepRequires("test.TestService.UpdateUser")),
	}
	w, err := f.WorkflowCreate(ctx, "test", steps...)
	if err != nil {
		t.Fatal(err)
	}

	req, err := jsonpbcodec.NewCodec().Marshal(&pb.LookupUserReq{Name: "vtolstov"})
	if err != nil {
		t.Fatal(err)
	}

	id, err := w.Execute(ctx, &flow.Message{Body: req}, flow.ExecuteTimeout(2*time.Second))
	keys, err := s.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	keys, err = store.NewNamespaceStore(s, filepath.Join("workflows", id)).List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	keys, err = store.NewNamespaceStore(s, filepath.Join("steps", id)).List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_ = keys
	t.Logf("execution id: %s, result: %v", id, err)

}
