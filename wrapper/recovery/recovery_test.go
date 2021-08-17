package recovery_test

import (
	"context"
	"fmt"
	"testing"

	cli "github.com/unistack-org/micro-client-grpc/v3"
	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	rrouter "github.com/unistack-org/micro-router-register/v3"
	srv "github.com/unistack-org/micro-server-grpc/v3"
	recwrapper "github.com/unistack-org/micro-wrapper-recovery/v3"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/errors"
	"github.com/unistack-org/micro/v3/register"
	"github.com/unistack-org/micro/v3/router"
	"github.com/unistack-org/micro/v3/server"
)

type Test interface {
	Method(ctx context.Context, in *TestRequest, opts ...client.CallOption) (*TestResponse, error)
}

type TestRequest struct {
	IsPanic bool
}
type TestResponse struct{}

type testHandler struct{}

func (t *testHandler) Method(ctx context.Context, req *TestRequest, rsp *TestResponse) error {
	if req.IsPanic {
		panic("panic here")
	}
	return nil
}

func TestRecovery(t *testing.T) {
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

	rfn := func(ctx context.Context, req server.Request, rsp interface{}, err error) error {
		return errors.BadRequest("id-1234567890", "handled panic: %v", err)
	}

	s := srv.NewServer(
		server.Codec("application/grpc+json", jsoncodec.NewCodec()),
		server.Codec("application/json", jsoncodec.NewCodec()),
		server.Name(name),
		server.Version(version),
		server.ID(id),
		server.Register(reg),
		server.Broker(brk),
		server.WrapHandler(
			recwrapper.NewHandlerWrapper(rfn),
		),
	)

	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	type Test struct {
		*testHandler
	}

	if err := s.Handle(s.NewHandler(&Test{new(testHandler)})); err != nil {
		t.Fatal(err)
	}

	if err := s.Start(); err != nil {
		t.Fatalf("Unexpected error starting server: %v", err)
	}
	defer func() {
		_ = s.Stop()
	}()

	req := c.NewRequest(name, "Test.Method", &TestRequest{IsPanic: true}, client.RequestContentType("application/json"))
	rsp := TestResponse{}

	err := c.Call(context.TODO(), req, &rsp)
	if err == nil {
		t.Fatalf("panic happens, but handler not return err")
	}

	fmt.Printf("%v\n", err)
}
