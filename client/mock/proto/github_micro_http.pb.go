package pb

import (
	context "context"
	http "net/http"

	v3 "go.unistack.org/micro-client-http/v3"
	v31 "go.unistack.org/micro-server-http/v3"
	client "go.unistack.org/micro/v3/client"
	server "go.unistack.org/micro/v3/server"
)

type githubClient struct {
	c    client.Client
	name string
}

func NewGithubClient(name string, c client.Client) GithubClient {
	return &githubClient{c: c, name: name}
}

func (c *githubClient) LookupUser(ctx context.Context, req *LookupUserReq, opts ...client.CallOption) (*LookupUserRsp, error) {
	errmap := make(map[string]interface{}, 1)
	errmap["default"] = &Error{}
	opts = append(opts,
		v3.ErrorMap(errmap),
	)
	opts = append(opts,
		v3.Method(http.MethodGet),
		v3.Path("/users/{username}"),
	)
	rsp := &LookupUserRsp{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "Github.LookupUser", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type githubServer struct {
	GithubServer
}

func (h *githubServer) LookupUser(ctx context.Context, req *LookupUserReq, rsp *LookupUserRsp) error {
	return h.GithubServer.LookupUser(ctx, req, rsp)
}

func RegisterGithubServer(s server.Server, sh GithubServer, opts ...server.HandlerOption) error {
	type github interface {
		LookupUser(ctx context.Context, req *LookupUserReq, rsp *LookupUserRsp) error
	}
	type Github struct {
		github
	}
	h := &githubServer{sh}
	var nopts []server.HandlerOption
	nopts = append(nopts, v31.HandlerEndpoints(GithubServerEndpoints))
	return s.Handle(s.NewHandler(&Github{h}, append(nopts, opts...)...))
}
