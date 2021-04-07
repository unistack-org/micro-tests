package reflect

import (
	"testing"

	pb "github.com/unistack-org/micro-tests/util/reflect/proto"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

func TestMergeBool(t *testing.T) {
	type str struct {
		Bool bool `json:"bool"`
	}

	mp := make(map[string]interface{})
	mp["bool"] = "true"
	s := &str{}

	if err := MergeMap(s, mp, []string{"json"}); err != nil {
		t.Fatal(err)
	}

	if !s.Bool {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

	mp["bool"] = "false"

	if err := MergeMap(s, mp, []string{"json"}); err != nil {
		t.Fatal(err)
	}

	if s.Bool {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

	mp["bool"] = 1

	if err := MergeMap(s, mp, []string{"json"}); err != nil {
		t.Fatal(err)
	}

	if !s.Bool {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

}

func TestMergeString(t *testing.T) {
	type str struct {
		Bool string `json:"bool"`
	}

	mp := make(map[string]interface{})
	mp["bool"] = true
	s := &str{}

	if err := MergeMap(s, mp, []string{"json"}); err != nil {
		t.Fatal(err)
	}

	if s.Bool != "true" {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

	mp["bool"] = false

	if err := MergeMap(s, mp, []string{"json"}); err != nil {
		t.Fatal(err)
	}

	if s.Bool != "false" {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

}

func TestMergeMap(t *testing.T) {
	dst := &pb.CallReq{
		Name: "name_old",
		Req:  "req_old",
	}

	mp := make(map[string]interface{})
	mp["name"] = "name_new"
	mp["req"] = "req_new"
	mp["arg2"] = 1
	mp["nested.string_args"] = []string{"args1", "args2"}
	mp["nested.uint64_args"] = []uint64{1, 2, 3}

	mp = rutil.FlattenMap(mp)

	if err := MergeMap(dst, mp, []string{"protobuf"}); err != nil {
		t.Fatal(err)
	}

	if dst.Name != "name_new" || dst.Req != "req_new" || dst.Arg2 != 1 {
		t.Fatalf("merge error: %v", dst)
	}

	if dst.Nested == nil ||
		len(dst.Nested.StringArgs) != 2 || dst.Nested.StringArgs[0] != "args1" ||
		len(dst.Nested.Uint64Args) != 3 || dst.Nested.Uint64Args[2].Value != 3 {
		t.Fatalf("merge error: %v", dst.Nested)
	}
}
