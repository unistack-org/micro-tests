// Code generated by protoc-gen-go-micro. DO NOT EDIT.
// protoc-gen-go-micro version: v3.5.3
// source: test.proto

package pb

import (
	context "context"
	v3 "go.unistack.org/micro-client-http/v3"
	v31 "go.unistack.org/micro-server-http/v3"
	api "go.unistack.org/micro/v3/api"
	client "go.unistack.org/micro/v3/client"
	codec "go.unistack.org/micro/v3/codec"
	server "go.unistack.org/micro/v3/server"
	http "net/http"
)

type testServiceClient struct {
	c    client.Client
	name string
}

func NewTestServiceClient(name string, c client.Client) TestServiceClient {
	return &testServiceClient{c: c, name: name}
}

func (c *testServiceClient) TestEndpoint(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error) {
	opts = append(opts,
		v3.Method(http.MethodGet),
		v3.Path("/users/test"),
	)
	opts = append(opts,
		v3.Header("client_uid", "true"),
		v3.Cookie("csrftoken", "true"),
	)
	rsp := &Response{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "TestService.TestEndpoint", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *testServiceClient) UserByID(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error) {
	opts = append(opts,
		v3.Method(http.MethodGet),
		v3.Path("/users/{id}"),
	)
	rsp := &Response{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "TestService.UserByID", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *testServiceClient) UserImageByID(ctx context.Context, req *Request, opts ...client.CallOption) (*codec.Frame, error) {
	opts = append(opts,
		v3.Method(http.MethodGet),
		v3.Path("/users/{id}/image"),
	)
	rsp := &codec.Frame{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "TestService.UserImageByID", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *testServiceClient) UploadFile(ctx context.Context, req *RequestImage, opts ...client.CallOption) (*ResponseImage, error) {
	opts = append(opts,
		v3.Method(http.MethodPost),
		v3.Path("/users/image/upload"),
	)
	rsp := &ResponseImage{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "TestService.UploadFile", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *testServiceClient) KzAmlRs(ctx context.Context, req *RequestAml, opts ...client.CallOption) (*ResponseAml, error) {
	opts = append(opts,
		v3.Method(http.MethodPost),
		v3.Path("/aml"),
	)
	rsp := &ResponseAml{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "TestService.KzAmlRs", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type testServiceServer struct {
	TestServiceServer
}

func (h *testServiceServer) TestEndpoint(ctx context.Context, req *Request, rsp *Response) error {
	v31.FillRequest(ctx, req,
		v31.Header("client_uid", "true"),
		v31.Cookie("csrftoken", "true"),
	)
	return h.TestServiceServer.TestEndpoint(ctx, req, rsp)
}

func (h *testServiceServer) UserByID(ctx context.Context, req *Request, rsp *Response) error {
	return h.TestServiceServer.UserByID(ctx, req, rsp)
}

func (h *testServiceServer) UserImageByID(ctx context.Context, req *Request, rsp *codec.Frame) error {
	return h.TestServiceServer.UserImageByID(ctx, req, rsp)
}

func (h *testServiceServer) UploadFile(ctx context.Context, req *RequestImage, rsp *ResponseImage) error {
	return h.TestServiceServer.UploadFile(ctx, req, rsp)
}

func (h *testServiceServer) KzAmlRs(ctx context.Context, req *RequestAml, rsp *ResponseAml) error {
	return h.TestServiceServer.KzAmlRs(ctx, req, rsp)
}

func RegisterTestServiceServer(s server.Server, sh TestServiceServer, opts ...server.HandlerOption) error {
	type testService interface {
		TestEndpoint(ctx context.Context, req *Request, rsp *Response) error
		UserByID(ctx context.Context, req *Request, rsp *Response) error
		UserImageByID(ctx context.Context, req *Request, rsp *codec.Frame) error
		UploadFile(ctx context.Context, req *RequestImage, rsp *ResponseImage) error
		KzAmlRs(ctx context.Context, req *RequestAml, rsp *ResponseAml) error
	}
	type TestService struct {
		testService
	}
	h := &testServiceServer{sh}
	var nopts []server.HandlerOption
	for _, endpoint := range TestServiceEndpoints {
		nopts = append(nopts, api.WithEndpoint(&endpoint))
	}
	return s.Handle(s.NewHandler(&TestService{h}, append(nopts, opts...)...))
}