// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v5.29.2
// source: test.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CallReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Req    string  `protobuf:"bytes,2,opt,name=req,proto3" json:"req,omitempty"`
	Arg1   string  `protobuf:"bytes,3,opt,name=arg1,proto3" json:"arg1,omitempty"`
	Arg2   uint64  `protobuf:"varint,4,opt,name=arg2,proto3" json:"arg2,omitempty"`
	Nested *Nested `protobuf:"bytes,5,opt,name=nested,proto3" json:"nested,omitempty"`
}

func (x *CallReq) Reset() {
	*x = CallReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CallReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallReq) ProtoMessage() {}

func (x *CallReq) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallReq.ProtoReflect.Descriptor instead.
func (*CallReq) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0}
}

func (x *CallReq) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CallReq) GetReq() string {
	if x != nil {
		return x.Req
	}
	return ""
}

func (x *CallReq) GetArg1() string {
	if x != nil {
		return x.Arg1
	}
	return ""
}

func (x *CallReq) GetArg2() uint64 {
	if x != nil {
		return x.Arg2
	}
	return 0
}

func (x *CallReq) GetNested() *Nested {
	if x != nil {
		return x.Nested
	}
	return nil
}

type Nested struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StringArgs []string                  `protobuf:"bytes,1,rep,name=string_args,json=stringArgs,proto3" json:"string_args,omitempty"`
	Uint64Args []*wrapperspb.UInt64Value `protobuf:"bytes,2,rep,name=uint64_args,json=uint64Args,proto3" json:"uint64_args,omitempty"`
}

func (x *Nested) Reset() {
	*x = Nested{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Nested) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Nested) ProtoMessage() {}

func (x *Nested) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Nested.ProtoReflect.Descriptor instead.
func (*Nested) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{1}
}

func (x *Nested) GetStringArgs() []string {
	if x != nil {
		return x.StringArgs
	}
	return nil
}

func (x *Nested) GetUint64Args() []*wrapperspb.UInt64Value {
	if x != nil {
		return x.Uint64Args
	}
	return nil
}

var File_test_proto protoreflect.FileDescriptor

var file_test_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x72, 0x65,
	0x66, 0x6c, 0x65, 0x63, 0x74, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x80, 0x01, 0x0a, 0x07, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65,
	0x71, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x71, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x72, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x72, 0x67, 0x31, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x72, 0x67, 0x31, 0x12, 0x12, 0x0a, 0x04, 0x61,
	0x72, 0x67, 0x32, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x61, 0x72, 0x67, 0x32, 0x12,
	0x27, 0x0a, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x72, 0x65, 0x66, 0x6c, 0x65, 0x63, 0x74, 0x2e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64,
	0x52, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x22, 0x68, 0x0a, 0x06, 0x4e, 0x65, 0x73, 0x74,
	0x65, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x61, 0x72, 0x67,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x41,
	0x72, 0x67, 0x73, 0x12, 0x3d, 0x0a, 0x0b, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x5f, 0x61, 0x72,
	0x67, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36,
	0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x41, 0x72,
	0x67, 0x73, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x6f, 0x2e, 0x75, 0x6e, 0x69, 0x73, 0x74, 0x61, 0x63,
	0x6b, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x2d, 0x74, 0x65, 0x73, 0x74,
	0x73, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x72, 0x65, 0x66, 0x6c, 0x65, 0x63, 0x74, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_test_proto_rawDescOnce sync.Once
	file_test_proto_rawDescData = file_test_proto_rawDesc
)

func file_test_proto_rawDescGZIP() []byte {
	file_test_proto_rawDescOnce.Do(func() {
		file_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_proto_rawDescData)
	})
	return file_test_proto_rawDescData
}

var file_test_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_test_proto_goTypes = []interface{}{
	(*CallReq)(nil),                // 0: reflect.CallReq
	(*Nested)(nil),                 // 1: reflect.Nested
	(*wrapperspb.UInt64Value)(nil), // 2: google.protobuf.UInt64Value
}
var file_test_proto_depIdxs = []int32{
	1, // 0: reflect.CallReq.nested:type_name -> reflect.Nested
	2, // 1: reflect.Nested.uint64_args:type_name -> google.protobuf.UInt64Value
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_test_proto_init() }
func file_test_proto_init() {
	if File_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CallReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_test_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Nested); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_proto_goTypes,
		DependencyIndexes: file_test_proto_depIdxs,
		MessageInfos:      file_test_proto_msgTypes,
	}.Build()
	File_test_proto = out.File
	file_test_proto_rawDesc = nil
	file_test_proto_goTypes = nil
	file_test_proto_depIdxs = nil
}
