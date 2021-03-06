// Code generated by protoc-gen-micro
// source: test.proto
package pb

import (
	context "context"
	api "github.com/unistack-org/micro/v3/api"
	client "github.com/unistack-org/micro/v3/client"
	codec "github.com/unistack-org/micro/v3/codec"
)

func NewTestServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{
		&api.Endpoint{
			Name:    "TestService.TestEndpoint",
			Path:    []string{"/users/test"},
			Method:  []string{"GET"},
			Handler: "rpc",
		},
		&api.Endpoint{
			Name:    "TestService.UserByID",
			Path:    []string{"/users/{id}"},
			Method:  []string{"GET"},
			Handler: "rpc",
		},
		&api.Endpoint{
			Name:    "TestService.UserImageByID",
			Path:    []string{"/users/{id}/image"},
			Method:  []string{"GET"},
			Handler: "rpc",
		},
		&api.Endpoint{
			Name:    "TestService.UploadFile",
			Path:    []string{"/users/image/upload"},
			Method:  []string{"POST"},
			Handler: "rpc",
		},
		&api.Endpoint{
			Name:    "TestService.KzAmlRs",
			Path:    []string{"/aml"},
			Method:  []string{"POST"},
			Handler: "rpc",
		},
	}
}

type TestServiceClient interface {
	TestEndpoint(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error)
	UserByID(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error)
	UserImageByID(ctx context.Context, req *Request, opts ...client.CallOption) (*codec.Frame, error)
	UploadFile(ctx context.Context, req *RequestImage, opts ...client.CallOption) (*ResponseImage, error)
	KzAmlRs(ctx context.Context, req *RequestAml, opts ...client.CallOption) (*ResponseAml, error)
}

type TestServiceServer interface {
	TestEndpoint(ctx context.Context, req *Request, rsp *Response) error
	UserByID(ctx context.Context, req *Request, rsp *Response) error
	UserImageByID(ctx context.Context, req *Request, rsp *codec.Frame) error
	UploadFile(ctx context.Context, req *RequestImage, rsp *ResponseImage) error
	KzAmlRs(ctx context.Context, req *RequestAml, rsp *ResponseAml) error
}
