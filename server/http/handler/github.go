package handler

import (
	"context"
	"net/http"

	httpsrv "go.unistack.org/micro-server-http/v3"
	pb "go.unistack.org/micro-tests/client/http/proto"
)

type GithubHandler struct{}

func NewGithubHandler() *GithubHandler {
	return &GithubHandler{}
}

func (h *GithubHandler) LookupUser(ctx context.Context, req *pb.LookupUserReq, rsp *pb.LookupUserRsp) error {
	if req.GetUsername() == "" || req.GetUsername() != "vtolstov" {
		httpsrv.SetRspCode(ctx, http.StatusBadRequest)
		return httpsrv.SetError(&pb.Error{Message: "name is not correct"})
	}
	rsp.Name = "Vasiliy Tolstov"
	httpsrv.SetRspCode(ctx, http.StatusOK)
	return nil
}
