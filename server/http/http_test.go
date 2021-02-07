package http_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	wrapperpb "github.com/golang/protobuf/ptypes/wrappers"
	httpcli "github.com/unistack-org/micro-client-http/v3"
	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	jsonpbcodec "github.com/unistack-org/micro-codec-jsonpb/v3"
	memory "github.com/unistack-org/micro-register-memory/v3"
	httpsrv "github.com/unistack-org/micro-server-http/v3"
	pb "github.com/unistack-org/micro-tests/server/http/proto"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/metadata"
	"github.com/unistack-org/micro/v3/server"
)

type Handler struct {
	t *testing.T
}

func NewServerHandlerWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			fmt.Printf("wrap ctx: %#+v req: %#+v\n", ctx, req)
			return fn(ctx, req, rsp)
		}
	}
}

func (h *Handler) Call(ctx context.Context, req *pb.CallReq, rsp *pb.CallRsp) error {
	if req.Nested == nil {
		h.t.Fatalf("invalid reflect merging")
	}
	if len(req.Nested.Uint64Args) != 3 || req.Nested.Uint64Args[2].Value != 3 {
		h.t.Fatalf("invalid reflect merging")
	}
	md, ok := metadata.FromContext(ctx)
	if !ok {
		h.t.Fatalf("context without metadata")
	}
	if _, ok := md.Get("User-Agent"); !ok {
		h.t.Fatalf("context metadata does not have User-Agent header")
	}
	if req.Name != "my_name" {
		h.t.Fatalf("invalid req received: %#+v", req)
	}
	rsp.Rsp = "name_my_name"
	httpsrv.SetRspCode(ctx, http.StatusCreated)
	return nil
}

func (h *Handler) CallError(ctx context.Context, req *pb.CallReq1, rsp *pb.CallRsp1) error {
	httpsrv.SetRspCode(ctx, http.StatusBadRequest)
	return &pb.Error{Msg: "my_error"}
	return nil
}

func TestNativeClientServer(t *testing.T) {
	reg := memory.NewRegister()
	ctx := context.Background()

	var mwfOk bool
	mwf := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mwfOk = true
			next.ServeHTTP(w, r)
		})
	}

	// create server
	srv := httpsrv.NewServer(
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("application/json", jsonpbcodec.NewCodec()),
		httpsrv.Middleware(mwf),
		//server.WrapHandler(NewServerHandlerWrapper()),
	)

	h := &Handler{t: t}
	pb.RegisterTestHandler(srv, h)

	// start server
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}

	// lookup server
	service, err := reg.LookupService(ctx, "helloworld")
	if err != nil {
		t.Fatal(err)
	}

	if len(service) != 1 {
		t.Fatalf("Expected 1 service got %d: %+v", len(service), service)
	}

	if len(service[0].Nodes) != 1 {
		t.Fatalf("Expected 1 node got %d: %+v", len(service[0].Nodes), service[0].Nodes)
	}

	cli := client.NewClientCallOptions(httpcli.NewClient(client.ContentType("application/json"), client.Codec("application/json", jsonpbcodec.NewCodec())), client.WithAddress(fmt.Sprintf("http://%s", service[0].Nodes[0].Address)))

	svc := pb.NewTestService("helloworld", cli)
	rsp, err := svc.Call(ctx, &pb.CallReq{
		Name: "my_name",
		Nested: &pb.Nested{Uint64Args: []*wrapperpb.UInt64Value{
			&wrapperpb.UInt64Value{Value: 1},
			&wrapperpb.UInt64Value{Value: 2},
			&wrapperpb.UInt64Value{Value: 3},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	if rsp.Rsp != "name_my_name" {
		t.Fatalf("invalid response: %#+v\n", rsp)
	}

	if !mwfOk {
		t.Fatalf("http middleware not works")
	}
	// stop server
	if err := srv.Stop(); err != nil {
		t.Fatal(err)
	}

}

func TestNativeServer(t *testing.T) {
	reg := memory.NewRegister()
	ctx := context.Background()

	// create server
	srv := httpsrv.NewServer(
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("application/json", jsoncodec.NewCodec()),
		//server.WrapHandler(NewServerHandlerWrapper()),
	)

	h := &Handler{t: t}
	pb.RegisterTestHandler(srv, h)

	// start server
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}

	// lookup server
	service, err := reg.LookupService(ctx, "helloworld")
	if err != nil {
		t.Fatal(err)
	}

	if len(service) != 1 {
		t.Fatalf("Expected 1 service got %d: %+v", len(service), service)
	}

	if len(service[0].Nodes) != 1 {
		t.Fatalf("Expected 1 node got %d: %+v", len(service[0].Nodes), service[0].Nodes)
	}

	// make request
	rsp, err := http.Post(fmt.Sprintf("http://%s/v1/test/call/my_name?req=key&arg1=arg1&arg2=12345&nested.string_args=str1,str2&nested.uint64_args=1,2,3", service[0].Nodes[0].Address), "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	if rsp.StatusCode != http.StatusCreated {
		t.Fatalf("invalid status received: %#+v\n", rsp)
	}

	b, err := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	if s := string(b); s != `{"rsp":"name_my_name"}` {
		t.Fatalf("Expected response %s, got %s", `{"rsp":"name_my_name"}`, s)
	}

	// make request with error
	rsp, err = http.Post(fmt.Sprintf("http://%s/v1/test/callerror/my_name", service[0].Nodes[0].Address), "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	if rsp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid status received: %#+v\n", rsp)
	}

	b, err = ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	if s := string(b); s != `{"msg":"my_error"}` {
		t.Fatalf("Expected response %s, got %s", `{"msg":"my_error"}`, s)
	}

	// stop server
	if err := srv.Stop(); err != nil {
		t.Fatal(err)
	}

}

func TestHTTPServer(t *testing.T) {
	reg := memory.NewRegister()
	ctx := context.Background()

	// create server
	srv := httpsrv.NewServer(server.Register(reg))

	// create server mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello world`))
	})

	// create handler
	hd := srv.NewHandler(mux)

	// register handler
	if err := srv.Handle(hd); err != nil {
		t.Fatal(err)
	}

	// start server
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}

	// lookup server
	service, err := reg.LookupService(ctx, server.DefaultName)
	if err != nil {
		t.Fatal(err)
	}

	if len(service) != 1 {
		t.Fatalf("Expected 1 service got %d: %+v", len(service), service)
	}

	if len(service[0].Nodes) != 1 {
		t.Fatalf("Expected 1 node got %d: %+v", len(service[0].Nodes), service[0].Nodes)
	}

	// make request
	rsp, err := http.Get(fmt.Sprintf("http://%s", service[0].Nodes[0].Address))
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if s := string(b); s != "hello world" {
		t.Fatalf("Expected response %s, got %s", "hello world", s)
	}

	// stop server
	if err := srv.Stop(); err != nil {
		t.Fatal(err)
	}
}
