package grpc_test

import (
	"context"
	"io"
	"net"
	"net/http"
	"testing"

	gclient "go.unistack.org/micro-client-grpc/v3"
	protocodec "go.unistack.org/micro-codec-proto/v3"
	regRouter "go.unistack.org/micro-router-register/v3"
	gserver "go.unistack.org/micro-server-grpc/v3"
	httpsrv "go.unistack.org/micro-server-http/v3"
	gpb "go.unistack.org/micro-tests/server/grpc/gproto"
	pb "go.unistack.org/micro-tests/server/grpc/proto"
	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/errors"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/router"
	"go.unistack.org/micro/v3/server"
	health "go.unistack.org/micro/v3/server/health"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	gmetadata "google.golang.org/grpc/metadata"
)

type testServer struct {
	pb.UnimplementedTestServer
}

type testnServer struct {
	pb.UnimplementedTestServer
}

func NewServerHandlerWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// fmt.Printf("wrap ctx: %#+v req: %#+v\n", ctx, req)
			return fn(ctx, req, rsp)
		}
	}
}

func (g *testServer) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.InternalServerError("xxx", "missing metadata")
	}
	if req.Name == "Error" {
		return &errors.Error{ID: "id", Code: 99, Detail: "detail"}
	}
	rsp.Msg = "Hello " + req.Name
	rsp.Broken = &pb.Broken{Field: "12345"}

	return nil
}

func (g *testnServer) Call(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	_, ok := gmetadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.InternalServerError("xxx", "missing metadata")
	}

	if req.Name == "Error" {
		return nil, &errors.Error{ID: "id", Code: 99, Detail: "detail"}
	}
	rsp := &pb.Response{}

	for i := 0; i < 650; i++ {
		rsp.Msg += "Hello " + req.Name
	}
	return rsp, nil
}

func TestGRPCServer(t *testing.T) {
	var err error
	codec.DefaultMaxMsgSize = 8 * 1024 * 1024
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = logger.DefaultLogger.Init(logger.WithLevel(logger.ErrorLevel))
	r := register.NewRegister()
	b := broker.NewBroker(broker.Register(r))
	s := gserver.NewServer(
		server.Codec("application/grpc+proto", protocodec.NewCodec()),
		server.Codec("application/grpc", protocodec.NewCodec()),
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

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	pb.RegisterTestServer(gs, &testnServer{})
	go func() {
		if err := gs.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	}()

	// lookup server
	service, err := r.LookupService(ctx, "helloworld")
	if err != nil {
		t.Fatal(err)
	}

	if len(service) != 1 {
		t.Fatalf("Expected 1 service got %d: %+v", len(service), service)
	}

	if len(service[0].Nodes) != 1 {
		t.Fatalf("Expected 1 node got %d: %+v", len(service[0].Nodes), service[0].Nodes)
	}

	// create client
	gc := gclient.NewClient(
		client.ContentType("application/grpc"),
		client.Codec("application/grpc", protocodec.NewCodec()), client.Router(rtr), client.Register(r), client.Broker(b))

	c := gpb.NewTestClient("helloworld", gc)

	var md metadata.Metadata
	t.Logf("call micro via micro")
	rq := &pb.Request{Name: "John"}
	for i := 0; i < 1500; i++ {
		rq.Name += "name"
	}
	_, err = c.Call(ctx, rq,
		client.WithResponseMetadata(&md),
		gclient.CallOptions(grpc.UseCompressor(gzip.Name)))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Logf("response md %#+v", md)

	ngcli, err := grpc.DialContext(ctx,
		// lis.Addr().String(),
		service[0].Nodes[0].Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer ngcli.Close()

	var gmd gmetadata.MD
	ngrpcsvc := pb.NewTestClient(ngcli)
	t.Logf("call micro via native")
	if _, err = ngrpcsvc.Call(ctx, rq,
		grpc.UseCompressor(gzip.Name),
		grpc.Header(&gmd)); err != nil {
		t.Fatal(err)
	}
	t.Logf("gmd %#+v\n", gmd)

	nxgcli, err := grpc.DialContext(ctx,
		lis.Addr().String(),
		// service[0].Nodes[0].Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer nxgcli.Close()

	ngrpcsvc = pb.NewTestClient(nxgcli)
	t.Logf("call native via native")
	if _, err := ngrpcsvc.Call(ctx, rq,
		grpc.UseCompressor(gzip.Name),
		grpc.Header(&gmd)); err != nil {
		t.Fatal(err)
	}
	t.Logf("gmd %#+v\n", gmd)

	//rsp := rpb.ServerReflectionResponse{}
	//req := c.NewRequest("helloworld", "Test.ServerReflectionInfo", &rpb.ServerReflectionRequest{}, client.StreamingRequest())
	//if err := c.Call(context.TODO(), req, &rsp); err != nil {
	//	t.Fatal(err)
	//}

	//	select {}
}
