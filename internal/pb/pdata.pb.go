// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.4
// source: pdata.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Pdata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pname   string `protobuf:"bytes,1,opt,name=pname,proto3" json:"pname,omitempty"`
	Ptype   string `protobuf:"bytes,2,opt,name=ptype,proto3" json:"ptype,omitempty"`
	Pdata   []byte `protobuf:"bytes,3,opt,name=pdata,proto3" json:"pdata,omitempty"`
	KeyHash []byte `protobuf:"bytes,4,opt,name=keyHash,proto3" json:"keyHash,omitempty"`
	ID      int64  `protobuf:"varint,5,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *Pdata) Reset() {
	*x = Pdata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pdata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pdata) ProtoMessage() {}

func (x *Pdata) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pdata.ProtoReflect.Descriptor instead.
func (*Pdata) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{0}
}

func (x *Pdata) GetPname() string {
	if x != nil {
		return x.Pname
	}
	return ""
}

func (x *Pdata) GetPtype() string {
	if x != nil {
		return x.Ptype
	}
	return ""
}

func (x *Pdata) GetPdata() []byte {
	if x != nil {
		return x.Pdata
	}
	return nil
}

func (x *Pdata) GetKeyHash() []byte {
	if x != nil {
		return x.KeyHash
	}
	return nil
}

func (x *Pdata) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

type PdataEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ID   int64  `protobuf:"varint,2,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *PdataEntry) Reset() {
	*x = PdataEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PdataEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PdataEntry) ProtoMessage() {}

func (x *PdataEntry) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PdataEntry.ProtoReflect.Descriptor instead.
func (*PdataEntry) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{1}
}

func (x *PdataEntry) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PdataEntry) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

// Add
type AddPdataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pdata *Pdata `protobuf:"bytes,1,opt,name=pdata,proto3" json:"pdata,omitempty"`
}

func (x *AddPdataRequest) Reset() {
	*x = AddPdataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddPdataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddPdataRequest) ProtoMessage() {}

func (x *AddPdataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddPdataRequest.ProtoReflect.Descriptor instead.
func (*AddPdataRequest) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{2}
}

func (x *AddPdataRequest) GetPdata() *Pdata {
	if x != nil {
		return x.Pdata
	}
	return nil
}

type AddPdataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *AddPdataResponse) Reset() {
	*x = AddPdataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddPdataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddPdataResponse) ProtoMessage() {}

func (x *AddPdataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddPdataResponse.ProtoReflect.Descriptor instead.
func (*AddPdataResponse) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{3}
}

func (x *AddPdataResponse) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

// Get
type GetPdataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pname string `protobuf:"bytes,1,opt,name=pname,proto3" json:"pname,omitempty"`
}

func (x *GetPdataRequest) Reset() {
	*x = GetPdataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPdataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPdataRequest) ProtoMessage() {}

func (x *GetPdataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPdataRequest.ProtoReflect.Descriptor instead.
func (*GetPdataRequest) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{4}
}

func (x *GetPdataRequest) GetPname() string {
	if x != nil {
		return x.Pname
	}
	return ""
}

type GetPdataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pdata *Pdata `protobuf:"bytes,1,opt,name=pdata,proto3" json:"pdata,omitempty"`
}

func (x *GetPdataResponse) Reset() {
	*x = GetPdataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPdataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPdataResponse) ProtoMessage() {}

func (x *GetPdataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPdataResponse.ProtoReflect.Descriptor instead.
func (*GetPdataResponse) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{5}
}

func (x *GetPdataResponse) GetPdata() *Pdata {
	if x != nil {
		return x.Pdata
	}
	return nil
}

// Update
type UpdatePdataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pdata *Pdata `protobuf:"bytes,1,opt,name=pdata,proto3" json:"pdata,omitempty"`
}

func (x *UpdatePdataRequest) Reset() {
	*x = UpdatePdataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePdataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePdataRequest) ProtoMessage() {}

func (x *UpdatePdataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePdataRequest.ProtoReflect.Descriptor instead.
func (*UpdatePdataRequest) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{6}
}

func (x *UpdatePdataRequest) GetPdata() *Pdata {
	if x != nil {
		return x.Pdata
	}
	return nil
}

type UpdatePdataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *UpdatePdataResponse) Reset() {
	*x = UpdatePdataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePdataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePdataResponse) ProtoMessage() {}

func (x *UpdatePdataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePdataResponse.ProtoReflect.Descriptor instead.
func (*UpdatePdataResponse) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{7}
}

func (x *UpdatePdataResponse) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

// List
type ListPdataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ptype string `protobuf:"bytes,1,opt,name=ptype,proto3" json:"ptype,omitempty"`
}

func (x *ListPdataRequest) Reset() {
	*x = ListPdataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPdataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPdataRequest) ProtoMessage() {}

func (x *ListPdataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPdataRequest.ProtoReflect.Descriptor instead.
func (*ListPdataRequest) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{8}
}

func (x *ListPdataRequest) GetPtype() string {
	if x != nil {
		return x.Ptype
	}
	return ""
}

type ListPdataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PdataEtnry []*PdataEntry `protobuf:"bytes,1,rep,name=pdataEtnry,proto3" json:"pdataEtnry,omitempty"`
}

func (x *ListPdataResponse) Reset() {
	*x = ListPdataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPdataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPdataResponse) ProtoMessage() {}

func (x *ListPdataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPdataResponse.ProtoReflect.Descriptor instead.
func (*ListPdataResponse) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{9}
}

func (x *ListPdataResponse) GetPdataEtnry() []*PdataEntry {
	if x != nil {
		return x.PdataEtnry
	}
	return nil
}

// Delete
type DeletePdataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PdataID int64 `protobuf:"varint,1,opt,name=pdataID,proto3" json:"pdataID,omitempty"`
}

func (x *DeletePdataRequest) Reset() {
	*x = DeletePdataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeletePdataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePdataRequest) ProtoMessage() {}

func (x *DeletePdataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeletePdataRequest.ProtoReflect.Descriptor instead.
func (*DeletePdataRequest) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{10}
}

func (x *DeletePdataRequest) GetPdataID() int64 {
	if x != nil {
		return x.PdataID
	}
	return 0
}

type DeletePdataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *DeletePdataResponse) Reset() {
	*x = DeletePdataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pdata_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeletePdataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePdataResponse) ProtoMessage() {}

func (x *DeletePdataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pdata_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeletePdataResponse.ProtoReflect.Descriptor instead.
func (*DeletePdataResponse) Descriptor() ([]byte, []int) {
	return file_pdata_proto_rawDescGZIP(), []int{11}
}

func (x *DeletePdataResponse) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

var File_pdata_proto protoreflect.FileDescriptor

var file_pdata_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x70, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73, 0x0a,
	0x05, 0x50, 0x64, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x70, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x70, 0x64, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x6b, 0x65, 0x79, 0x48,
	0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x48, 0x61,
	0x73, 0x68, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x49, 0x44, 0x22, 0x30, 0x0a, 0x0a, 0x50, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x49, 0x44, 0x22, 0x2f, 0x0a, 0x0f, 0x41, 0x64, 0x64, 0x50, 0x64, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x05, 0x70, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x50, 0x64, 0x61, 0x74, 0x61, 0x52, 0x05,
	0x70, 0x64, 0x61, 0x74, 0x61, 0x22, 0x2e, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x50, 0x64, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x27, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x50, 0x64, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x30,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x50, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1c, 0x0a, 0x05, 0x70, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x06, 0x2e, 0x50, 0x64, 0x61, 0x74, 0x61, 0x52, 0x05, 0x70, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x32, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x05, 0x70, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x50, 0x64, 0x61, 0x74, 0x61, 0x52, 0x05, 0x70,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x31, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x64,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x28, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x50,
	0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x70,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x74, 0x79, 0x70,
	0x65, 0x22, 0x40, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x0a, 0x70, 0x64, 0x61, 0x74, 0x61, 0x45,
	0x74, 0x6e, 0x72, 0x79, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x50, 0x64, 0x61,
	0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x70, 0x64, 0x61, 0x74, 0x61, 0x45, 0x74,
	0x6e, 0x72, 0x79, 0x22, 0x2e, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x64, 0x61,
	0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x64, 0x61,
	0x74, 0x61, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x70, 0x64, 0x61, 0x74,
	0x61, 0x49, 0x44, 0x22, 0x31, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x64, 0x61,
	0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xa1, 0x02, 0x0a, 0x0b, 0x50, 0x72, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x31, 0x0a, 0x08, 0x41, 0x64, 0x64, 0x50, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x10, 0x2e, 0x41, 0x64, 0x64, 0x50, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x41, 0x64, 0x64, 0x50, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x31, 0x0a, 0x08, 0x47, 0x65, 0x74,
	0x50, 0x64, 0x61, 0x74, 0x61, 0x12, 0x10, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x64, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x64, 0x61,
	0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x0b,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x64, 0x61, 0x74, 0x61, 0x12, 0x13, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x50, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x14, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74,
	0x50, 0x64, 0x61, 0x74, 0x61, 0x12, 0x11, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x64, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50,
	0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3a,
	0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x64, 0x61, 0x74, 0x61, 0x12, 0x13, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x14, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x64, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pdata_proto_rawDescOnce sync.Once
	file_pdata_proto_rawDescData = file_pdata_proto_rawDesc
)

func file_pdata_proto_rawDescGZIP() []byte {
	file_pdata_proto_rawDescOnce.Do(func() {
		file_pdata_proto_rawDescData = protoimpl.X.CompressGZIP(file_pdata_proto_rawDescData)
	})
	return file_pdata_proto_rawDescData
}

var file_pdata_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_pdata_proto_goTypes = []interface{}{
	(*Pdata)(nil),               // 0: Pdata
	(*PdataEntry)(nil),          // 1: PdataEntry
	(*AddPdataRequest)(nil),     // 2: AddPdataRequest
	(*AddPdataResponse)(nil),    // 3: AddPdataResponse
	(*GetPdataRequest)(nil),     // 4: GetPdataRequest
	(*GetPdataResponse)(nil),    // 5: GetPdataResponse
	(*UpdatePdataRequest)(nil),  // 6: UpdatePdataRequest
	(*UpdatePdataResponse)(nil), // 7: UpdatePdataResponse
	(*ListPdataRequest)(nil),    // 8: ListPdataRequest
	(*ListPdataResponse)(nil),   // 9: ListPdataResponse
	(*DeletePdataRequest)(nil),  // 10: DeletePdataRequest
	(*DeletePdataResponse)(nil), // 11: DeletePdataResponse
}
var file_pdata_proto_depIdxs = []int32{
	0,  // 0: AddPdataRequest.pdata:type_name -> Pdata
	0,  // 1: GetPdataResponse.pdata:type_name -> Pdata
	0,  // 2: UpdatePdataRequest.pdata:type_name -> Pdata
	1,  // 3: ListPdataResponse.pdataEtnry:type_name -> PdataEntry
	2,  // 4: PrivateData.AddPdata:input_type -> AddPdataRequest
	4,  // 5: PrivateData.GetPdata:input_type -> GetPdataRequest
	6,  // 6: PrivateData.UpdatePdata:input_type -> UpdatePdataRequest
	8,  // 7: PrivateData.ListPdata:input_type -> ListPdataRequest
	10, // 8: PrivateData.DeletePdata:input_type -> DeletePdataRequest
	3,  // 9: PrivateData.AddPdata:output_type -> AddPdataResponse
	5,  // 10: PrivateData.GetPdata:output_type -> GetPdataResponse
	7,  // 11: PrivateData.UpdatePdata:output_type -> UpdatePdataResponse
	9,  // 12: PrivateData.ListPdata:output_type -> ListPdataResponse
	11, // 13: PrivateData.DeletePdata:output_type -> DeletePdataResponse
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_pdata_proto_init() }
func file_pdata_proto_init() {
	if File_pdata_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pdata_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pdata); i {
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
		file_pdata_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PdataEntry); i {
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
		file_pdata_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddPdataRequest); i {
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
		file_pdata_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddPdataResponse); i {
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
		file_pdata_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPdataRequest); i {
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
		file_pdata_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPdataResponse); i {
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
		file_pdata_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePdataRequest); i {
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
		file_pdata_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePdataResponse); i {
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
		file_pdata_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPdataRequest); i {
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
		file_pdata_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPdataResponse); i {
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
		file_pdata_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeletePdataRequest); i {
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
		file_pdata_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeletePdataResponse); i {
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
			RawDescriptor: file_pdata_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pdata_proto_goTypes,
		DependencyIndexes: file_pdata_proto_depIdxs,
		MessageInfos:      file_pdata_proto_msgTypes,
	}.Build()
	File_pdata_proto = out.File
	file_pdata_proto_rawDesc = nil
	file_pdata_proto_goTypes = nil
	file_pdata_proto_depIdxs = nil
}
