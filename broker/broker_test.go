package broker

import (
	"testing"

	jsoncodec "github.com/unistack-org/micro-codec-json/v3"
	"github.com/unistack-org/micro/v3/broker"
)

func TestBrokerMessage(t *testing.T) {
	c := jsoncodec.NewCodec()

	buf := []byte(`{"Header":{"Content-Type":"application\/json","Micro-Topic":"notify"},"Body":{"Class":{}}}`)
	msg := &broker.Message{}

	if err := c.Unmarshal(buf, msg); err != nil {
		t.Fatal(err)
	}
}
