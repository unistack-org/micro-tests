package pb

import (
	context "context"

	v3 "go.unistack.org/micro-server-http/v3"
	client "go.unistack.org/micro/v3/client"
)

var GithubName = "Github"

var GithubServerEndpoints = []v3.EndpointMetadata{
	{
		Name:   "Github.LookupUser",
		Path:   "/users/{username}",
		Method: "GET",
		Body:   "",
		Stream: false,
	},
}

type GithubClient interface {
	LookupUser(ctx context.Context, req *LookupUserReq, opts ...client.CallOption) (*LookupUserRsp, error)
}

type GithubServer interface {
	LookupUser(ctx context.Context, req *LookupUserReq, rsp *LookupUserRsp) error
}
