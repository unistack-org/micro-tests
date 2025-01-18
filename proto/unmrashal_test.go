package pb

import (
	"testing"

	cp "go.unistack.org/micro-codec-proto/v3"
)

func TestMarshalUnmarshal(t *testing.T) {
	c := cp.NewCodec()
	msg2 := &Message2{Items: []*Item2{{Key1: "akey1", Key2: "akey2"}, {Key1: "bkey1", Key2: "bkey2"}}}
	buf, err := c.Marshal(msg2)
	if err != nil {
		t.Fatal(err)
	}

	msg1 := &Message1{}
	err = c.Unmarshal(buf, msg1)
	if err != nil {
		t.Fatal(err)
	}
	/*
		for _, item := range msg1.Items {
			fmt.Printf("item %#+v\n", item)
		}
	*/
}
