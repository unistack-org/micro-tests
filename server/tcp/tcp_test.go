package tcp_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	tcp "go.unistack.org/micro-server-tcp/v3"
	"go.unistack.org/micro/v3/broker"
	mregister "go.unistack.org/micro/v3/register/memory"
	"go.unistack.org/micro/v3/server"
)

type testHandler struct {
	done  chan struct{}
	cherr chan error
}

func TestTCPServer(t *testing.T) {
	ctx := context.Background()

	reg := mregister.NewRegister()
	if err := reg.Init(); err != nil {
		t.Fatal(err)
	}

	brk := broker.NewBroker(broker.Register(reg))
	if err := brk.Init(); err != nil {
		t.Fatal(err)
	}
	// create server
	srv := tcp.NewServer(server.Register(reg), server.Broker(brk), server.Address("127.0.0.1:65000"))

	// create handler
	h := &testHandler{cherr: make(chan error), done: make(chan struct{})}

	// register handler
	if err := srv.Handle(srv.NewHandler(h)); err != nil {
		t.Fatal(err)
	}

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

	c, err := net.DialTimeout("tcp", srv.Options().Address, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if _, err = c.Write([]byte("test")); err != nil {
		t.Fatal(err)
	}

	select {
	case <-h.done:
		break
	case err := <-h.cherr:
		t.Fatal(err)
	}

	// stop server
	if err := srv.Stop(); err != nil {
		t.Fatal(err)
	}
}

func (h *testHandler) Serve(c net.Conn) {
	var n int
	var err error

	defer c.Close()

	buf := make([]byte, 1024*8) // 8k buffer

	for {
		n, err = c.Read(buf)
		if err != nil && err == io.EOF {
			return
		} else if err != nil {
			h.cherr <- err
		}
		fmt.Printf("%s", buf[:n])
		close(h.done)
	}
}
