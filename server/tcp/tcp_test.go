package tcp_test

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	bmemory "github.com/unistack-org/micro-broker-memory"
	rmemory "github.com/unistack-org/micro-registry-memory"
	tcp "github.com/unistack-org/micro-server-tcp"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/server"
)

type testHandler struct {
	done chan struct{}
}

func TestTCPServer(t *testing.T) {
	reg := rmemory.NewRegistry()
	brk := bmemory.NewBroker(broker.Registry(reg))
	// create server
	srv := tcp.NewServer(server.Registry(reg), server.Broker(brk), server.Address("127.0.0.1:65000"))

	// create handler
	h := &testHandler{done: make(chan struct{})}

	// register handler
	if err := srv.Handle(srv.NewHandler(h)); err != nil {
		t.Fatal(err)
	}

	// start server
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}

	// lookup server
	service, err := reg.GetService(server.DefaultName)
	if err != nil {
		t.Fatal(err)
	}

	if len(service) != 1 {
		t.Fatalf("Expected 1 service got %d: %+v", len(service), service)
	}

	if len(service[0].Nodes) != 1 {
		t.Fatalf("Expected 1 node got %d: %+v", len(service[0].Nodes), service[0].Nodes)
	}

	go func() {
		<-h.done
		// stop server
		if err := srv.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	c, err := net.DialTimeout("tcp", srv.Options().Address, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if _, err = c.Write([]byte("test")); err != nil {
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
			logger.Fatal(err)
		}
		fmt.Printf("%s", buf[:n])
		close(h.done)
	}
}
