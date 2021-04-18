package http_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"

	httpcli "github.com/unistack-org/micro-client-http/v3"
	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	jsonpbcodec "github.com/unistack-org/micro-codec-jsonpb/v3"
	urlencodecodec "github.com/unistack-org/micro-codec-urlencode/v3"
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

func multipartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%#+v\n", r)
}

func upload(client *http.Client, url string, values map[string]io.Reader) error {
	var err error
	b := bytes.NewBuffer(nil)
	w := multipart.NewWriter(b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if fw, err = w.CreateFormFile(key, key); err != nil {
			return err
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check the response
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		err = fmt.Errorf("bad status: %s", res.Status)
	}

	return err
}

func TestMultipart(t *testing.T) {
	reg := register.NewRegister()
	ctx := context.Background()

	// create server
	srv := httpsrv.NewServer(
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("application/json", jsoncodec.NewCodec()),
		server.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
		httpsrv.PathHandler("/upload", multipartHandler),
	)

	if err := srv.Init(); err != nil {
		t.Fatal(err)
	}

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

	t.Logf("test multipart upload")
	values := make(map[string]io.Reader, 2)
	values["first.txt"] = bytes.NewReader([]byte("first content"))
	values["second.txt"] = bytes.NewReader([]byte("second content"))
	err = upload(http.DefaultClient, "http://"+service[0].Nodes[0].Address+"/upload", values)
	if err != nil {
		t.Fatal(err)
	}
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

func (h *Handler) CallRepeated(ctx context.Context, req *pb.CallReq, rsp *pb.CallRsp) error {
	if len(req.Ids) != 2 || req.Ids[0] != "123" {
		h.t.Fatalf("invalid reflect merging")
	}
	rsp.Rsp = "name_my_name"
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

func TestNativeFormUrlencoded(t *testing.T) {
	reg := register.NewRegister()
	ctx := context.Background()

	// create server
	srv := httpsrv.NewServer(
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("application/json", jsoncodec.NewCodec()),
		server.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
		//server.WrapHandler(NewServerHandlerWrapper()),
	)

	if err := srv.Init(); err != nil {
		t.Fatal(err)
	}
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

	t.Logf("test net/http client with application/x-www-form-urlencoded")
	data := url.Values{}
	data.Set("req", "fookey")
	data.Set("arg1", "arg1val")
	data.Add("nested.uint64_args", "1")
	data.Add("nested.uint64_args", "2")
	data.Add("nested.uint64_args", "3")
	// make request
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/v1/test/call/my_name", service[0].Nodes[0].Address), strings.NewReader(data.Encode())) // URL-encoded payload
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()

	if err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if rsp.StatusCode != http.StatusCreated {
		t.Fatalf("invalid status received: %#+v\n%s\n", rsp, b)
	}

	if s := string(b); s != `{"rsp":"name_my_name"}` {
		t.Fatalf("Expected response %s, got %s", `{"rsp":"name_my_name"}`, s)
	}

	if v := rsp.Header.Get("My-Key"); v != "my-val" {
		t.Fatalf("empty response header: %#+v", rsp.Header)
	}

	t.Logf("test native client with application/x-www-form-urlencoded")
	cli := client.NewClientCallOptions(
		httpcli.NewClient(
			client.ContentType("application/x-www-form-urlencoded"),
			client.Codec("application/json", jsonpbcodec.NewCodec()),
			client.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
		),
		client.WithAddress(fmt.Sprintf("http://%s", service[0].Nodes[0].Address)))

	svc1 := pb.NewTestClient("helloworld", cli)
	nrsp, err := svc1.Call(ctx, &pb.CallReq{
		Name: "my_name",
		Arg1: "arg1val",
		Nested: &pb.Nested{Uint64Args: []*wrapperspb.UInt64Value{
			&wrapperspb.UInt64Value{Value: 1},
			&wrapperspb.UInt64Value{Value: 2},
			&wrapperspb.UInt64Value{Value: 3},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	if nrsp.Rsp != "name_my_name" {
		t.Fatalf("invalid response: %#+v\n", nrsp)
	}

	// stop server
	if err := srv.Stop(); err != nil {
		t.Fatal(err)
	}

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
		server.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
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
		server.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
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
	rsp, err := http.Post(fmt.Sprintf("http://%s/v1/test/call/my_name?req=key&arg1=arg1&arg2=12345&nested.string_args=str1&nested.string_args=str2&nested.uint64_args=1&nested.uint64_args=2&nested.uint64_args=3", service[0].Nodes[0].Address), "application/json", nil)
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

	rsp, err = http.Post(fmt.Sprintf("http://%s/v1/test/call_repeated/?ids=123&ids=321", service[0].Nodes[0].Address), "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	if rsp.StatusCode != http.StatusCreated {
		buf, err := io.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("invalid status received: %#+v err: %v\n", rsp, err)
		}
		t.Fatalf("invalid status received: %#+v buf: %s\n", rsp, buf)
	}

	b, err = ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	if s := string(b); s != `{"rsp":"name_my_name"}` {
		t.Fatalf("Expected response %s, got %s", `{"rsp":"name_my_name"}`, s)
	}

	c := client.NewClientCallOptions(httpcli.NewClient(client.ContentType("application/json"), client.Codec("application/json", jsoncodec.NewCodec())), client.WithAddress("http://"+service[0].Nodes[0].Address))
	pbc := pb.NewTestClient("test", c)

	prsp, err := pbc.CallRepeated(context.TODO(), &pb.CallReq{Ids: []string{"123", "321"}})
	if err != nil {
		t.Fatal(err)
	}

	if prsp.Rsp != "name_my_name" {
		t.Fatalf("invalid rsp received: %#+v\n", rsp)
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
		server.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
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
		server.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
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
