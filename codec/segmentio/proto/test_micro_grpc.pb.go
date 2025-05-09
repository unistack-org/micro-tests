// Code generated by protoc-gen-go-micro. DO NOT EDIT.
// protoc-gen-go-micro version: v3.10.4
// source: test.proto

package pb

import (
	context "context"
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

func (c *testClient) Call(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error) {
	rsp := &Response{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "Test.Call", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type testServer struct {
	TestServer
}

func (h *testServer) Call(ctx context.Context, req *Request, rsp *Response) error {
	return h.TestServer.Call(ctx, req, rsp)
}

func RegisterTestServer(s server.Server, sh TestServer, opts ...server.HandlerOption) error {
	type test interface {
		Call(ctx context.Context, req *Request, rsp *Response) error
	}
	type Test struct {
		test
	}
	h := &testServer{sh}
	return s.Handle(s.NewHandler(&Test{h}, opts...))
}
