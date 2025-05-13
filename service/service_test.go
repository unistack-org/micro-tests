package service

import (
	"testing"

	httpcli "go.unistack.org/micro-client-http/v4"
	httpsrv "go.unistack.org/micro-server-http/v4"
	"go.unistack.org/micro/v4"
	"go.unistack.org/micro/v4/server"
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
