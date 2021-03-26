package service

import (
	"testing"

	httpcli "github.com/unistack-org/micro-client-http/v3"
	httpsrv "github.com/unistack-org/micro-server-http/v3"
	"github.com/unistack-org/micro/v3"
	"github.com/unistack-org/micro/v3/server"
)

func TestHTTPService(t *testing.T) {
	svc := micro.NewService(
		micro.Server(httpsrv.NewServer(server.Address("127.0.0.1:0"))),
		micro.Client(httpcli.NewClient()),
	)

	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}

}
