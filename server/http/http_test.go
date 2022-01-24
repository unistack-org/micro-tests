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
	"sync"
	"testing"

	httpcli "go.unistack.org/micro-client-http/v3"
	jsoncodec "go.unistack.org/micro-codec-json/v3"
	jsonpbcodec "go.unistack.org/micro-codec-jsonpb/v3"
	urlencodecodec "go.unistack.org/micro-codec-urlencode/v3"
	xmlcodec "go.unistack.org/micro-codec-xml/v3"
	vmeter "go.unistack.org/micro-meter-victoriametrics/v3"
	httpsrv "go.unistack.org/micro-server-http/v3"
	pb "go.unistack.org/micro-tests/server/http/proto"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/logger"
	lwrapper "go.unistack.org/micro/v3/logger/wrapper"
	"go.unistack.org/micro/v3/metadata"
	handler "go.unistack.org/micro/v3/meter/handler"
	mwrapper "go.unistack.org/micro/v3/meter/wrapper"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/server"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type Handler struct {
	t *testing.T
}

func multipartHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("%#+v\n", r)
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
		server.Address("127.0.0.1:0"),
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("application/json", jsoncodec.NewCodec()),
		server.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
		httpsrv.PathHandler(http.MethodPost, "/upload", multipartHandler),
	)

	if err := srv.Init(); err != nil {
		t.Fatal(err)
	}

	h := &Handler{t: t}
	if err := pb.RegisterTestServer(srv, h); err != nil {
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

	// t.Logf("test multipart upload")
	values := make(map[string]io.Reader, 2)
	values["first.txt"] = bytes.NewReader([]byte("first content"))
	values["second.txt"] = bytes.NewReader([]byte("second content"))
	err = upload(http.DefaultClient, "http://"+service[0].Nodes[0].Address+"/upload", values)
	if err != nil {
		t.Fatal(err)
	}
}

func NewServerHandlerWrapper(t *testing.T) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				t.Fatal("metadata empty")
			}
			if v, ok := md.Get("Authorization"); ok && v == "test" {
				nmd := metadata.New(1)
				nmd.Set("my-key", "my-val")
				nmd.Set("Content-Type", "text/xml")
				metadata.SetOutgoingContext(ctx, nmd)
				httpsrv.SetRspCode(ctx, http.StatusUnauthorized)
				return httpsrv.SetError(&pb.CallRsp{Rsp: "name_my_name"})
			}

			if v, ok := md.Get("Test-Content-Type"); ok && v != "" {
				nmd := metadata.New(1)
				nmd.Set("my-key", "my-val")
				nmd.Set("Content-Type", v)
				metadata.SetOutgoingContext(ctx, nmd)
			}

			return fn(ctx, req, rsp)
		}
	}
}

func (h *Handler) CallDouble(ctx context.Context, req *pb.CallReq, rsp *pb.CallRsp) error {
	rsp.Rsp = "name_double"
	httpsrv.SetRspCode(ctx, http.StatusCreated)
	return nil
}

func (h *Handler) CallRepeatedString(ctx context.Context, req *pb.CallReq, rsp *pb.CallRsp) error {
	if len(req.StringIds) != 2 || req.StringIds[0] != "123" {
		h.t.Fatalf("invalid reflect merging, strings_ids invalid: %v", req.StringIds)
	}
	rsp.Rsp = "name_my_name"
	httpsrv.SetRspCode(ctx, http.StatusCreated)
	return nil
}

func (h *Handler) CallRepeatedInt64(ctx context.Context, req *pb.CallReq, rsp *pb.CallRsp) error {
	if len(req.Int64Ids) != 2 || req.Int64Ids[0] != 123 {
		h.t.Fatalf("invalid reflect merging, int64_ids invalid: %v", req.Int64Ids)
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
	if req.Clientid != "1234567890" {
		h.t.Fatalf("invalid req recevided %#+v", req)
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
		server.Address("127.0.0.1:0"),
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("application/json", jsoncodec.NewCodec()),
		server.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
	// server.WrapHandler(NewServerHandlerWrapper()),
	)

	if err := srv.Init(); err != nil {
		t.Fatal(err)
	}
	h := &Handler{t: t}
	if err := pb.RegisterTestServer(srv, h); err != nil {
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

	data := url.Values{}
	data.Set("req", "fookey")
	data.Set("arg1", "arg1val")
	data.Add("nested.uint64_args", "1")
	data.Add("nested.uint64_args", "2")
	data.Add("nested.uint64_args", "3")
	// make request
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/v1/test/call/my_name", service[0].Nodes[0].Address), strings.NewReader(data.Encode())) // URL-encoded payload
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Clientid", "1234567890")
	req.AddCookie(&http.Cookie{Name: "Csrftoken", Value: "csrftoken"})
	// req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if err != nil {
		t.Fatalf("test net/http client with application/x-www-form-urlencoded err: %v", err)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("test net/http client with application/x-www-form-urlencoded err: %v", err)
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

	cli := client.NewClientCallOptions(
		httpcli.NewClient(
			client.ContentType("application/x-www-form-urlencoded"),
			client.Codec("application/json", jsonpbcodec.NewCodec()),
			client.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
		),
		client.WithAddress(fmt.Sprintf("http://%s", service[0].Nodes[0].Address)))

	svc1 := pb.NewTestClient("helloworld", cli)
	nrsp, err := svc1.Call(ctx, &pb.CallReq{
		Name:      "my_name",
		Arg1:      "arg1val",
		Clientid:  "1234567890",
		Csrftoken: "csrftoken",
		Nested: &pb.Nested{Uint64Args: []*wrapperspb.UInt64Value{
			&wrapperspb.UInt64Value{Value: 1},
			&wrapperspb.UInt64Value{Value: 2},
			&wrapperspb.UInt64Value{Value: 3},
		}},
	})
	if err != nil {
		t.Fatalf("test native client with application/x-www-form-urlencoded err: %v", err)
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
		server.Address("127.0.0.1:0"),
		server.Meter(m),
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("application/json", jsonpbcodec.NewCodec()),
		server.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
		server.WrapHandler(mwrapper.NewHandlerWrapper(mwrapper.Meter(m))),
		server.WrapHandler(lwrapper.NewServerHandlerWrapper(lwrapper.WithEnabled(false), lwrapper.WithLevel(logger.ErrorLevel))),
		httpsrv.Middleware(mwf),
		server.WrapHandler(NewServerHandlerWrapper(t)),
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
		Name:      "my_name",
		Clientid:  "1234567890",
		Csrftoken: "csrftoken",
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
		Clientid:  "1234567890",
		Csrftoken: "csrftoken",
		Nested: &pb.Nested{Uint64Args: []*wrapperspb.UInt64Value{
			&wrapperspb.UInt64Value{Value: 1},
			&wrapperspb.UInt64Value{Value: 2},
			&wrapperspb.UInt64Value{Value: 3},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	hr, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://%s/v1/test/call/my_name", service[0].Nodes[0].Address), bytes.NewReader(hb))
	if err != nil {
		t.Fatalf("test rsp code from net/http client to native micro http server err: %v", err)
	}
	hr.Header.Set("Content-Type", "application/json")

	hrsp, err := http.DefaultClient.Do(hr)
	if err != nil {
		t.Fatalf("test rsp code from net/http client to native micro http server err: %v", err)
	}
	defer func() {
		_ = hrsp.Body.Close()
	}()

	_, err = io.ReadAll(hrsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if hrsp.StatusCode != 201 {
		t.Fatalf("invalid rsp code %#+v", hrsp)
	}

	svc2 := pb.NewTestDoubleClient("helloworld", cli)
	rsp, err = svc2.CallDouble(ctx, &pb.CallReq{
		Name: "my_name",
	})
	if err != nil {
		t.Fatalf("test second server err: %v", err)
	}

	if rsp.Rsp != "name_double" {
		t.Fatalf("test second server invalid response: %#+v\n", rsp)
	}

	hrsp, err = http.Get(fmt.Sprintf("http://%s/metrics", service[0].Nodes[0].Address))
	if err != nil {
		t.Fatal(err)
	}

	buf, err := io.ReadAll(hrsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(buf), `micro_server_request_total`) {
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
		server.Address("127.0.0.1:0"),
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("text/xml", xmlcodec.NewCodec()),
		server.Codec("application/json", jsoncodec.NewCodec()),
		server.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
		server.WrapHandler(NewServerHandlerWrapper(t)),
	)

	h := &Handler{t: t}
	if err := pb.RegisterTestServer(srv, h); err != nil {
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

	// make request
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/v1/test/call/my_name?req=key&arg1=arg1&arg2=12345&nested.string_args=str1&nested.string_args=str2&nested.uint64_args=1&nested.uint64_args=2&nested.uint64_args=3", service[0].Nodes[0].Address), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "test")
	req.Header.Set("Content-Type", "application/json")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if rsp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("invalid status received: %#+v\n", rsp)
	}

	b, err := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	if s := string(b); s != `<CallRsp><Rsp>name_my_name</Rsp></CallRsp>` {
		t.Fatalf("Expected response %s, got %s", `<CallRsp><Rsp>name_my_name</Rsp></CallRsp>`, s)
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

	rsp, err = http.Post(fmt.Sprintf("http://%s/v1/test/call_repeated_string?string_ids=123&string_ids=321", service[0].Nodes[0].Address), "application/json", nil)
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

	prsp, err := pbc.CallRepeatedString(context.TODO(), &pb.CallReq{StringIds: []string{"123", "321"}})
	if err != nil {
		t.Fatalf("test with string_ids err: %v", err)
	}

	if prsp.Rsp != "name_my_name" {
		t.Fatalf("invalid rsp received: %#+v\n", rsp)
	}

	prsp, err = pbc.CallRepeatedInt64(context.TODO(), &pb.CallReq{Int64Ids: []int64{123, 321}})
	if err != nil {
		t.Fatalf("test with int64_ids err: %v", err)
	}

	if prsp.Rsp != "name_my_name" {
		t.Fatalf("invalid rsp received: %#+v\n", rsp)
	}

	// Test-Content-Type

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
		server.Address("127.0.0.1:0"),
		server.Register(reg),
		server.Codec("application/json", jsoncodec.NewCodec()),
		server.Codec("application/x-www-form-urlencoded", urlencodecodec.NewCodec()),
	)

	// create server mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`hello world`))
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

type handlerSwapper struct {
	mu      sync.RWMutex
	handler http.Handler
}

func (h *handlerSwapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	handler := h.handler
	h.mu.RUnlock()
	handler.ServeHTTP(w, r)
}

func TestHTTPServer(t *testing.T) {
	reg := register.NewRegister()
	ctx := context.Background()

	// create server mux
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/first", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`hello world`))
	})
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/second", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`hello world`))
	})

	h := &handlerSwapper{handler: mux1}
	// create server
	srv := httpsrv.NewServer(
		server.Address("127.0.0.1:0"),
		server.Register(reg),
		httpsrv.Server(&http.Server{Handler: h}),
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
	rsp, err := http.Get(fmt.Sprintf("http://%s/first", service[0].Nodes[0].Address))
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

	rsp, err = http.Get(fmt.Sprintf("http://%s/second", service[0].Nodes[0].Address))
	if err != nil {
		t.Fatal(err)
	}
	if rsp.StatusCode != 404 {
		t.Fatal("second route must not exists")
	}
	h.mu.Lock()
	h.handler = mux2
	h.mu.Unlock()

	rsp, err = http.Get(fmt.Sprintf("http://%s/first", service[0].Nodes[0].Address))
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 404 {
		t.Fatal("first route must not exists")
	}

	rsp, err = http.Get(fmt.Sprintf("http://%s/second", service[0].Nodes[0].Address))
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()

	b, err = ioutil.ReadAll(rsp.Body)
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
