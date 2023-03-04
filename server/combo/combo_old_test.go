//go:build ignore

package combo_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	// drpccli "go.unistack.org/micro-client-drpc/v3"
	grpccli "go.unistack.org/micro-client-grpc/v3"
	httpcli "go.unistack.org/micro-client-http/v3"
	jsonpbcodec "go.unistack.org/micro-codec-jsonpb/v3"
	protocodec "go.unistack.org/micro-codec-proto/v3"
	httpsrv "go.unistack.org/micro-server-http/v3"
	// mdpb "go.unistack.org/micro-tests/server/combo/mdpb"
	mgpb "go.unistack.org/micro-tests/server/combo/mgpb"
	mhpb "go.unistack.org/micro-tests/server/combo/mhpb"

	// ndpb "go.unistack.org/micro-tests/server/combo/ndpb"
	ngpb "go.unistack.org/micro-tests/server/combo/ngpb"
	pb "go.unistack.org/micro-tests/server/combo/proto"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
	gmetadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	//"storj.io/drpc"
	//"storj.io/drpc/drpcconn"
	//"storj.io/drpc/drpchttp"
)

type Handler struct {
	t *testing.T
}

const (
	grpcDefaultContentType = "application/grpc+proto"
	// drpcDefaultContentType = "application/drpc+proto"
	httpDefaultContentType = "application/json"
)

type wrapMicroCodec struct{ codec.Codec }

func (w *wrapMicroCodec) Name() string {
	return w.Codec.String()
}

func (w *wrapMicroCodec) Marshal(v interface{}) ([]byte, error) {
	return w.Codec.Marshal(v)
}

func (w *wrapMicroCodec) Unmarshal(d []byte, v interface{}) error {
	return w.Codec.Unmarshal(d, v)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	w.Header().Add("Content-Type", httpDefaultContentType)
	w.WriteHeader(http.StatusOK)
	_ = jsonpbcodec.NewCodec().Write(w, nil, &pb.CallRsp{Rsp: "name_my_name"})
}

func (h *Handler) ServeGRPC(_ interface{}, stream grpc.ServerStream) error {
	ctx := stream.Context()

	fullMethod, ok := grpc.MethodFromServerStream(stream)
	if !ok {
		return status.Errorf(codes.Internal, "method does not exist in context")
	}

	serviceName, methodName, err := grpcServiceMethod(fullMethod)
	if err != nil {
		return status.New(codes.InvalidArgument, err.Error()).Err()
	}
	_, _ = serviceName, methodName
	// get grpc metadata
	gmd, ok := gmetadata.FromIncomingContext(stream.Context())
	if !ok {
		gmd = gmetadata.MD{}
	}

	md := metadata.New(len(gmd))
	for k, v := range gmd {
		md.Set(k, strings.Join(v, ", "))
	}

	// timeout for server deadline
	to, ok := md.Get("timeout")
	if ok {
		md.Del("timeout")
	}

	// get content type
	ct := grpcDefaultContentType

	if ctype, ok := md.Get("content-type"); ok {
		ct = ctype
	} else if ctype, ok := md.Get("x-content-type"); ok {
		ct = ctype
		md.Del("x-content-type")
	}

	_ = ct

	// get peer from context
	if p, ok := peer.FromContext(ctx); ok {
		md["Remote"] = p.Addr.String()
		ctx = peer.NewContext(ctx, p)
	}

	// create new context
	ctx = metadata.NewIncomingContext(ctx, md)

	// set the timeout if we have it
	if len(to) > 0 {
		if n, err := strconv.ParseUint(to, 10, 64); err == nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Duration(n))
			defer cancel()
		}
	}

	frame := &codec.Frame{}
	if err := stream.RecvMsg(frame); err != nil {
		return err
	}

	//	logger.Infof(ctx, "frame: %s", frame.Data)

	if err := stream.SendMsg(&pb.CallRsp{Rsp: "name_my_name"}); err != nil {
		return err
	}

	return nil
}

func grpcServiceMethod(m string) (string, string, error) {
	if len(m) == 0 {
		return "", "", fmt.Errorf("malformed method name: %q", m)
	}

	// grpc method
	if m[0] == '/' {
		// [ , Foo, Bar]
		// [ , package.Foo, Bar]
		// [ , a.package.Foo, Bar]
		parts := strings.Split(m, "/")
		if len(parts) != 3 || len(parts[1]) == 0 || len(parts[2]) == 0 {
			return "", "", fmt.Errorf("malformed method name: %q", m)
		}
		service := strings.Split(parts[1], ".")
		return service[len(service)-1], parts[2], nil
	}

	// non grpc method
	parts := strings.Split(m, ".")

	// expect [Foo, Bar]
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed method name: %q", m)
	}

	return parts[0], parts[1], nil
}

/*
func (h *Handler) ServeDRPC(stream drpc.Stream, rpc string) error {
	ctx := stream.Context()
	logger.Infof(ctx, "drpc: %#+v", rpc)
	return nil
}
*/

/*
func (h *Handler) HandleRPC(stream drpc.Stream, rpc string) error {
	return h.ServeDRPC(stream, rpc)
}
*/

func TestComboServer(t *testing.T) {
	reg := register.NewRegister()
	ctx := context.Background()

	h := &Handler{t: t}

	_ = logger.DefaultLogger.Init(logger.WithCallerSkipCount(3))
	encoding.RegisterCodec(&wrapMicroCodec{protocodec.NewCodec()})

	lis, err := net.Listen("tcp", fmt.Sprintf(":0"))
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	gsrv := grpc.NewServer(grpc.UnknownServiceHandler(h.ServeGRPC))
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := gsrv.Serve(lis); err != nil {
					log.Fatalf("failed to serve: %v", err)
				}
			}
		}
	}()

	comboHandler := newComboMux(h, gsrv, nil) // drpchttp.New(h))
	http2Server := &http2.Server{}
	hs := &http.Server{Handler: h2c.NewHandler(comboHandler, http2Server)}

	// create server
	srv := httpsrv.NewServer(
		server.Address("127.0.0.1:0"),
		server.Name("helloworld"),
		server.Register(reg),
		httpsrv.Server(hs),
	)

	// init server
	if err := srv.Init(); err != nil {
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

	mhcli := client.NewClientCallOptions(httpcli.NewClient(client.ContentType(httpDefaultContentType), client.Codec(httpDefaultContentType, jsonpbcodec.NewCodec())), client.WithAddress("http://"+service[0].Nodes[0].Address))

	mhttpsvc := mhpb.NewTestClient("helloworld", mhcli)

	mgcli := client.NewClientCallOptions(grpccli.NewClient(client.ContentType(grpcDefaultContentType), client.Codec(grpcDefaultContentType, protocodec.NewCodec())), client.WithAddress("http://"+service[0].Nodes[0].Address))

	mgrpcsvc := mgpb.NewTestClient("helloworld", mgcli)

	// mdcli := client.NewClientCallOptions(drpccli.NewClient(client.ContentType(drpcDefaultContentType), client.Codec(drpcDefaultContentType, protocodec.NewCodec())), client.WithAddress("http://"+service[0].Nodes[0].Address))

	// mdrpcsvc := mdpb.NewTestClient("helloworld", mdcli)

	t.Logf("call via micro grpc")
	rsp, err := mgrpcsvc.Call(ctx, &pb.CallReq{Req: "my_name"})
	if err != nil {
		t.Fatal(err)
	} else {
		if rsp.Rsp != "name_my_name" {
			t.Fatalf("invalid response: %#+v\n", rsp)
		}
	}

	ngcli, err := grpc.DialContext(ctx, service[0].Nodes[0].Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer ngcli.Close()

	ngrpcsvc := ngpb.NewTestClient(ngcli)
	t.Logf("call via native grpc")
	if rsp, err := ngrpcsvc.Call(ctx, &ngpb.CallReq{Req: "my_name"}); err != nil {
		t.Fatal(err)
	} else {
		if rsp.Rsp != "name_my_name" {
			t.Fatalf("invalid response: %#+v\n", rsp)
		}
	}

	t.Logf("call via micro http")
	if rsp, err := mhttpsvc.Call(ctx, &pb.CallReq{Req: "my_name"}); err != nil {
		t.Fatal(err)
	} else {
		if rsp.Rsp != "name_my_name" {
			t.Fatalf("invalid response: %#+v\n", rsp)
		}
	}

	/*
		tc, err := net.Dial("tcp", service[0].Nodes[0].Address)
		if err != nil {
			t.Fatal(err)
		}

		ndcli := drpcconn.New(tc)
		defer ndcli.Close()
		/*
			ndrpcsvc := ndpb.NewDRPCTestClient(ndcli)

			t.Logf("call via native drpc")
			if rsp, err := ndrpcsvc.Call(context.TODO(), &ndpb.CallReq{Req: "my_name"}); err != nil {
				t.Logf("native drpc err: %v", err)
				// t.Fatal(err)
			} else {
				if rsp.Rsp != "name_my_name" {
					t.Fatalf("invalid response: %#+v\n", rsp)
				}
			}

		t.Logf("call via micro drpc")
		if rsp, err = mdrpcsvc.Call(ctx, &pb.CallReq{Req: "my_name"}); err != nil {
			t.Logf("micro drpc err: %v", err)
			// t.Fatal(err)
		} else {
			if rsp.Rsp != "name_my_name" {
				t.Fatalf("invalid response: %#+v\n", rsp)
			}
		}*/
}

func newComboMux(httph http.Handler, grpch http.Handler, drpch http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 {
			ct := r.Header.Get("content-type")
			switch {
			case strings.HasPrefix(ct, "application/grpc"):
				if grpch != nil {
					grpch.ServeHTTP(w, r)
				}
				return
			case strings.HasPrefix(ct, "application/drpc"):
				if drpch != nil {
					drpch.ServeHTTP(w, r)
				}
				return
			}
		}

		httph.ServeHTTP(w, r)
	})
}
