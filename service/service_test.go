package service

import (
	"testing"

	httpcli "go.unistack.org/micro-client-http/v5"
	httpsrv "go.unistack.org/micro-server-http/v5"
	"go.unistack.org/micro/v5"
	"go.unistack.org/micro/v5/server"
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
