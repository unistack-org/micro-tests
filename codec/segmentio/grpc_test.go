package grpc_test

import (
	"context"
	"testing"

	gclient "go.unistack.org/micro-client-grpc/v3"
	//	protocodec "go.unistack.org/micro-codec-proto/v3"
	protocodec "go.unistack.org/micro-codec-segmentio/v3/proto"
	regRouter "go.unistack.org/micro-router-register/v3"
	gserver "go.unistack.org/micro-server-grpc/v3"
	gpb "go.unistack.org/micro-tests/server/grpc/gproto"
	pb "go.unistack.org/micro-tests/server/grpc/proto"
	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/errors"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/router"
	"go.unistack.org/micro/v3/server"
)

type testServer struct {
	pb.UnimplementedTestServer
}

func (g *testServer) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	if req.Name == "Error" {
		return &errors.Error{ID: "id", Code: 99, Detail: "detail"}
	}
	rsp.Msg = "Hello " + req.Name
	return nil
}

func TestGRPCServer(t *testing.T) {
	var err error

	r := register.NewRegister()
	b := broker.NewBroker(broker.Register(r))
	s := gserver.NewServer(
		server.Codec("application/grpc+proto", protocodec.NewCodec()),
		server.Codec("application/grpc", protocodec.NewCodec()),
		server.Address("127.0.0.1:0"),
		server.Register(r),
		server.Name("helloworld"),
		gserver.Reflection(true),
	)
	// create router
	rtr := regRouter.NewRouter(router.Register(r))

	h := &testServer{}
	err = gpb.RegisterTestServer(s, h)
	if err != nil {
		t.Fatalf("can't register handler: %v", err)
	}

	if err = s.Init(); err != nil {
		t.Fatal(err)
	}

	if err = s.Start(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = s.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	// create client
	c := gclient.NewClient(
		client.Codec("application/grpc+proto", protocodec.NewCodec()),
		client.Codec("application/grpc", protocodec.NewCodec()),
		client.Router(rtr),
		client.Register(r),
		client.Broker(b),
	)

	testMethods := []string{
		"Test.Call",
	}

	for _, method := range testMethods {
		req := c.NewRequest("helloworld", method, &pb.Request{
			Name: "John",
		})

		rsp := pb.Response{}

		err = c.Call(context.TODO(), req, &rsp)
		if err != nil {
			t.Fatalf("method: %s err: %v", method, err)
		}

		if rsp.Msg != "Hello John" {
			t.Fatalf("Got unexpected response %v", rsp.Msg)
		}
	}

	//rsp := rpb.ServerReflectionResponse{}
	//req := c.NewRequest("helloworld", "Test.ServerReflectionInfo", &rpb.ServerReflectionRequest{}, client.StreamingRequest())
	//if err := c.Call(context.TODO(), req, &rsp); err != nil {
	//	t.Fatal(err)
	//}

	//	select {}
}
