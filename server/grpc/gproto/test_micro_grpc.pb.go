// Code generated by protoc-gen-go-micro. DO NOT EDIT.
// protoc-gen-go-micro version: v3.5.3
// source: test.proto

package helloworld

import (
	context "context"
	proto "go.unistack.org/micro-tests/server/grpc/proto"
	api "go.unistack.org/micro/v3/api"
	client "go.unistack.org/micro/v3/client"
	server "go.unistack.org/micro/v3/server"
)

type testClient struct {
	c    client.Client
	name string
}

func NewTestClient(name string, c client.Client) TestClient {
	return &testClient{c: c, name: name}
}

func (c *testClient) Call(ctx context.Context, req *proto.Request, opts ...client.CallOption) (*proto.Response, error) {
	rsp := &proto.Response{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "Test.Call", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type testServer struct {
	TestServer
}

func (h *testServer) Call(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	return h.TestServer.Call(ctx, req, rsp)
}

func RegisterTestServer(s server.Server, sh TestServer, opts ...server.HandlerOption) error {
	type test interface {
		Call(ctx context.Context, req *proto.Request, rsp *proto.Response) error
	}
	type Test struct {
		test
	}
	h := &testServer{sh}
	var nopts []server.HandlerOption
	for _, endpoint := range TestEndpoints {
		nopts = append(nopts, api.WithEndpoint(&endpoint))
	}
	return s.Handle(s.NewHandler(&Test{h}, append(nopts, opts...)...))
}