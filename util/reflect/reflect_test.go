package reflect

import (
	"testing"

	pb "github.com/unistack-org/micro-tests/util/reflect/proto"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

func TestFieldName(t *testing.T) {
	name := rutil.FieldName("NestedArgs")
	if name != "nested_args" {
		t.Fatalf("%s != nested_args", name)
	}
}

func TestMergeBool(t *testing.T) {
	type str struct {
		Bool bool `json:"bool"`
	}

	mp := make(map[string]interface{})
	mp["bool"] = "true"
	s := &str{}

	if err := rutil.Merge(s, mp, rutil.Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if !s.Bool {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

	mp["bool"] = "false"

	if err := rutil.Merge(s, mp, rutil.Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if s.Bool {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

	mp["bool"] = 1

	if err := rutil.Merge(s, mp, rutil.Tags([]string{"json"})); err != nil {
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

	t.Logf("merge with true")
	if err := rutil.Merge(s, mp, rutil.Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if s.Bool != "true" {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

	mp["bool"] = false
	t.Logf("merge with false")
	if err := rutil.Merge(s, mp, rutil.Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if s.Bool != "false" {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

}

func TestMerge(t *testing.T) {
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

	if err := rutil.Merge(dst, mp, rutil.Tags([]string{"protobuf"})); err != nil {
		t.Fatal(err)
	}

	if dst.Name != "name_new" || dst.Req != "req_new" || dst.Arg2 != 1 {
		t.Fatalf("merge error: %#+v", dst)
	}

	if dst.Nested == nil || len(dst.Nested.Uint64Args) != 3 ||
		len(dst.Nested.StringArgs) != 2 || dst.Nested.StringArgs[0] != "args1" ||
		len(dst.Nested.Uint64Args) != 3 || dst.Nested.Uint64Args[2].Value != 3 {
		t.Fatalf("merge error: %#+v", dst.Nested)
	}

	nmp := make(map[string]interface{})
	nmp["nested.uint64_args"] = []uint64{4}
	nmp = rutil.FlattenMap(nmp)

	if err := rutil.Merge(dst, nmp, rutil.SliceAppend(true), rutil.Tags([]string{"protobuf"})); err != nil {
		t.Fatal(err)
	}

	if dst.Nested == nil || len(dst.Nested.Uint64Args) != 4 || dst.Nested.Uint64Args[3].Value != 4 {
		t.Fatalf("merge error: %#+v", dst.Nested)
	}
}

func TestMergeNested(t *testing.T) {
	type CallReqNested struct {
		StringArgs []string       `json:"string_args"`
		Uint64Args []uint64       `json:"uint64_args"`
		Nested     *CallReqNested `json:"nested2"`
	}

	type CallReq struct {
		Name   string         `json:"name"`
		Req    string         `json:"req"`
		Arg2   int            `json:"arg2"`
		Nested *CallReqNested `json:"nested"`
	}

	dst := &CallReq{
		Name: "name_old",
		Req:  "req_old",
	}

	mp := make(map[string]interface{})
	mp["name"] = "name_new"
	mp["req"] = "req_new"
	mp["arg2"] = 1
	mp["nested.string_args"] = []string{"args1", "args2"}
	mp["nested.uint64_args"] = []uint64{1, 2, 3}
	mp["nested.nested2.uint64_args"] = []uint64{1, 2, 3}

	mp = rutil.FlattenMap(mp)

	if err := rutil.Merge(dst, mp, rutil.Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if dst.Name != "name_new" || dst.Req != "req_new" || dst.Arg2 != 1 {
		t.Fatalf("merge error: %#+v", dst)
	}

	if dst.Nested == nil || len(dst.Nested.Uint64Args) != 3 ||
		len(dst.Nested.StringArgs) != 2 || dst.Nested.StringArgs[0] != "args1" ||
		len(dst.Nested.Uint64Args) != 3 || dst.Nested.Uint64Args[2] != 3 {
		t.Fatalf("merge error: %#+v", dst.Nested)
	}

	nmp := make(map[string]interface{})
	nmp["nested.uint64_args"] = []uint64{4}
	nmp = rutil.FlattenMap(nmp)

	if err := rutil.Merge(dst, nmp, rutil.SliceAppend(true), rutil.Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if dst.Nested == nil || len(dst.Nested.Uint64Args) != 4 || dst.Nested.Uint64Args[3] != 4 {
		t.Fatalf("merge error: %#+v", dst.Nested)
	}
}
