// Code generated by protoc-gen-micro
// source: test.proto
package helloworld

import (
	context "context"
	proto "github.com/unistack-org/micro-tests/client/grpc/proto"
	api "github.com/unistack-org/micro/v3/api"
	client "github.com/unistack-org/micro/v3/client"
	server "github.com/unistack-org/micro/v3/server"
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

func (c *testClient) Stream(ctx context.Context, opts ...client.CallOption) (Test_StreamClient, error) {
	stream, err := c.c.Stream(ctx, c.c.NewRequest(c.name, "Test.Stream", &proto.Request{}), opts...)
	if err != nil {
		return nil, err
	}
	return &testClientStream{stream}, nil
}

type testClientStream struct {
	stream client.Stream
}

func (s *testClientStream) Close() error {
	return s.stream.Close()
}

func (s *testClientStream) Context() context.Context {
	return s.stream.Context()
}

func (s *testClientStream) SendMsg(msg interface{}) error {
	return s.stream.Send(msg)
}

func (s *testClientStream) RecvMsg(msg interface{}) error {
	return s.stream.Recv(msg)
}

func (s *testClientStream) Send(msg *proto.Request) error {
	return s.stream.Send(msg)
}

func (s *testClientStream) Recv() (*proto.Response, error) {
	msg := &proto.Response{}
	if err := s.stream.Recv(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

type testServer struct {
	TestServer
}

func (h *testServer) Call(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	return h.TestServer.Call(ctx, req, rsp)
}

func (h *testServer) Stream(ctx context.Context, stream server.Stream) error {
	return h.TestServer.Stream(ctx, &testStreamStream{stream})
}

type testStreamStream struct {
	stream server.Stream
}

func (s *testStreamStream) Close() error {
	return s.stream.Close()
}

func (s *testStreamStream) Context() context.Context {
	return s.stream.Context()
}

func (s *testStreamStream) SendMsg(msg interface{}) error {
	return s.stream.Send(msg)
}

func (s *testStreamStream) RecvMsg(msg interface{}) error {
	return s.stream.Recv(msg)
}

func (s *testStreamStream) Send(msg *proto.Response) error {
	return s.stream.Send(msg)
}

func (s *testStreamStream) Recv() (*proto.Request, error) {
	msg := &proto.Request{}
	if err := s.stream.Recv(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func RegisterTestServer(s server.Server, sh TestServer, opts ...server.HandlerOption) error {
	type test interface {
		Call(ctx context.Context, req *proto.Request, rsp *proto.Response) error
		Stream(ctx context.Context, stream server.Stream) error
	}
	type Test struct {
		test
	}
	h := &testServer{sh}
	for _, endpoint := range NewTestEndpoints() {
		opts = append(opts, api.WithEndpoint(endpoint))
	}
	return s.Handle(s.NewHandler(&Test{h}, opts...))
}