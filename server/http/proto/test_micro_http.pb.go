// Code generated by protoc-gen-micro
// source: test.proto
package pb

import (
	"context"
	"fmt"
	//  "net/http"

	micro_client_http "github.com/unistack-org/micro-client-http/v3"
	micro_api "github.com/unistack-org/micro/v3/api"
	micro_client "github.com/unistack-org/micro/v3/client"
	micro_server "github.com/unistack-org/micro/v3/server"
)

var (
	_ micro_server.Option
	_ micro_client.Option
)

type testService struct {
	c    micro_client.Client
	name string
}

// Micro client stuff

// NewTestService create new service client
func NewTestService(name string, c micro_client.Client) TestService {
	return &testService{c: c, name: name}
}

func (c *testService) Call(ctx context.Context, req *CallReq, opts ...micro_client.CallOption) (*CallRsp, error) {
	errmap := make(map[string]interface{}, 1)
	errmap["default"] = &Error{}

	nopts := append(opts,
		micro_client_http.Method("POST"),
		micro_client_http.Path("/v1/test/call/{name}"),
		micro_client_http.Body("*"),
		micro_client_http.ErrorMap(errmap),
	)
	rsp := &CallRsp{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "Test.Call", req), rsp, nopts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *testService) CallError(ctx context.Context, req *CallReq1, opts ...micro_client.CallOption) (*CallRsp1, error) {
	errmap := make(map[string]interface{}, 1)
	errmap["default"] = &Error{}

	nopts := append(opts,
		micro_client_http.Method("POST"),
		micro_client_http.Path("/v1/test/callerror/{name}"),
		micro_client_http.Body("*"),
		micro_client_http.ErrorMap(errmap),
	)
	rsp := &CallRsp1{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "Test.CallError", req), rsp, nopts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// Micro server stuff

type testHandler struct {
	TestHandler
}

func (h *testHandler) Call(ctx context.Context, req *CallReq, rsp *CallRsp) error {
	return h.TestHandler.Call(ctx, req, rsp)
}

/*
func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("new request: %#+v\n", r)
  // HANDLE ALL STUFF
}
*/

func (h *testHandler) CallError(ctx context.Context, req *CallReq1, rsp *CallRsp1) error {
	return h.TestHandler.CallError(ctx, req, rsp)
}

/*
func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("new request: %#+v\n", r)
  // HANDLE ALL STUFF
}
*/

// Error method to satisfy error interface
func (e *Error) Error() string {
	return fmt.Sprintf("%#v", e)
}

// RegisterTestHandler registers server handler
func RegisterTestHandler(s micro_server.Server, sh TestHandler, opts ...micro_server.HandlerOption) error {
	type test interface {
		Call(context.Context, *CallReq, *CallRsp) error
		CallError(context.Context, *CallReq1, *CallRsp1) error
		//        ServeHTTP(http.ResponseWriter, *http.Request)
	}
	type Test struct {
		test
	}
	h := &testHandler{sh}
	for _, endpoint := range NewTestEndpoints() {
		opts = append(opts, micro_api.WithEndpoint(endpoint))
	}
	return s.Handle(s.NewHandler(&Test{h}, opts...))
}