package codec

import (
	"fmt"
	"testing"

	jsoncodec "go.unistack.org/micro-codec-json/v3"
	"go.unistack.org/micro/v3/codec"
)

func TestFrame(t *testing.T) {
	type FrameStruct struct {
		Frame *codec.Frame `json:"frame"`
		Name  string       `json:"name"`
	}
	dst := &FrameStruct{}
	data := []byte(`{"name":"test","frame": {"first":"second"}}`)
	c := jsoncodec.NewCodec()

	if err := c.Unmarshal(data, dst); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("xxx %s\n", dst.Frame)
}
