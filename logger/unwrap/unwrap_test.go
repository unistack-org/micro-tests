package unwrap_test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	"go.unistack.org/micro/v3/logger/unwrap"
)

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
