//go:build ignore
// +build ignore

package grpc

import (
	"context"
	"net"
	"testing"
	"time"

	pb "go.unistack.org/micro-tests/client/grpc/proto"
	"google.golang.org/grpc"
	pgrpc "google.golang.org/grpc"
)

func testPool(t *testing.T, size int, ttl time.Duration, idle int, ms int) {
	// setup server
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	ctx := context.Background()
	s := pgrpc.NewServer()
	pb.RegisterTestServer(s, &testServer{})

	go s.Serve(l)
	defer s.Stop()

	// zero pool
	p := newPool(size, ttl, idle, ms)

	for i := 0; i < 10; i++ {
		// get a conn
		cc, err := p.getConn(ctx, l.Addr().String(), grpc.WithInsecure())
		if err != nil {
			t.Fatal(err)
		}

		rsp := pb.Response{}

		err = cc.Invoke(context.TODO(), "/helloworld.Test/CallNative", &pb.Request{Name: "John"}, &rsp)
		if err != nil {
			t.Fatal(err)
		}

		if rsp.Msg != "Hello John" {
			t.Fatalf("Got unexpected response %v", rsp.Msg)
		}

		// release the conn
		p.release(l.Addr().String(), cc, nil)

		p.Lock()
		if i := p.conns[l.Addr().String()].count; i > size {
			p.Unlock()
			t.Fatalf("pool size %d is greater than expected %d", i, size)
		}
		p.Unlock()
	}
}

func TestGRPCPool(t *testing.T) {
	testPool(t, 0, time.Minute, 10, 2)
	testPool(t, 2, time.Minute, 10, 1)
}
