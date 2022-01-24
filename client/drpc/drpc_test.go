// +build ignore

package drpc

/*
import (
	"context"
	"net"
	"testing"

	drpc "go.unistack.org/micro-client-drpc/v3"
	protocodec "go.unistack.org/micro-codec-proto/v3"
	pb "go.unistack.org/micro-tests/client/drpc/proto"
	"go.unistack.org/micro/v3/client"
	"storj.io/drpc/drpcmux"
	"storj.io/drpc/drpcserver"
)

type TestServer struct{}

func (s *TestServer) Call(ctx context.Context, req *pb.CallReq) (*pb.CallRsp, error) {
	return &pb.CallRsp{Name: req.Name + " rsp"}, nil
}

func (s *TestServer) Hello(ctx context.Context, req *pb.CallReq) (*pb.CallRsp, error) {
	return &pb.CallRsp{Name: req.Name + " rsp"}, nil
}

func TestDrpc(t *testing.T) {
	t.Skip()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := drpcmux.New()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	err = pb.DRPCRegisterTestService(m, &TestServer{})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		if err := drpcserver.New(m).Serve(ctx, ln); err != nil {
			t.Fatal(err)
		}
	}()

	c := client.NewClientCallOptions(drpc.NewClient(client.Codec("application/drpc+proto", protocodec.NewCodec())), client.WithAddress(ln.Addr().String()))
	cli := pb.NewTestServiceClient("test", c)
	rsp, err := cli.Call(ctx, &pb.CallReq{Name: "test_name"})
	if err != nil {
		t.Fatal(err)
	}
	if rsp.Name != "test_name rsp" {
		t.Fatalf("unexpected rsp %#+v", rsp)
	}
}
*/
