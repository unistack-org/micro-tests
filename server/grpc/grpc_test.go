package grpc_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	gclient "github.com/unistack-org/micro-client-grpc/v3"
	protocodec "github.com/unistack-org/micro-codec-proto/v3"
	regRouter "github.com/unistack-org/micro-router-register/v3"
	gserver "github.com/unistack-org/micro-server-grpc/v3"
	httpsrv "github.com/unistack-org/micro-server-http/v3"
	gpb "github.com/unistack-org/micro-tests/server/grpc/gproto"
	pb "github.com/unistack-org/micro-tests/server/grpc/proto"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/errors"
	"github.com/unistack-org/micro/v3/register"
	"github.com/unistack-org/micro/v3/router"
	"github.com/unistack-org/micro/v3/server"
	health "github.com/unistack-org/micro/v3/server/health"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
)

type testServer struct {
	pb.UnimplementedTestServer
}

func NewServerHandlerWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			//fmt.Printf("wrap ctx: %#+v req: %#+v\n", ctx, req)
			return fn(ctx, req, rsp)
		}
	}
}

func (g *testServer) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	if req.Name == "Error" {
		return &errors.Error{Id: "id", Code: 99, Detail: "detail"}
	}
	rsp.Msg = "Hello " + req.Name
	rsp.Broken = &pb.Broken{Field: "12345"}

	return nil
}

func TestGRPCServer(t *testing.T) {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := register.NewRegister()
	b := broker.NewBroker(broker.Register(r))
	s := gserver.NewServer(
		server.Codec("application/grpc+proto", protocodec.NewCodec()),
		server.Address("127.0.0.1:0"),
		server.Register(r),
		server.Name("helloworld"),
		gserver.Reflection(true),
		server.WrapHandler(NewServerHandlerWrapper()),
	)
	// create router
	rtr := regRouter.NewRouter(router.Register(r))

	h := &testServer{}
	if err = gpb.RegisterTestServer(s, h); err != nil {
		t.Fatalf("can't register handler: %v", err)
	}

	srv := httpsrv.NewServer(
		server.Address("127.0.0.1:0"),
		server.Codec("text/plain", codec.NewCodec()),
	)
	if err = health.RegisterHealthServer(srv, health.NewHandler(health.Version("0.0.1"))); err != nil {
		t.Fatalf("cant register health handler: %v", err)
	}

	if err = srv.Init(); err != nil {
		t.Fatal(err)
	}

	if err = srv.Start(); err != nil {
		t.Fatal(err)
	}

	if err = s.Init(); err != nil {
		t.Fatal(err)
	}

	if err = s.Start(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = srv.Stop(); err != nil {
			t.Fatal(err)
		}

		if err = s.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	hr, err := http.NewRequestWithContext(ctx, "GET", "http://"+srv.Options().Address+"/version", nil)
	if err != nil {
		t.Fatal(err)
	}
	hr.Header.Set("Content-Type", "text/plain")
	rsp, err := http.DefaultClient.Do(hr)
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()
	buf, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	} else if string(buf) != "0.0.1" {
		t.Fatalf("unknown version returned from health handler: %s", buf)
	}

	// create client
	c := gclient.NewClient(client.Codec("application/grpc+proto", protocodec.NewCodec()), client.Router(rtr), client.Register(r), client.Broker(b))

	testMethods := []string{
		"Test.Call",
	}

	for _, method := range testMethods {
		req := c.NewRequest("helloworld", method, &pb.Request{
			Name: "John",
		})

		rsp := &pb.Response{}

		err = c.Call(context.TODO(), req, rsp)
		if err != nil {
			t.Fatalf("method: %s err: %v", method, err)
		}

		if rsp.Msg != "Hello John" {
			t.Fatalf("Got unexpected response %v", rsp.Msg)
		}

		enc := &jsonpb.MarshalOptions{EmitUnpopulated: true}
		_, err := enc.Marshal(rsp)
		if err != nil {
			t.Fatal(err)
		}
	}

	//rsp := rpb.ServerReflectionResponse{}
	//req := c.NewRequest("helloworld", "Test.ServerReflectionInfo", &rpb.ServerReflectionRequest{}, client.StreamingRequest())
	//if err := c.Call(context.TODO(), req, &rsp); err != nil {
	//	t.Fatal(err)
	//}

	//	select {}
}
