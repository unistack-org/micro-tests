// Code generated by protoc-gen-go-micro. DO NOT EDIT.
// protoc-gen-go-micro version: v3.10.4

package pb

import (
	protojson "google.golang.org/protobuf/encoding/protojson"
)

var (
	marshaler = protojson.MarshalOptions{}
)

func (m *Error) Error() string {
	buf, _ := marshaler.Marshal(m)
	return string(buf)
}
