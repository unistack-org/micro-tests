package http_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	httpcli "github.com/unistack-org/micro-client-http/v3"
	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	jsonpbcodec "github.com/unistack-org/micro-codec-jsonpb/v3"
	vmeter "github.com/unistack-org/micro-meter-victoriametrics/v3"
	httpsrv "github.com/unistack-org/micro-server-http/v3"
	pb "github.com/unistack-org/micro-tests/server/http/proto"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/logger"
	lwrapper "github.com/unistack-org/micro/v3/logger/wrapper"
	"github.com/unistack-org/micro/v3/metadata"
	handler "github.com/unistack-org/micro/v3/meter/handler"
	mwrapper "github.com/unistack-org/micro/v3/meter/wrapper"
	"github.com/unistack-org/micro/v3/register"
	"github.com/unistack-org/micro/v3/server"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type Handler struct {
	t *testing.T
}

func NewServerHandlerWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			fmt.Printf("wrap ctx: %s\n", req.Service())
			return fn(ctx, req, rsp)
		}
	}
}

func (h *Handler) CallDouble(ctx context.Context, req *pb.CallReq, rsp *pb.CallRsp) error {
	rsp.Rsp = "name_double"
	httpsrv.SetRspCode(ctx, http.StatusCreated)
	return nil
}

func (h *Handler) Call(ctx context.Context, req *pb.CallReq, rsp *pb.CallRsp) error {
	if req.Nested == nil {
		h.t.Fatalf("invalid reflect merging")
	}
	if len(req.Nested.Uint64Args) != 3 || req.Nested.Uint64Args[2].Value != 3 {
		h.t.Fatalf("invalid reflect merging")
	}
	md, ok := metadata.FromIncomingContext(ctx)
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
	md = metadata.New(1)
	md.Set("my-key", "my-val")
	metadata.SetOutgoingContext(ctx, md)
	return nil
}

func (h *Handler) CallError(ctx context.Context, req *pb.CallReq1, rsp *pb.CallRsp1) error {
	httpsrv.SetRspCode(ctx, http.StatusBadRequest)
	return httpsrv.SetError(&pb.Error{Msg: "my_error"})
}

func TestNativeClientServer(t *testing.T) {
	reg := register.NewRegister()
	ctx := context.Background()

	var mwfOk bool
	mwf := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mwfOk = true
			next.ServeHTTP(w, r)
		})
	}

	m := vmeter.NewMeter()
	// create server
	srv := httpsrv.NewServer(
		server.Meter(m),
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("application/json", jsonpbcodec.NewCodec()),
		server.WrapHandler(mwrapper.NewHandlerWrapper(mwrapper.Meter(m))),
		server.WrapHandler(lwrapper.NewServerHandlerWrapper(lwrapper.WithEnabled(true), lwrapper.WithLevel(logger.InfoLevel))),
		httpsrv.Middleware(mwf),
		server.WrapHandler(NewServerHandlerWrapper()),
	)

	h := &Handler{t: t}

	// init server
	if err := srv.Init(); err != nil {
		t.Fatal(err)
	}

	if err := pb.RegisterTestServer(srv, h); err != nil {
		t.Fatal(err)
	}
	if err := pb.RegisterTestDoubleServer(srv, h); err != nil {
		t.Fatal(err)
	}
	if err := handler.RegisterMeterServer(srv, handler.NewHandler(handler.Meter(srv.Options().Meter))); err != nil {
		t.Fatal(err)
	}
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

	svc1 := pb.NewTestClient("helloworld", cli)
	rsp, err := svc1.Call(ctx, &pb.CallReq{
		Name: "my_name",
		Nested: &pb.Nested{Uint64Args: []*wrapperspb.UInt64Value{
			&wrapperspb.UInt64Value{Value: 1},
			&wrapperspb.UInt64Value{Value: 2},
			&wrapperspb.UInt64Value{Value: 3},
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

	hb, err := jsonpbcodec.NewCodec().Marshal(&pb.CallReq{
		Nested: &pb.Nested{Uint64Args: []*wrapperspb.UInt64Value{
			&wrapperspb.UInt64Value{Value: 1},
			&wrapperspb.UInt64Value{Value: 2},
			&wrapperspb.UInt64Value{Value: 3},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("test rsp code from net/http client to native micro http server")
	hr, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://%s/v1/test/call/my_name", service[0].Nodes[0].Address), bytes.NewReader(hb))
	if err != nil {
		t.Fatal(err)
	}
	hr.Header.Set("Content-Type", "application/json")

	hrsp, err := http.DefaultClient.Do(hr)
	if err != nil {
		t.Fatal(err)
	}
	defer hrsp.Body.Close()
	buf, err := io.ReadAll(hrsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if hrsp.StatusCode != 201 {
		t.Fatalf("invalid rsp code %#+v", hrsp)
	}

	t.Logf("test second server")
	svc2 := pb.NewTestDoubleClient("helloworld", cli)
	rsp, err = svc2.CallDouble(ctx, &pb.CallReq{
		Name: "my_name",
	})
	if err != nil {
		t.Fatal(err)
	}

	if rsp.Rsp != "name_double" {
		t.Fatalf("invalid response: %#+v\n", rsp)
	}

	hrsp, err = http.Get(fmt.Sprintf("http://%s/metrics", service[0].Nodes[0].Address))
	if err != nil {
		t.Fatal(err)
	}

	buf, err = io.ReadAll(hrsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(buf), `micro_server_request_total{micro_status="success"}`) {
		t.Fatalf("rsp not contains metrics: %s", buf)
	}
	// stop server
	if err := srv.Stop(); err != nil {
		t.Fatal(err)
	}

}

func TestNativeServer(t *testing.T) {
	reg := register.NewRegister()
	ctx := context.Background()

	// create server
	srv := httpsrv.NewServer(
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("application/json", jsoncodec.NewCodec()),
		//server.WrapHandler(NewServerHandlerWrapper()),
	)

	h := &Handler{t: t}
	pb.RegisterTestServer(srv, h)

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

	if v := rsp.Header.Get("My-Key"); v != "my-val" {
		t.Fatalf("empty response header: %#+v", rsp.Header)
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

func TestHTTPHandler(t *testing.T) {
	reg := register.NewRegister()
	ctx := context.Background()

	// create server
	srv := httpsrv.NewServer(
		server.Register(reg),
		server.Codec("application/json", jsoncodec.NewCodec()),
	)

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

func TestHTTPServer(t *testing.T) {
	reg := register.NewRegister()
	ctx := context.Background()

	// create server mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello world`))
	})

	// create server
	srv := httpsrv.NewServer(
		server.Register(reg),
		httpsrv.Server(&http.Server{Handler: mux}),
		server.Codec("application/json", jsoncodec.NewCodec()),
	)

	if err := srv.Init(); err != nil {
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
