package combo_test

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"testing"

	grpccli "go.unistack.org/micro-client-grpc/v3"
	httpcli "go.unistack.org/micro-client-http/v3"
	jsonpbcodec "go.unistack.org/micro-codec-jsonpb/v3"
	protocodec "go.unistack.org/micro-codec-proto/v3"
	grpcsrv "go.unistack.org/micro-server-grpc/v3"
	httpsrv "go.unistack.org/micro-server-http/v3"
	mgpb "go.unistack.org/micro-tests/server/combo/mgpb"
	mhpb "go.unistack.org/micro-tests/server/combo/mhpb"
	ngpb "go.unistack.org/micro-tests/server/combo/ngpb"
	pb "go.unistack.org/micro-tests/server/combo/proto"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed swagger-ui
var assets embed.FS

type Handler struct {
	t *testing.T
}

const (
	grpcDefaultContentType = "application/grpc+proto"
	httpDefaultContentType = "application/json"
)

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

func (h *Handler) Call(ctx context.Context, req *pb.CallReq, rsp *pb.CallRsp) error {
	rsp.Rsp = "name_my_name"
	return nil
}

func TestComboServer(t *testing.T) {
	reg := register.NewRegister()
	ctx := context.Background()

	h := &Handler{t: t}

	_ = logger.DefaultLogger.Init(logger.WithCallerSkipCount(3))

	// create grpc server
	gsrv := grpcsrv.NewServer(
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("application/json", jsonpbcodec.NewCodec()),
		server.Codec("application/grpc", protocodec.NewCodec()),
		server.Codec("application/grpc+proto", protocodec.NewCodec()),
		server.Codec("application/grpc+json", jsonpbcodec.NewCodec()),
	)

	// init grpc server
	if err := gsrv.Init(); err != nil {
		t.Fatalf("grpc err: %v", err)
	}

	if err := mgpb.RegisterTestServer(gsrv, h); err != nil {
		t.Fatalf("grpc err: %v", err)
	}

	swaggerdir, _ := fs.Sub(assets, "swagger-ui")

	// create http server
	hsrv := httpsrv.NewServer(
		server.Address("127.0.0.1:0"),
		server.Name("helloworld"),
		server.Register(reg),
		server.Codec("application/json", jsonpbcodec.NewCodec()),
		httpsrv.PathHandler(http.MethodGet, "/swagger-ui/*", http.StripPrefix("/swagger-ui", http.FileServer(http.FS(swaggerdir))).ServeHTTP),
	)

	// fill http server handler struct
	hs := &http.Server{Handler: h2c.NewHandler(newComboMux(hsrv, gsrv.GRPCServer(), nil), &http2.Server{})}

	// init http server
	if err := hsrv.Init(httpsrv.Server(hs)); err != nil {
		t.Fatal(err)
	}

	if err := mhpb.RegisterTestServer(hsrv, h); err != nil {
		t.Fatalf("grpc err: %v", err)
	}

	// start http server
	if err := hsrv.Start(); err != nil {
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

	t.Logf("call via micro grpc")
	if rsp, err := mgrpcsvc.Call(ctx, &pb.CallReq{Req: "my_name"}); err != nil {
		t.Fatal(err)
	} else if rsp.Rsp != "name_my_name" {
		t.Fatalf("invalid response: %#+v\n", rsp)
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
	} else if rsp.Rsp != "name_my_name" {
		t.Fatalf("invalid response: %#+v\n", rsp)
	}

	t.Logf("call via micro http")
	if rsp, err := mhttpsvc.Call(ctx, &pb.CallReq{Req: "my_name"}); err != nil {
		t.Fatal(err)
	} else if rsp.Rsp != "name_my_name" {
		t.Fatalf("invalid response: %#+v\n", rsp)
	}

	var hreq *http.Request
	hreq, err = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/swagger-ui/index.html", service[0].Nodes[0].Address), nil)
	if err != nil {
		t.Fatal(err)
	}
	var hrsp *http.Response
	hrsp, err = http.DefaultClient.Do(hreq)
	if err != nil || hrsp.StatusCode != http.StatusOK {
		t.Fatalf("error rsp: %v err: %v", hrsp, err)
	}

	hreq, err = http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/Call", service[0].Nodes[0].Address), strings.NewReader(`{"req":"my_name"}`))
	if err != nil {
		t.Fatal(err)
	}
	hreq.Header.Add("Content-Type", "application/json")

	hrsp, err = http.DefaultClient.Do(hreq)
	if err != nil || hrsp.StatusCode != http.StatusOK {
		t.Fatalf("error rsp: %v err: %v", hrsp, err)
	}
	defer hrsp.Body.Close()
	buf, err := io.ReadAll(hrsp.Body)
	if err != nil {
		t.Fatalf("read body fail: %v", err)
	}
	rsp := &pb.CallRsp{}
	if err = jsonpbcodec.NewCodec().Unmarshal(buf, rsp); err != nil {
		t.Fatal(err)
	} else if rsp.Rsp != "name_my_name" {
		t.Fatalf("invalid response: %#+v\n", rsp)
	}
}
