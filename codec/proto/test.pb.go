// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.6
// source: test.proto

package pb

import (
	codec "github.com/unistack-org/micro/v3/codec"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_test_proto protoreflect.FileDescriptor

var file_test_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x68, 0x65,
	0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x1a, 0x11, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x2f,
	0x66, 0x72, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x38, 0x0a, 0x04, 0x54,
	0x65, 0x73, 0x74, 0x12, 0x30, 0x0a, 0x04, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x12, 0x2e, 0x6d, 0x69,
	0x63, 0x72, 0x6f, 0x2e, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x2e, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x1a,
	0x12, 0x2e, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x2e, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x2e, 0x46, 0x72,
	0x61, 0x6d, 0x65, 0x22, 0x00, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6e, 0x69, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x2d, 0x6f, 0x72, 0x67,
	0x2f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x2d, 0x74, 0x65, 0x73, 0x74, 0x73, 0x2f, 0x63, 0x6f, 0x64,
	0x65, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var file_test_proto_goTypes = []interface{}{
	(*codec.Frame)(nil), // 0: micro.codec.Frame
}
var file_test_proto_depIdxs = []int32{
	0, // 0: helloworld.Test.Call:input_type -> micro.codec.Frame
	0, // 1: helloworld.Test.Call:output_type -> micro.codec.Frame
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_test_proto_init() }
func file_test_proto_init() {
	if File_test_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_test_proto_goTypes,
		DependencyIndexes: file_test_proto_depIdxs,
	}.Build()
	File_test_proto = out.File
	file_test_proto_rawDesc = nil
	file_test_proto_goTypes = nil
	file_test_proto_depIdxs = nil
}