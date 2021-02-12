package handler

import (
	"context"
	httpsrv "github.com/unistack-org/micro-server-http/v3"
	pb "github.com/unistack-org/micro-tests/client/http/proto"
	"net/http"
)

type GithubHandler struct{}

func NewGithubHandler() *GithubHandler {
	return &GithubHandler{}
}

func (h *GithubHandler) LookupUser(ctx context.Context, req *pb.LookupUserReq, rsp *pb.LookupUserRsp) error {
	if req.GetUsername() == "" || req.GetUsername() != "vtolstov" {
		httpsrv.SetRspCode(ctx, http.StatusBadRequest)
		return &pb.Error{Message: "name is not correct"}
	}
	rsp.Name = "Vasiliy Tolstov"
	httpsrv.SetRspCode(ctx, http.StatusOK)
	return nil
}
