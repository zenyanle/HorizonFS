// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.1
// source: pkg/proto/metadata_node.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ChunkInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChunkId       string                 `protobuf:"bytes,1,opt,name=chunk_id,json=chunkId,proto3" json:"chunk_id,omitempty"`
	ChunkLocation string                 `protobuf:"bytes,2,opt,name=chunk_location,json=chunkLocation,proto3" json:"chunk_location,omitempty"`
	ChunkStatus   string                 `protobuf:"bytes,3,opt,name=chunk_status,json=chunkStatus,proto3" json:"chunk_status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChunkInfo) Reset() {
	*x = ChunkInfo{}
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChunkInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChunkInfo) ProtoMessage() {}

func (x *ChunkInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChunkInfo.ProtoReflect.Descriptor instead.
func (*ChunkInfo) Descriptor() ([]byte, []int) {
	return file_pkg_proto_metadata_node_proto_rawDescGZIP(), []int{0}
}

func (x *ChunkInfo) GetChunkId() string {
	if x != nil {
		return x.ChunkId
	}
	return ""
}

func (x *ChunkInfo) GetChunkLocation() string {
	if x != nil {
		return x.ChunkLocation
	}
	return ""
}

func (x *ChunkInfo) GetChunkStatus() string {
	if x != nil {
		return x.ChunkStatus
	}
	return ""
}

type FileMetadata struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Chunks        []*ChunkInfo           `protobuf:"bytes,2,rep,name=chunks,proto3" json:"chunks,omitempty"` // <-- 这里是关键
	TotalChunks   int64                  `protobuf:"varint,3,opt,name=total_chunks,json=totalChunks,proto3" json:"total_chunks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FileMetadata) Reset() {
	*x = FileMetadata{}
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileMetadata) ProtoMessage() {}

func (x *FileMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileMetadata.ProtoReflect.Descriptor instead.
func (*FileMetadata) Descriptor() ([]byte, []int) {
	return file_pkg_proto_metadata_node_proto_rawDescGZIP(), []int{1}
}

func (x *FileMetadata) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *FileMetadata) GetChunks() []*ChunkInfo {
	if x != nil {
		return x.Chunks
	}
	return nil
}

func (x *FileMetadata) GetTotalChunks() int64 {
	if x != nil {
		return x.TotalChunks
	}
	return 0
}

type ProposeSetRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Index         int64                  `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Value         *ChunkInfo             `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProposeSetRequest) Reset() {
	*x = ProposeSetRequest{}
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProposeSetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposeSetRequest) ProtoMessage() {}

func (x *ProposeSetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposeSetRequest.ProtoReflect.Descriptor instead.
func (*ProposeSetRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_metadata_node_proto_rawDescGZIP(), []int{2}
}

func (x *ProposeSetRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ProposeSetRequest) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *ProposeSetRequest) GetValue() *ChunkInfo {
	if x != nil {
		return x.Value
	}
	return nil
}

type ProposeSetResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage  string                 `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProposeSetResponse) Reset() {
	*x = ProposeSetResponse{}
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProposeSetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposeSetResponse) ProtoMessage() {}

func (x *ProposeSetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposeSetResponse.ProtoReflect.Descriptor instead.
func (*ProposeSetResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_metadata_node_proto_rawDescGZIP(), []int{3}
}

func (x *ProposeSetResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *ProposeSetResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type GetMetadataRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Index         int64                  `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMetadataRequest) Reset() {
	*x = GetMetadataRequest{}
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMetadataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetadataRequest) ProtoMessage() {}

func (x *GetMetadataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetadataRequest.ProtoReflect.Descriptor instead.
func (*GetMetadataRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_metadata_node_proto_rawDescGZIP(), []int{4}
}

func (x *GetMetadataRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *GetMetadataRequest) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

type GetMetadataResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Value         *ChunkInfo             `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Found         bool                   `protobuf:"varint,2,opt,name=found,proto3" json:"found,omitempty"`
	ErrorMessage  string                 `protobuf:"bytes,3,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMetadataResponse) Reset() {
	*x = GetMetadataResponse{}
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMetadataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetadataResponse) ProtoMessage() {}

func (x *GetMetadataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetadataResponse.ProtoReflect.Descriptor instead.
func (*GetMetadataResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_metadata_node_proto_rawDescGZIP(), []int{5}
}

func (x *GetMetadataResponse) GetValue() *ChunkInfo {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *GetMetadataResponse) GetFound() bool {
	if x != nil {
		return x.Found
	}
	return false
}

func (x *GetMetadataResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type ProposeDeleteRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Index         int64                  `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProposeDeleteRequest) Reset() {
	*x = ProposeDeleteRequest{}
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProposeDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposeDeleteRequest) ProtoMessage() {}

func (x *ProposeDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposeDeleteRequest.ProtoReflect.Descriptor instead.
func (*ProposeDeleteRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_metadata_node_proto_rawDescGZIP(), []int{6}
}

func (x *ProposeDeleteRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ProposeDeleteRequest) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

type ProposeDeleteResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorMessage  string                 `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProposeDeleteResponse) Reset() {
	*x = ProposeDeleteResponse{}
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProposeDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposeDeleteResponse) ProtoMessage() {}

func (x *ProposeDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposeDeleteResponse.ProtoReflect.Descriptor instead.
func (*ProposeDeleteResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_metadata_node_proto_rawDescGZIP(), []int{7}
}

func (x *ProposeDeleteResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *ProposeDeleteResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type GetFileMetadataRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFileMetadataRequest) Reset() {
	*x = GetFileMetadataRequest{}
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFileMetadataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileMetadataRequest) ProtoMessage() {}

func (x *GetFileMetadataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileMetadataRequest.ProtoReflect.Descriptor instead.
func (*GetFileMetadataRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_metadata_node_proto_rawDescGZIP(), []int{8}
}

func (x *GetFileMetadataRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type GetFileMetadataResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Metadata      *FileMetadata          `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"` // <-- 接收包含 repeated 字段的消息
	Found         bool                   `protobuf:"varint,2,opt,name=found,proto3" json:"found,omitempty"`
	ErrorMessage  string                 `protobuf:"bytes,3,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFileMetadataResponse) Reset() {
	*x = GetFileMetadataResponse{}
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFileMetadataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileMetadataResponse) ProtoMessage() {}

func (x *GetFileMetadataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_metadata_node_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileMetadataResponse.ProtoReflect.Descriptor instead.
func (*GetFileMetadataResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_metadata_node_proto_rawDescGZIP(), []int{9}
}

func (x *GetFileMetadataResponse) GetMetadata() *FileMetadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *GetFileMetadataResponse) GetFound() bool {
	if x != nil {
		return x.Found
	}
	return false
}

func (x *GetFileMetadataResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

var File_pkg_proto_metadata_node_proto protoreflect.FileDescriptor

const file_pkg_proto_metadata_node_proto_rawDesc = "" +
	"\n" +
	"\x1dpkg/proto/metadata_node.proto\x12\x05proto\"p\n" +
	"\tChunkInfo\x12\x19\n" +
	"\bchunk_id\x18\x01 \x01(\tR\achunkId\x12%\n" +
	"\x0echunk_location\x18\x02 \x01(\tR\rchunkLocation\x12!\n" +
	"\fchunk_status\x18\x03 \x01(\tR\vchunkStatus\"m\n" +
	"\fFileMetadata\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12(\n" +
	"\x06chunks\x18\x02 \x03(\v2\x10.proto.ChunkInfoR\x06chunks\x12!\n" +
	"\ftotal_chunks\x18\x03 \x01(\x03R\vtotalChunks\"c\n" +
	"\x11ProposeSetRequest\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05index\x18\x02 \x01(\x03R\x05index\x12&\n" +
	"\x05value\x18\x03 \x01(\v2\x10.proto.ChunkInfoR\x05value\"S\n" +
	"\x12ProposeSetResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12#\n" +
	"\rerror_message\x18\x02 \x01(\tR\ferrorMessage\"<\n" +
	"\x12GetMetadataRequest\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05index\x18\x02 \x01(\x03R\x05index\"x\n" +
	"\x13GetMetadataResponse\x12&\n" +
	"\x05value\x18\x01 \x01(\v2\x10.proto.ChunkInfoR\x05value\x12\x14\n" +
	"\x05found\x18\x02 \x01(\bR\x05found\x12#\n" +
	"\rerror_message\x18\x03 \x01(\tR\ferrorMessage\">\n" +
	"\x14ProposeDeleteRequest\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05index\x18\x02 \x01(\x03R\x05index\"V\n" +
	"\x15ProposeDeleteResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12#\n" +
	"\rerror_message\x18\x02 \x01(\tR\ferrorMessage\"*\n" +
	"\x16GetFileMetadataRequest\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\"\x85\x01\n" +
	"\x17GetFileMetadataResponse\x12/\n" +
	"\bmetadata\x18\x01 \x01(\v2\x13.proto.FileMetadataR\bmetadata\x12\x14\n" +
	"\x05found\x18\x02 \x01(\bR\x05found\x12#\n" +
	"\rerror_message\x18\x03 \x01(\tR\ferrorMessage2\xb8\x02\n" +
	"\x0fMetadataService\x12A\n" +
	"\n" +
	"ProposeSet\x12\x18.proto.ProposeSetRequest\x1a\x19.proto.ProposeSetResponse\x12D\n" +
	"\vGetMetadata\x12\x19.proto.GetMetadataRequest\x1a\x1a.proto.GetMetadataResponse\x12J\n" +
	"\rProposeDelete\x12\x1b.proto.ProposeDeleteRequest\x1a\x1c.proto.ProposeDeleteResponse\x12P\n" +
	"\x0fGetFileMetadata\x12\x1d.proto.GetFileMetadataRequest\x1a\x1e.proto.GetFileMetadataResponseB\fZ\n" +
	"/pkg/protob\x06proto3"

var (
	file_pkg_proto_metadata_node_proto_rawDescOnce sync.Once
	file_pkg_proto_metadata_node_proto_rawDescData []byte
)

func file_pkg_proto_metadata_node_proto_rawDescGZIP() []byte {
	file_pkg_proto_metadata_node_proto_rawDescOnce.Do(func() {
		file_pkg_proto_metadata_node_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pkg_proto_metadata_node_proto_rawDesc), len(file_pkg_proto_metadata_node_proto_rawDesc)))
	})
	return file_pkg_proto_metadata_node_proto_rawDescData
}

var file_pkg_proto_metadata_node_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_pkg_proto_metadata_node_proto_goTypes = []any{
	(*ChunkInfo)(nil),               // 0: proto.ChunkInfo
	(*FileMetadata)(nil),            // 1: proto.FileMetadata
	(*ProposeSetRequest)(nil),       // 2: proto.ProposeSetRequest
	(*ProposeSetResponse)(nil),      // 3: proto.ProposeSetResponse
	(*GetMetadataRequest)(nil),      // 4: proto.GetMetadataRequest
	(*GetMetadataResponse)(nil),     // 5: proto.GetMetadataResponse
	(*ProposeDeleteRequest)(nil),    // 6: proto.ProposeDeleteRequest
	(*ProposeDeleteResponse)(nil),   // 7: proto.ProposeDeleteResponse
	(*GetFileMetadataRequest)(nil),  // 8: proto.GetFileMetadataRequest
	(*GetFileMetadataResponse)(nil), // 9: proto.GetFileMetadataResponse
}
var file_pkg_proto_metadata_node_proto_depIdxs = []int32{
	0, // 0: proto.FileMetadata.chunks:type_name -> proto.ChunkInfo
	0, // 1: proto.ProposeSetRequest.value:type_name -> proto.ChunkInfo
	0, // 2: proto.GetMetadataResponse.value:type_name -> proto.ChunkInfo
	1, // 3: proto.GetFileMetadataResponse.metadata:type_name -> proto.FileMetadata
	2, // 4: proto.MetadataService.ProposeSet:input_type -> proto.ProposeSetRequest
	4, // 5: proto.MetadataService.GetMetadata:input_type -> proto.GetMetadataRequest
	6, // 6: proto.MetadataService.ProposeDelete:input_type -> proto.ProposeDeleteRequest
	8, // 7: proto.MetadataService.GetFileMetadata:input_type -> proto.GetFileMetadataRequest
	3, // 8: proto.MetadataService.ProposeSet:output_type -> proto.ProposeSetResponse
	5, // 9: proto.MetadataService.GetMetadata:output_type -> proto.GetMetadataResponse
	7, // 10: proto.MetadataService.ProposeDelete:output_type -> proto.ProposeDeleteResponse
	9, // 11: proto.MetadataService.GetFileMetadata:output_type -> proto.GetFileMetadataResponse
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_pkg_proto_metadata_node_proto_init() }
func file_pkg_proto_metadata_node_proto_init() {
	if File_pkg_proto_metadata_node_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pkg_proto_metadata_node_proto_rawDesc), len(file_pkg_proto_metadata_node_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_proto_metadata_node_proto_goTypes,
		DependencyIndexes: file_pkg_proto_metadata_node_proto_depIdxs,
		MessageInfos:      file_pkg_proto_metadata_node_proto_msgTypes,
	}.Build()
	File_pkg_proto_metadata_node_proto = out.File
	file_pkg_proto_metadata_node_proto_goTypes = nil
	file_pkg_proto_metadata_node_proto_depIdxs = nil
}
