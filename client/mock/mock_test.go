package mock

import (
	"context"
	"testing"
	"time"

	"go.unistack.org/micro-client-mock/v3"
	jsoncodec "go.unistack.org/micro-codec-json/v3"
	pb "go.unistack.org/micro-tests/client/mock/proto"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/errors"
)

func TestCallWithoutError(t *testing.T) {
	c := mock.NewClient(client.ContentType("application/json"), client.Codec("application/json", jsoncodec.NewCodec()))

	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	reqbuf := []byte(`{"username": "vtolstov"}`)
	rspbuf := []byte(`{"name": "Vasiliy Tolstov"}`)
	er := c.ExpectRequest(c.NewRequest("github", "Github.LookupUser", reqbuf))
	er.WillReturnResponse("application/json", rspbuf)
	er.WillDelayFor(10 * time.Millisecond)

	gh := pb.NewGithubClient("github", c)

	rsp, err := gh.LookupUser(context.TODO(), &pb.LookupUserReq{Username: "vtolstov"})
	if err != nil {
		t.Fatal(err)
	}

	if rsp.Name != "Vasiliy Tolstov" {
		t.Fatalf("invalid rsp received: %#+v\n", rsp)
	}

	if err := c.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCallWithtError(t *testing.T) {
	c := mock.NewClient(client.ContentType("application/json"), client.Codec("application/json", jsoncodec.NewCodec()))

	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	reqbuf := []byte(`{"username": "vtolstov"}`)
	rspbuf := []byte(`{"name": "Vasiliy Tolstov"}`)
	er := c.ExpectRequest(c.NewRequest("github", "Github.LookupUser", reqbuf))
	er.WillReturnResponse("application/json", rspbuf)
	er.WillDelayFor(10 * time.Millisecond)
	er.WillReturnError(errors.InternalServerError("test", "internal server error"))

	gh := pb.NewGithubClient("github", c)

	rsp, err := gh.LookupUser(context.TODO(), &pb.LookupUserReq{Username: "vtolstov"})
	if err == nil || rsp != nil {
		t.Fatal("call must return error")
	}

	if err := c.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
