// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: test.proto

package pb

import (
	reflect "reflect"
	sync "sync"

	_ "go.unistack.org/micro-proto/v3/api"
	_ "go.unistack.org/micro-proto/v3/openapiv3"
	_ "go.unistack.org/micro-proto/v3/tag"
	codec "go.unistack.org/micro/v3/codec"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RequestAml struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KzAmlNewOnlineRs *KZAmlNewOnlineRs `protobuf:"bytes,1,opt,name=kzAmlNewOnlineRs,proto3" json:"kzAmlNewOnlineRs,omitempty" xml:"KZAmlNewOnlineRs"`
}

func (x *RequestAml) Reset() {
	*x = RequestAml{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestAml) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestAml) ProtoMessage() {}

func (x *RequestAml) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use RequestAml.ProtoReflect.Descriptor instead.
func (*RequestAml) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0}
}

func (x *RequestAml) GetKzAmlNewOnlineRs() *KZAmlNewOnlineRs {
	if x != nil {
		return x.KzAmlNewOnlineRs
	}
	return nil
}

type KZAmlNewOnlineRs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rquid      string  `protobuf:"bytes,1,opt,name=rquid,proto3" json:"rquid,omitempty" xml:"RqUID"`
	Rqtm       string  `protobuf:"bytes,2,opt,name=rqtm,proto3" json:"rqtm,omitempty" xml:"RqTm"`
	Status     *Status `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty" xml:"Status"`
	TerrStatus int64   `protobuf:"varint,4,opt,name=terr_status,json=terrStatus,proto3" json:"terr_status,omitempty" xml:"TerrStatus"`
	AmlStatus  int64   `protobuf:"varint,5,opt,name=aml_status,json=amlStatus,proto3" json:"aml_status,omitempty" xml:"AMLStatus"`
}

func (x *KZAmlNewOnlineRs) Reset() {
	*x = KZAmlNewOnlineRs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KZAmlNewOnlineRs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KZAmlNewOnlineRs) ProtoMessage() {}

func (x *KZAmlNewOnlineRs) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use KZAmlNewOnlineRs.ProtoReflect.Descriptor instead.
func (*KZAmlNewOnlineRs) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{1}
}

func (x *KZAmlNewOnlineRs) GetRquid() string {
	if x != nil {
		return x.Rquid
	}
	return ""
}

func (x *KZAmlNewOnlineRs) GetRqtm() string {
	if x != nil {
		return x.Rqtm
	}
	return ""
}

func (x *KZAmlNewOnlineRs) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *KZAmlNewOnlineRs) GetTerrStatus() int64 {
	if x != nil {
		return x.TerrStatus
	}
	return 0
}

func (x *KZAmlNewOnlineRs) GetAmlStatus() int64 {
	if x != nil {
		return x.AmlStatus
	}
	return 0
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode int64 `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty" xml:"StatusCode"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{2}
}

func (x *Status) GetStatusCode() int64 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

type ResponseAml struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resp string `protobuf:"bytes,1,opt,name=resp,proto3" json:"resp,omitempty"`
}

func (x *ResponseAml) Reset() {
	*x = ResponseAml{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseAml) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseAml) ProtoMessage() {}

func (x *ResponseAml) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseAml.ProtoReflect.Descriptor instead.
func (*ResponseAml) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{3}
}

func (x *ResponseAml) GetResp() string {
	if x != nil {
		return x.Resp
	}
	return ""
}

type RequestImage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Image    []byte `protobuf:"bytes,1,opt,name=image,proto3" json:"image,omitempty"`
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	DocType  string `protobuf:"bytes,3,opt,name=doc_type,json=docType,proto3" json:"doc_type,omitempty"`
}

func (x *RequestImage) Reset() {
	*x = RequestImage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestImage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestImage) ProtoMessage() {}

func (x *RequestImage) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestImage.ProtoReflect.Descriptor instead.
func (*RequestImage) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{4}
}

func (x *RequestImage) GetImage() []byte {
	if x != nil {
		return x.Image
	}
	return nil
}

func (x *RequestImage) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *RequestImage) GetDocType() string {
	if x != nil {
		return x.DocType
	}
	return ""
}

type ResponseImage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ResponseImage) Reset() {
	*x = ResponseImage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseImage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseImage) ProtoMessage() {}

func (x *ResponseImage) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseImage.ProtoReflect.Descriptor instead.
func (*ResponseImage) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{5}
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{6}
}

func (x *Request) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncodedId string `protobuf:"bytes,1,opt,name=encoded_id,json=encodedId,proto3" json:"encoded_id,omitempty" xml:"encoded_id,attr"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{7}
}

func (x *Response) GetEncodedId() string {
	if x != nil {
		return x.EncodedId
	}
	return ""
}

var File_test_proto protoreflect.FileDescriptor

var file_test_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x74, 0x65,
	0x73, 0x74, 0x1a, 0x15, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x63, 0x6f, 0x64, 0x65, 0x63,
	0x2f, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x74, 0x61,
	0x67, 0x2f, 0x74, 0x61, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x6f, 0x70, 0x65,
	0x6e, 0x61, 0x70, 0x69, 0x76, 0x33, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a, 0x0a, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x41, 0x6d, 0x6c, 0x12, 0x5f, 0x0a, 0x10, 0x6b, 0x7a, 0x41, 0x6d, 0x6c, 0x4e,
	0x65, 0x77, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x4b, 0x5a, 0x41, 0x6d, 0x6c, 0x4e, 0x65, 0x77,
	0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x73, 0x42, 0x1b, 0x9a, 0x84, 0x9e, 0x03, 0x16, 0x78,
	0x6d, 0x6c, 0x3a, 0x22, 0x4b, 0x5a, 0x41, 0x6d, 0x6c, 0x4e, 0x65, 0x77, 0x4f, 0x6e, 0x6c, 0x69,
	0x6e, 0x65, 0x52, 0x73, 0x22, 0x52, 0x10, 0x6b, 0x7a, 0x41, 0x6d, 0x6c, 0x4e, 0x65, 0x77, 0x4f,
	0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x73, 0x22, 0x85, 0x02, 0x0a, 0x10, 0x4b, 0x5a, 0x41, 0x6d,
	0x6c, 0x4e, 0x65, 0x77, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x73, 0x12, 0x26, 0x0a, 0x05,
	0x72, 0x71, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x10, 0x9a, 0x84, 0x9e,
	0x03, 0x0b, 0x78, 0x6d, 0x6c, 0x3a, 0x22, 0x52, 0x71, 0x55, 0x49, 0x44, 0x22, 0x52, 0x05, 0x72,
	0x71, 0x75, 0x69, 0x64, 0x12, 0x23, 0x0a, 0x04, 0x72, 0x71, 0x74, 0x6d, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x0f, 0x9a, 0x84, 0x9e, 0x03, 0x0a, 0x78, 0x6d, 0x6c, 0x3a, 0x22, 0x52, 0x71,
	0x54, 0x6d, 0x22, 0x52, 0x04, 0x72, 0x71, 0x74, 0x6d, 0x12, 0x37, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x74, 0x65, 0x73, 0x74,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x11, 0x9a, 0x84, 0x9e, 0x03, 0x0c, 0x78, 0x6d,
	0x6c, 0x3a, 0x22, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x36, 0x0a, 0x0b, 0x74, 0x65, 0x72, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x42, 0x15, 0x9a, 0x84, 0x9e, 0x03, 0x10, 0x78, 0x6d,
	0x6c, 0x3a, 0x22, 0x54, 0x65, 0x72, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x52, 0x0a,
	0x74, 0x65, 0x72, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x33, 0x0a, 0x0a, 0x61, 0x6d,
	0x6c, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x42, 0x14,
	0x9a, 0x84, 0x9e, 0x03, 0x0f, 0x78, 0x6d, 0x6c, 0x3a, 0x22, 0x41, 0x4d, 0x4c, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x52, 0x09, 0x61, 0x6d, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x40, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x36, 0x0a, 0x0b, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x42, 0x15,
	0x9a, 0x84, 0x9e, 0x03, 0x10, 0x78, 0x6d, 0x6c, 0x3a, 0x22, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x43, 0x6f, 0x64, 0x65, 0x22, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64,
	0x65, 0x22, 0x21, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x41, 0x6d, 0x6c,
	0x12, 0x12, 0x0a, 0x04, 0x72, 0x65, 0x73, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x72, 0x65, 0x73, 0x70, 0x22, 0x5c, 0x0a, 0x0c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x6f, 0x63, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x6f, 0x63, 0x54, 0x79,
	0x70, 0x65, 0x22, 0x0f, 0x0a, 0x0d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x49, 0x6d,
	0x61, 0x67, 0x65, 0x22, 0x19, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x45,
	0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x6e,
	0x63, 0x6f, 0x64, 0x65, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1a,
	0x9a, 0x84, 0x9e, 0x03, 0x15, 0x78, 0x6d, 0x6c, 0x3a, 0x22, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65,
	0x64, 0x5f, 0x69, 0x64, 0x2c, 0x61, 0x74, 0x74, 0x72, 0x22, 0x52, 0x09, 0x65, 0x6e, 0x63, 0x6f,
	0x64, 0x65, 0x64, 0x49, 0x64, 0x32, 0xf2, 0x02, 0x0a, 0x0b, 0x54, 0x65, 0x73, 0x74, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x42, 0x0a, 0x0c, 0x54, 0x65, 0x73, 0x74, 0x45, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x0d, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x13, 0xb2, 0xea, 0xff, 0xf9, 0x01, 0x0d, 0x12, 0x0b, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x12, 0x3e, 0x0a, 0x08, 0x55, 0x73, 0x65,
	0x72, 0x42, 0x79, 0x49, 0x44, 0x12, 0x0d, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x13, 0xb2, 0xea, 0xff, 0xf9, 0x01, 0x0d, 0x12, 0x0b, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x4d, 0x0a, 0x0d, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x42, 0x79, 0x49, 0x44, 0x12, 0x0d, 0x2e, 0x74, 0x65, 0x73,
	0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x6d, 0x69, 0x63, 0x72,
	0x6f, 0x2e, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x2e, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x22, 0x19, 0xb2,
	0xea, 0xff, 0xf9, 0x01, 0x13, 0x12, 0x11, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x69,
	0x64, 0x7d, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x52, 0x0a, 0x0a, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x12, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x1a, 0x13, 0x2e, 0x74, 0x65, 0x73,
	0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x22,
	0x1b, 0xb2, 0xea, 0xff, 0xf9, 0x01, 0x15, 0x22, 0x13, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x2f, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x3c, 0x0a, 0x07,
	0x4b, 0x7a, 0x41, 0x6d, 0x6c, 0x52, 0x73, 0x12, 0x10, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x6d, 0x6c, 0x1a, 0x11, 0x2e, 0x74, 0x65, 0x73, 0x74,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x41, 0x6d, 0x6c, 0x22, 0x0c, 0xb2, 0xea,
	0xff, 0xf9, 0x01, 0x06, 0x22, 0x04, 0x2f, 0x61, 0x6d, 0x6c, 0x42, 0x30, 0x5a, 0x09, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x62, 0xaa, 0x84, 0x9e, 0x03, 0x20, 0x12, 0x1e, 0x0a, 0x0a,
	0x74, 0x65, 0x73, 0x74, 0x20, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x09, 0x74, 0x65, 0x73, 0x74,
	0x20, 0x64, 0x65, 0x73, 0x63, 0x32, 0x05, 0x30, 0x2e, 0x30, 0x2e, 0x35, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
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

var (
	file_test_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
	file_test_proto_goTypes  = []interface{}{
		(*RequestAml)(nil),       // 0: test.RequestAml
		(*KZAmlNewOnlineRs)(nil), // 1: test.KZAmlNewOnlineRs
		(*Status)(nil),           // 2: test.Status
		(*ResponseAml)(nil),      // 3: test.ResponseAml
		(*RequestImage)(nil),     // 4: test.RequestImage
		(*ResponseImage)(nil),    // 5: test.ResponseImage
		(*Request)(nil),          // 6: test.Request
		(*Response)(nil),         // 7: test.Response
		(*codec.Frame)(nil),      // 8: micro.codec.Frame
	}
)

var file_test_proto_depIdxs = []int32{
	1, // 0: test.RequestAml.kzAmlNewOnlineRs:type_name -> test.KZAmlNewOnlineRs
	2, // 1: test.KZAmlNewOnlineRs.status:type_name -> test.Status
	6, // 2: test.TestService.TestEndpoint:input_type -> test.Request
	6, // 3: test.TestService.UserByID:input_type -> test.Request
	6, // 4: test.TestService.UserImageByID:input_type -> test.Request
	4, // 5: test.TestService.UploadFile:input_type -> test.RequestImage
	0, // 6: test.TestService.KzAmlRs:input_type -> test.RequestAml
	7, // 7: test.TestService.TestEndpoint:output_type -> test.Response
	7, // 8: test.TestService.UserByID:output_type -> test.Response
	8, // 9: test.TestService.UserImageByID:output_type -> micro.codec.Frame
	5, // 10: test.TestService.UploadFile:output_type -> test.ResponseImage
	3, // 11: test.TestService.KzAmlRs:output_type -> test.ResponseAml
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
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
			switch v := v.(*RequestAml); i {
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
			switch v := v.(*KZAmlNewOnlineRs); i {
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
		file_test_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
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
		file_test_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseAml); i {
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
		file_test_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestImage); i {
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
		file_test_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseImage); i {
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
		file_test_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_test_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
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