package router_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	rpc "github.com/unistack-org/micro-api-handler-rpc/v3"
	rregister "github.com/unistack-org/micro-api-router-register/v3"
	rstatic "github.com/unistack-org/micro-api-router-static/v3"
	bmemory "github.com/unistack-org/micro-broker-memory/v3"
	gcli "github.com/unistack-org/micro-client-grpc/v3"
	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	protocodec "github.com/unistack-org/micro-codec-proto/v3"
	rmemory "github.com/unistack-org/micro-register-memory/v3"
	regRouter "github.com/unistack-org/micro-router-register/v3"
	gsrv "github.com/unistack-org/micro-server-grpc/v3"
	pb "github.com/unistack-org/micro-tests/server/grpc/proto"
	"github.com/unistack-org/micro/v3/api"
	"github.com/unistack-org/micro/v3/api/handler"
	"github.com/unistack-org/micro/v3/api/router"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/client"
	rt "github.com/unistack-org/micro/v3/router"
	"github.com/unistack-org/micro/v3/server"
)

// server is used to implement helloworld.GreeterServer.
type testServer struct {
}

// TestHello implements helloworld.GreeterServer
func (s *testServer) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	if req.Name == "Timeout" {
		time.Sleep(2 * time.Second)
		rsp.Msg = "Timeout"
		return nil
	}
	rsp.Msg = "Hello " + req.Uuid
	return nil
}

func initial(t *testing.T) (server.Server, client.Client) {
	//logger.DefaultLogger = logger.NewLogger(logger.WithLevel(logger.TraceLevel))
	r := rmemory.NewRegister()
	if err := r.Init(); err != nil {
		t.Fatal(err)
	}

	b := bmemory.NewBroker(broker.Register(r))
	if err := b.Init(); err != nil {
		t.Fatal(err)
	}

	// create a new client
	s := gsrv.NewServer(
		server.Codec("application/grpc+proto", protocodec.NewCodec()),
		server.Codec("application/grpc+json", protocodec.NewCodec()),
		server.Codec("application/json", jsoncodec.NewCodec()),
		server.Name("foo"),
		server.Broker(b),
		server.Register(r),
		server.RegisterInterval(1*time.Second),
	)

	rtr := regRouter.NewRouter(
		rt.Register(r),
	)

	if err := rtr.Init(); err != nil {
		t.Fatal(err)
	}

	// create a new server
	c := gcli.NewClient(
		client.Codec("application/grpc+proto", protocodec.NewCodec()),
		client.Codec("application/grpc+json", protocodec.NewCodec()),
		client.Codec("application/json", jsoncodec.NewCodec()),
		client.Register(r),
		client.Router(rtr),
		client.Broker(b),
	)

	h := &testServer{}
	pb.RegisterTestHandler(s, h)

	if err := s.Init(); err != nil {
		t.Fatalf("failed to init: %v", err)
	}

	if err := s.Start(); err != nil {
		t.Fatalf("failed to start: %v", err)
	}

	return s, c
}

func check(t *testing.T, addr string, path string, expected string, timeout bool) {
	var r io.Reader

	if timeout {
		r = bytes.NewBuffer([]byte(`{"name":"Timeout"}`))
	}
	req, err := http.NewRequest("POST", fmt.Sprintf(path, addr), r)
	if err != nil {
		t.Fatalf("Failed to created http.Request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Timeout", "1") // set timeout to 1s
	rsp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatalf("Failed to created http.Request: %v", err)
	}
	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	jsonMsg := expected
	if string(buf) != jsonMsg {
		t.Fatalf("invalid message received, parsing error %s != %s", buf, jsonMsg)
	}
}

func TestApiTimeout(t *testing.T) {
	s, c := initial(t)
	defer s.Stop()

	router := rregister.NewRouter(
		router.WithHandler(rpc.Handler),
		router.WithRegister(s.Options().Register),
	)
	if err := router.Init(); err != nil {
		t.Fatal(err)
	}

	hrpc := rpc.NewHandler(
		handler.WithClient(c),
		handler.WithRouter(router),
	)
	hsrv := &http.Server{
		Handler:        hrpc,
		Addr:           "127.0.0.1:6543",
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		IdleTimeout:    20 * time.Second,
		MaxHeaderBytes: 1024 * 1024 * 1, // 1Mb
	}

	go func() {
		log.Println(hsrv.ListenAndServe())
	}()

	defer hsrv.Close()
	time.Sleep(1 * time.Second)
	check(t, hsrv.Addr, "http://%s/api/v0/test/call/TEST", `{"Id":"go.micro.client","Code":408,"Detail":"context deadline exceeded","Status":"Request Timeout"}`, true)
}

func TestRouterRegisterPcre(t *testing.T) {
	s, c := initial(t)
	defer s.Stop()

	router := rregister.NewRouter(
		router.WithHandler(rpc.Handler),
		router.WithRegister(s.Options().Register),
	)
	if err := router.Init(); err != nil {
		t.Fatal(err)
	}

	hrpc := rpc.NewHandler(
		handler.WithClient(c),
		handler.WithRouter(router),
	)
	hsrv := &http.Server{
		Handler:        hrpc,
		Addr:           "127.0.0.1:6543",
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		IdleTimeout:    20 * time.Second,
		MaxHeaderBytes: 1024 * 1024 * 1, // 1Mb
	}

	go func() {
		log.Println(hsrv.ListenAndServe())
	}()

	defer hsrv.Close()
	time.Sleep(1 * time.Second)
	check(t, hsrv.Addr, "http://%s/api/v0/test/call/TEST", `{"msg":"Hello "}`, false)
}

func TestRouterStaticPcre(t *testing.T) {
	s, c := initial(t)
	defer s.Stop()

	router := rstatic.NewRouter(
		router.WithHandler(rpc.Handler),
		router.WithRegister(s.Options().Register),
	)
	if err := router.Init(); err != nil {
		t.Fatal(err)
	}

	err := router.Register(&api.Endpoint{
		Name:    "foo.Test.Call",
		Method:  []string{"POST"},
		Path:    []string{"^/api/v0/test/call/?$"},
		Handler: "rpc",
	})
	if err != nil {
		t.Fatal(err)
	}

	hrpc := rpc.NewHandler(
		handler.WithClient(c),
		handler.WithRouter(router),
	)
	hsrv := &http.Server{
		Handler:        hrpc,
		Addr:           "127.0.0.1:6543",
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		IdleTimeout:    20 * time.Second,
		MaxHeaderBytes: 1024 * 1024 * 1, // 1Mb
	}

	go func() {
		log.Println(hsrv.ListenAndServe())
	}()
	defer hsrv.Close()

	time.Sleep(1 * time.Second)
	check(t, hsrv.Addr, "http://%s/api/v0/test/call", `{"msg":"Hello "}`, false)
}

func TestRouterStaticGpath(t *testing.T) {
	s, c := initial(t)
	defer s.Stop()

	router := rstatic.NewRouter(
		router.WithHandler(rpc.Handler),
		router.WithRegister(s.Options().Register),
	)

	err := router.Register(&api.Endpoint{
		Name:    "foo.Test.Call",
		Method:  []string{"POST"},
		Path:    []string{"/api/v0/test/call/{uuid}"},
		Handler: "rpc",
	})
	if err != nil {
		t.Fatal(err)
	}

	hrpc := rpc.NewHandler(
		handler.WithClient(c),
		handler.WithRouter(router),
	)
	hsrv := &http.Server{
		Handler:        hrpc,
		Addr:           "127.0.0.1:6543",
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		IdleTimeout:    20 * time.Second,
		MaxHeaderBytes: 1024 * 1024 * 1, // 1Mb
	}

	go func() {
		log.Println(hsrv.ListenAndServe())
	}()
	defer hsrv.Close()

	time.Sleep(1 * time.Second)
	check(t, hsrv.Addr, "http://%s/api/v0/test/call/TEST", `{"msg":"Hello TEST"}`, false)
}

func TestRouterStaticPcreInvalid(t *testing.T) {
	var ep *api.Endpoint
	var err error

	s, c := initial(t)
	defer s.Stop()

	router := rstatic.NewRouter(
		router.WithHandler(rpc.Handler),
		router.WithRegister(s.Options().Register),
	)

	ep = &api.Endpoint{
		Name:    "foo.Test.Call",
		Method:  []string{"POST"},
		Path:    []string{"^/api/v0/test/call/?"},
		Handler: "rpc",
	}

	err = router.Register(ep)
	if err == nil {
		t.Fatalf("invalid endpoint %v", ep)
	}

	ep = &api.Endpoint{
		Name:    "foo.Test.Call",
		Method:  []string{"POST"},
		Path:    []string{"/api/v0/test/call/?$"},
		Handler: "rpc",
	}

	err = router.Register(ep)
	if err == nil {
		t.Fatalf("invalid endpoint %v", ep)
	}

	_ = c
}
