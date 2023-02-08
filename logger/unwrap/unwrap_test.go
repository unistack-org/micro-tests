package unwrap_test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	pb "go.unistack.org/micro-tests/client/grpc/proto"
	"go.unistack.org/micro/v3/logger/unwrap"
)

func TestProtoMessage(t *testing.T) {
	type Response struct {
		Val *pb.Response `logger:"take"`
	}
	val := &Response{Val: &pb.Response{Msg: "test"}}

	buf := fmt.Sprintf("%#v", unwrap.Unwrap(val, unwrap.Tagged(true)))
	cmp := `&unwrap_test.Response{Val:(*helloworld.Response){Msg:"test"}}`
	if strings.Compare(buf, cmp) != 0 {
		t.Fatalf("not proper written \n%s\n%s", cmp, buf)
	}
}

func TestWrappers(t *testing.T) {
	type CustomerInfo struct {
		MainPhone    *wrappers.StringValue `logger:"take"`
		BankClientId string                `logger:"take"`
		NullString   sql.NullString        `logger:"take"`
	}

	c := &CustomerInfo{MainPhone: &wrappers.StringValue{Value: "+712334"}, BankClientId: "12345", NullString: sql.NullString{String: "test"}}

	buf := fmt.Sprintf("%#v", unwrap.Unwrap(c, unwrap.Tagged(true)))
	cmp := `&unwrap_test.CustomerInfo{MainPhone:(*wrapperspb.StringValue){Value:"+712334"}, BankClientId:"12345", NullString:(sql.NullString){String:"test"}}`
	if strings.Compare(buf, cmp) != 0 {
		t.Fatalf("not proper written \n%s\n%s", cmp, buf)
	}
}
