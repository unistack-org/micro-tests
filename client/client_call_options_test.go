// +build ignore

package client_test

import (
	"context"
	"testing"

	"go.unistack.org/micro/v4/broker"
	bmemory "go.unistack.org/micro/v4/broker/memory"
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/client/grpc"
	tmemory "go.unistack.org/micro/v4/network/transport/memory"
	rmemory "go.unistack.org/micro/v4/register/memory"
	"go.unistack.org/micro/v4/router"
	rtreg "go.unistack.org/micro/v4/router/register"
	"go.unistack.org/micro/v4/server"
	grpcsrv "go.unistack.org/micro/v4/server/grpc"
	cw "go.unistack.org/micro/v4/util/client"
)

type TestFoo struct{}

type TestReq struct{}

type TestRsp struct {
	Data string
}

func (h *TestFoo) Bar(ctx context.Context, req *TestReq, rsp *TestRsp) error {
	rsp.Data = "pass"
	return nil
}

func TestStaticClient(t *testing.T) {
	var err error

	req := grpc.NewClient().NewRequest(
		"go.micro.service.foo",
		"TestFoo.Bar",
		&TestReq{},
		client.WithContentType("application/json"),
	)
	rsp := &TestRsp{}

	reg := rmemory.NewRegister()
	brk := bmemory.NewBroker(broker.Register(reg))
	tr := tmemory.NewTransport()
	rtr := rtreg.NewRouter(router.Register(reg))

	srv := grpcsrv.NewServer(
		server.Broker(brk),
		server.Register(reg),
		server.Name("go.micro.service.foo"),
		server.Address("127.0.0.1:0"),
		server.Transport(tr),
	)
	if err = srv.Handle(srv.NewHandler(&TestFoo{})); err != nil {
		t.Fatal(err)
	}

	if err = srv.Start(); err != nil {
		t.Fatal(err)
	}

	cli := grpc.NewClient(
		client.Router(rtr),
		client.Broker(brk),
		client.Transport(tr),
	)

	w1 := cw.Static("xxx_localhost:12345", cli)
	if err = w1.Call(context.TODO(), req, nil); err == nil {
		t.Fatal("address xxx_#localhost:12345 must not exists and call must be failed")
	}

	w2 := cw.Static(srv.Options().Address, cli)
	if err = w2.Call(context.TODO(), req, rsp); err != nil {
		t.Fatal(err)
	} else if rsp.Data != "pass" {
		t.Fatalf("something wrong with response: %#+v", rsp)
	}
}
