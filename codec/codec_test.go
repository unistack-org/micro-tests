package codec_test

import (
	"testing"

	grpc "go.unistack.org/micro-codec-grpc/v3"
	json "go.unistack.org/micro-codec-json/v3"
	proto "go.unistack.org/micro-codec-proto/v3"
	"go.unistack.org/micro/v3/codec"
)

type testRWC struct{}

func (rwc *testRWC) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (rwc *testRWC) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (rwc *testRWC) Close() error {
	return nil
}

func getCodecs() map[string]codec.Codec {
	return map[string]codec.Codec{
		"bytes": codec.NewCodec(),
		"grpc":  grpc.NewCodec(),
		"json":  json.NewCodec(),
		"proto": proto.NewCodec(),
	}
}

func Test_WriteEmptyBody(t *testing.T) {
	rw := &testRWC{}
	for name, c := range getCodecs() {
		err := c.Write(rw, &codec.Message{
			Type:   codec.Error,
			Header: map[string]string{},
		}, nil)
		if err != nil {
			t.Fatalf("codec %s - expected no error when writing empty/nil body: %s", name, err)
		}
	}
}
